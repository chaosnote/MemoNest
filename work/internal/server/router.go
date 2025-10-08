package server

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/server/middleware"
	"idv/chris/MemoNest/internal/service"
)

// RegisterRoutes 註冊所有路由
func RegisterRoutes(engine *gin.Engine, logger *zap.Logger, deps service.Deps) {
	engine.GET("/", vesion(deps))
	engine.GET("/health", health(deps))

	r := engine.Group("/tools", middleware.IPCheckMiddleware()) // 工具
	r.GET("/uuid", gen_uuid(deps))

	r = engine.Group("/api/v1") // 對外 api
	r.GET("/show_list", show_list(deps))
	r.POST("/login", login(deps))
	r.POST("/logout", logout(deps))
}

//-----------------------------------------------

func vesion(service.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"error": "", "message": "Test"})
	}
}

func health(service.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"error": "", "message": "OK"})
	}
}

//-----------------------------------------------

func gen_uuid(service.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error": "",
			"uuid":  uuid.NewString(),
		})
	}
}

//-----------------------------------------------

func show_list(deps service.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		allCategories := deps.Mariadb.GetCategories()

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
				if tempNode.ParentID == "00000000-0000-0000-0000-000000000000" {
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
			if cat.ParentID == "00000000-0000-0000-0000-000000000000" {
				rootNodes = append(rootNodes, currentNode)
			} else {
				if parent, ok := nodesMap[cat.ParentID]; ok {
					parent.Children = append(parent.Children, currentNode)
				}
			}
		}

		data := gin.H{"Categories": rootNodes}

		templates := template.Must(template.ParseFiles(filepath.Join("./assets", "templates", "categories.html")))
		err := templates.Execute(c.Writer, data)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

// 註冊或登入玩家
func login(service.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// 玩家登出
func logout(service.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 gin 取得 POST JSON 資料
	}
}
