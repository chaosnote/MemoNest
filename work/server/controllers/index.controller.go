package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	xxx "idv/chris/MemoNest/adapter/http"
	"idv/chris/MemoNest/service"
	"idv/chris/MemoNest/utils"
)

type IndexController struct {
	Debug bool
}

func (ic *IndexController) entry(c *gin.Context) {
	if !ic.Debug {
		helper := xxx.NewGinSession(c)
		if helper.IsLogin() {
			c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "未登入"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "已登入"})
		return
	}

	dir := filepath.Join("./web", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "index", "index.html")},
		Pattern: []string{},
	}

	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	e = tmpl.ExecuteTemplate(c.Writer, "index.html", gin.H{
		"Title": "測試頁",
		"Login": []string{"/api/v1/member/login", "/api/v1/member/logout"},
	})
	if e != nil {
		return
	}
}

func (ic *IndexController) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "OK"})
}

func NewIndexController(engine *gin.Engine, di service.DI) {
	c := &IndexController{
		Debug: di.Config.Gin.Mode == "debug",
	}

	engine.GET("/", c.entry)
	engine.GET("/health", c.health)
}
