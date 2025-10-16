package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/server/controllers/share"
	"idv/chris/MemoNest/internal/server/middleware"
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

func (nh *ArticleHelper) add(node_id string) (id int, e error) {
	t := time.Now().UTC()
	row := nh.db.QueryRow("CALL insert_article(?, ?, ?, ?, ?) ;", "", "", t, t, node_id)
	e = row.Scan(&id)
	if e != nil {
		return
	}
	return
}

func (nh *ArticleHelper) del(id int) (e error) {
	_, e = nh.db.Exec(`DELETE from articles where RowID = ? ;`, id)
	if e != nil {
		return
	}
	return
}

func (nh *ArticleHelper) update(row_id int, title, content string) error {
	t := time.Now().UTC()
	query := `UPDATE articles SET Title = ?, Content = ?, UpdateDt =? WHERE RowID = ?;`
	_, e := nh.db.Exec(query, title, content, t, row_id)
	if e != nil {
		return e
	}
	return nil
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
			&article.RowID,
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

func (nh *ArticleHelper) list() (articles []model.Article) {
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
		ORDER BY a.UpdateDt DESC ;
	`)
	if e != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		e = rows.Scan(
			&article.RowID,
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
		Page:    []string{filepath.Join(dir, "page", "article", "fresh.html")},
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

	e = tmpl.ExecuteTemplate(c.Writer, "fresh.html", gin.H{
		"Title":        "增加文章",
		"NodeMap":      node_map,
		"ArticleTitle": "請輸入文章標題",
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
		Title   string `json:"title"`
		Content string `json:"content"`
		NodeID  string `json:"node_id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))
	var row_id int
	row_id, e = ac.helper.add(param.NodeID)
	if e != nil {
		return
	}

	article_id := fmt.Sprintf("%v", row_id)

	s := sessions.Default(c)
	account := s.Get(model.SK_ACCOUNT).(string)
	content := share.ProcessBase64Images(account, article_id, param.Content)
	e = ac.helper.update(row_id, param.Title, content)
	if e != nil {
		return
	}
	share.CleanupUnusedImages(account, article_id, content)

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (ac *ArticleController) del(c *gin.Context) {
	const msg = "del"
	logger := utils.NewFileLogger("./dist/logs/article/del", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		Txt string `json:"id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	s := sessions.Default(c)
	account := s.Get(model.SK_ACCOUNT).(string)
	key := []byte(s.Get(model.SK_AES_KEY).(string))
	output, _ := utils.AesDecrypt(param.Txt, key)

	var id int
	id, e = strconv.Atoi(output)
	if e != nil {
		return
	}
	e = ac.helper.del(id)
	if e != nil {
		return
	}

	share.DelImageDir(account, fmt.Sprintf("%v", id))

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (ac *ArticleController) edit(c *gin.Context) {
	const msg = "edit"
	logger := utils.NewFileLogger("./dist/logs/article/edit", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()

	s := sessions.Default(c)
	key := []byte(s.Get(model.SK_AES_KEY).(string))
	txt, e := utils.AesDecrypt(c.Param("id"), key)
	if e != nil {
		return
	}

	var id int
	id, e = strconv.Atoi(txt)
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
		Funcs: map[string]any{
			"trans": func(id, data string) string {
				return strings.ReplaceAll(data, model.IMG_ENCRYPT, id)
			},
		},
	}

	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	e = tmpl.ExecuteTemplate(c.Writer, "edit.html", gin.H{
		"Title":          "修改文章",
		"ID":             c.Param("id"),
		"PathName":       data.PathName,
		"ArticleTitle":   data.Title,
		"ArticleContent": data.Content,
	})
	if e != nil {
		return
	}
}

func (ac *ArticleController) renew(c *gin.Context) {
	const msg = "renew"
	logger := utils.NewFileLogger("./dist/logs/article/renew", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		ID      string `json:"id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	s := sessions.Default(c)
	account := s.Get(model.SK_ACCOUNT).(string)
	key := []byte(s.Get(model.SK_AES_KEY).(string))
	article_id, _ := utils.AesDecrypt(param.ID, key)

	var row_id int
	row_id, e = strconv.Atoi(article_id)
	if e != nil {
		return
	}

	content := share.ProcessBase64Images(account, article_id, param.Content)
	e = ac.helper.update(row_id, param.Title, content)
	if e != nil {
		return
	}
	share.CleanupUnusedImages(account, article_id, content)

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (ac *ArticleController) list(c *gin.Context) {
	const msg = "list"
	logger := utils.NewFileLogger("./dist/logs/article/list", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()

	s := sessions.Default(c)
	key := []byte(s.Get(model.SK_AES_KEY).(string))

	list := ac.helper.list()

	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "article", "list.html")},
		Pattern: []string{},
		Funcs: map[string]any{
			"format": func(t time.Time) string {
				loc, _ := time.LoadLocation("Asia/Taipei")
				return t.In(loc).Format("2006-01-02 15:04")
			},
			"encrypt": func(id int) string {
				output, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", id)), key)
				return string(output)
			},
			"trans": func(id int, data string) template.HTML {
				output, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", id)), key)
				return template.HTML(strings.ReplaceAll(data, model.IMG_ENCRYPT, output))
			},
		},
	}
	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	e = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title": "文章清單",
		"List":  list,
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
	r.Use(middleware.MustLoginMiddleware(di))
	r.GET("/fresh", c.fresh)
	r.GET("/list", c.list)
	r.POST("/add", c.add)
	r.POST("/del", c.del)
	r.POST("/renew", c.renew)
	r.GET("/edit/:id", c.edit)
}
