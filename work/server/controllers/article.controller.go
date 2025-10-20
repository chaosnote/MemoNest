package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/model"
	"idv/chris/MemoNest/server/controllers/share"
	"idv/chris/MemoNest/server/middleware"
	"idv/chris/MemoNest/service"
	"idv/chris/MemoNest/utils"
)

type ArticleController struct {
	repo repo.ArticleRepository
}

func (ac *ArticleController) fresh(c *gin.Context) {
	dir := filepath.Join("./web", "templates")
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

	tmp_list, e := ac.repo.GetAllNode()
	if e != nil {
		return
	}
	_, node_map := share.GenNodeInfo(tmp_list)
	menu_list, menu_map := share.GetMenu()
	e = tmpl.ExecuteTemplate(c.Writer, "fresh.html", gin.H{
		"Title":        "增加文章",
		"Menu":         menu_list,
		"Children":     menu_list[menu_map[share.MK_ARTICLE]].Children,
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
	row_id, e = ac.repo.Add(param.NodeID)
	if e != nil {
		return
	}

	article_id := fmt.Sprintf("%v", row_id)

	helper := share.NewSessionHelper(c)
	account := helper.GetAccount()
	content := share.ProcessBase64Images(account, article_id, param.Content)
	e = ac.repo.Update(row_id, param.Title, content)
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
		ID string `json:"id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	helper := share.NewSessionHelper(c)
	account := helper.GetAccount()
	aes_key := []byte(helper.GetAESKey())
	plain_text, _ := utils.AesDecrypt(param.ID, aes_key)

	var id int
	id, e = strconv.Atoi(plain_text)
	if e != nil {
		return
	}
	e = ac.repo.Delete(id)
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

	helper := share.NewSessionHelper(c)
	aes_key := []byte(helper.GetAESKey())
	plain_text, e := utils.AesDecrypt(c.Param("id"), aes_key)
	if e != nil {
		return
	}

	var id int
	id, e = strconv.Atoi(plain_text)
	if e != nil {
		return
	}

	list, e := ac.repo.Get(id)
	if e != nil {
		return
	}
	if len(list) == 0 {
		e = fmt.Errorf("無指定資料")
		return
	}
	data := list[0]

	dir := filepath.Join("./web", "templates")
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
	menu_list, menu_map := share.GetMenu()
	e = tmpl.ExecuteTemplate(c.Writer, "edit.html", gin.H{
		"Title":          "修改文章",
		"Menu":           menu_list,
		"Children":       menu_list[menu_map[share.MK_ARTICLE]].Children,
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

	helper := share.NewSessionHelper(c)
	account := helper.GetAccount()
	aes_key := []byte(helper.GetAESKey())
	article_id, _ := utils.AesDecrypt(param.ID, aes_key)

	var row_id int
	row_id, e = strconv.Atoi(article_id)
	if e != nil {
		return
	}

	content := share.ProcessBase64Images(account, article_id, param.Content)
	e = ac.repo.Update(row_id, param.Title, content)
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

	q := c.Query("q")

	var list []model.Article
	if len(q) > 0 {
		list, e = ac.repo.Query(q)
	} else {
		list, e = ac.repo.List()
	}
	if e != nil {
		return
	}

	helper := share.NewSessionHelper(c)
	aes_key := []byte(helper.GetAESKey())

	dir := filepath.Join("./web", "templates")
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
				cipher_text, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", id)), aes_key)
				return string(cipher_text)
			},
			"trans": func(id int, data string) template.HTML {
				output, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", id)), aes_key)
				return template.HTML(strings.ReplaceAll(data, model.IMG_ENCRYPT, output))
			},
		},
	}
	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	menu_list, menu_map := share.GetMenu()
	e = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title":    "文章清單",
		"Menu":     menu_list,
		"Children": menu_list[menu_map[share.MK_ARTICLE]].Children,
		"List":     list,
	})
	if e != nil {
		return
	}
}

func NewArticleController(rg *gin.RouterGroup, di service.DI, repo repo.ArticleRepository) {
	c := &ArticleController{
		repo: repo,
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
