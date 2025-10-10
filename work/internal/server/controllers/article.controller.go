package controllers

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"
)

type ArticleHelper struct {
	db *sql.DB
}

func (nh *ArticleHelper) getAllNode() (categories []model.Category) {
	rows, err := nh.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return
		}
		categories = append(categories, c)
	}

	return
}

//-----------------------------------------------

type ArticleController struct {
	helper *ArticleHelper
}

func (u *ArticleController) add(c *gin.Context) {
	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "article", "add.html")},
		Pattern: []string{},
	}
	tmpl, e := utils.RenderTemplate(config)
	defer func() {
		if e != nil {
			http.Error(c.Writer, e.Error(), http.StatusInternalServerError)
		}
	}()
	if e != nil {
		return
	}

	root_id := uuid.Nil.String()
	source := u.helper.getAllNode()
	node_map := make(map[string]*model.CategoryNode)

	// 第一次遍歷：建立節點地圖
	for _, cat := range source {
		node_map[cat.NodeID] = &model.CategoryNode{
			Category: cat,
		}
	}
	// 第二次遍歷：建立樹狀結構並生成路徑
	for _, cat := range source {
		current_node := node_map[cat.NodeID]

		// 建立完整路徑
		path_seg := []string{current_node.PathName}
		temp_node := current_node
		for {
			if temp_node.ParentID == root_id {
				break
			}
			parent, ok := node_map[temp_node.ParentID]
			if !ok {
				break
			}
			path_seg = append([]string{parent.PathName}, path_seg...) // 將父節點名稱加到最前面
			temp_node = parent
		}
		current_node.Path = "/" + strings.Join(path_seg, "/")
	}

	e = tmpl.ExecuteTemplate(c.Writer, "add.html", gin.H{
		"Title":   "增加文章",
		"NodeMap": node_map,
	})
	if e != nil {
		return
	}
}

func NewArticleController(rg *gin.RouterGroup, di service.DI) {
	c := &ArticleController{
		helper: &ArticleHelper{
			db: di.MariaDB,
		},
	}
	r := rg.Group("/article")
	r.GET("/add", c.add)
}
