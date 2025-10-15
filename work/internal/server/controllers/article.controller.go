package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/server/controllers/share"
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

func (nh *ArticleHelper) add(title, content, node_id string) error {
	_, e := nh.db.Exec("CALL insert_article(?, ?, NOW(), NOW(), ?) ;", title, content, node_id)
	if e != nil {
		return e
	}
	return e
}

func (nh *ArticleHelper) get(id int) (articles []model.Article) {
	rows, e := nh.db.Query(`
		SELECT 
			a.RowID AS ArticleRowID,
			a.Title,
			a.Content,
			a.NodeID,
			c.PathName,
			a.UpdateDt
		FROM articles as a
		JOIN categories as c ON a.NodeID = c.NodeID
		WHERE a.RowID = ? ;
	`, id)
	if e != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		e = rows.Scan(
			&article.ArticleRowID,
			&article.Title,
			&article.Content,
			&article.NodeID,
			&article.PathName,
			&article.UpdateDt,
		)
		if e != nil {
			return
		}
		articles = append(articles, article)
	}
	return
}

//-----------------------------------------------

type ArticleController struct {
	helper *ArticleHelper
}

func (ac *ArticleController) fresh(c *gin.Context) {
	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "article", "edit.html")},
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

	_, node_map := share.GenNodeInfo(ac.helper.getAllNode())

	e = tmpl.ExecuteTemplate(c.Writer, "edit.html", gin.H{
		"Title":          "增加文章",
		"UsePicker":      true,
		"NodeMap":        node_map,
		"ArticleTitle":   "請輸入文章標題",
		"ArticleContent": "<p><br></p>",
	})
	if e != nil {
		return
	}
}

func (ac *ArticleController) add(c *gin.Context) {
	const msg = "add"
	logger := utils.NewFileLogger("./dist/logs/article/add", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		ArticleID string `json:"node_id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	const userID = "tester"
	content := share.ProcessBase64Images(userID, param.ArticleID, param.Content)
	e = ac.helper.add(
		param.Title,
		content,
		param.ArticleID,
	)
	if e != nil {
		return
	}

	share.CleanupUnusedImages(userID, param.ArticleID, content)

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (ac *ArticleController) edit(c *gin.Context) {
	const msg = "edit"
	logger := utils.NewFileLogger("./dist/logs/article/add", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()

	var id int
	id, e = strconv.Atoi(c.Param("id"))
	if e != nil {
		return
	}

	list := ac.helper.get(id)
	if len(list) == 0 {
		e = fmt.Errorf("無指定資料")
		return
	}
	data := list[0]

	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "article", "edit.html")},
		Pattern: []string{},
	}

	tmpl, e := utils.RenderTemplate(config)
	e = tmpl.ExecuteTemplate(c.Writer, "edit.html", gin.H{
		"Title":          "增加文章",
		"UsePicker":      false,
		"ArticleTitle":   data.Title,
		"ArticleContent": data.Content,
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
	r.GET("/fresh", c.fresh)
	r.POST("/add", c.add)
	r.GET("/edit/:id", c.edit)
}
