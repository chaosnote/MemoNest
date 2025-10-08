package controllers

import (
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/service"
)

type ToolsController struct{}

// 顯示節點路徑
func (u *ToolsController) show_node(di service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		const parent_id = "00000000-0000-0000-0000-000000000000"
		allCategories := di.Mariadb.GetCategories()

		nodesMap := make(map[string]*model.CategoryNode)
		var rootNodes []*model.CategoryNode

		// 第一次遍歷：建立節點地圖
		for _, cat := range allCategories {
			nodesMap[cat.NodeID] = &model.CategoryNode{
				Category: cat,
			}
		}

		// 第二次遍歷：建立樹狀結構並生成路徑
		for _, cat := range allCategories {
			currentNode := nodesMap[cat.NodeID]

			// 建立完整路徑
			pathParts := []string{currentNode.PathName}
			tempNode := currentNode
			for {
				if tempNode.ParentID == parent_id {
					break
				}
				if parent, ok := nodesMap[tempNode.ParentID]; ok {
					pathParts = append([]string{parent.PathName}, pathParts...) // 將父節點名稱加到最前面
					tempNode = parent
				} else {
					break // 父節點不存在，停止回溯
				}
			}
			currentNode.Path = "/" + strings.Join(pathParts, "/")

			// 處理樹狀結構
			if cat.ParentID == parent_id {
				rootNodes = append(rootNodes, currentNode)
			} else {
				if parent, ok := nodesMap[cat.ParentID]; ok {
					parent.Children = append(parent.Children, currentNode)
				}
			}
		}

		data := gin.H{"Categories": rootNodes}

		templates := template.Must(template.ParseFiles(filepath.Join("./assets", "templates", "tools", "show_node.html")))
		err := templates.Execute(c.Writer, data)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (u *ToolsController) gen_uuid(service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error": "",
			"uuid":  uuid.NewString(),
		})
	}
}

func NewToolsController(rg *gin.RouterGroup, di service.DI) {
	c := &ToolsController{}
	r := rg.Group("/tools")
	r.GET("/show_node", c.show_node(di))
	r.GET("/uuid", c.gen_uuid(di))
}
