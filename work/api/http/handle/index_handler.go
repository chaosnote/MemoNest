package handle

import (
	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type IndexHandler struct {
	Debug   bool
	UC      *usecase.IndexUsecase
	Session service.Session
}

func (h *IndexHandler) Entry(c *gin.Context) {
	if !h.Debug {
		h.Session.Init(c)
		if h.Session.IsLogin() {
			c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "未登入"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "已登入"})
		return
	}

	dir := filepath.Join("./web", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "index", "logged_in.html")},
		Pattern: []string{},
	}

	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	e = tmpl.ExecuteTemplate(c.Writer, "logged_in.html", gin.H{
		"Title": "測試頁",
		"Login": []string{"/api/v1/member/login", "/api/v1/member/logout"},
	})
	if e != nil {
		return
	}
}

func (h *IndexHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}
