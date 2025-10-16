package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"
)

type IndexController struct {
	Debug bool
}

func (ic *IndexController) entry(c *gin.Context) {
	if !ic.Debug {
		s := sessions.Default(c)
		flag, ok := s.Get(model.SK_IS_LOGIN).(bool)
		if !ok || !flag {
			c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "未登入"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "已登入"})
		return
	}

	dir := filepath.Join("./assets", "templates")
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
		"Title":   "測試頁",
		"Login":   []string{"/api/v1/member/login", "/api/v1/member/logout"},
		"Article": []string{"/api/v1/article/fresh", "/api/v1/article/list"},
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
