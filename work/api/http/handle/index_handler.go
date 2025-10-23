package handle

import (
	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IndexHandler struct {
	Debug   bool
	UC      *usecase.IndexUsecase
	Session service.Session
}

func (h *IndexHandler) Entry(c *gin.Context) {
	const msg = "entry"
	logger := utils.NewFileLogger("./dist/logs/index/entry", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()

	h.Session.Init(c)
	dir := filepath.Join("./web", "templates")
	if h.Session.IsLogin() {
		config := utils.TemplateConfig{
			Layout:  filepath.Join(dir, "layout", "share.html"),
			Page:    []string{filepath.Join(dir, "page", "index", "logged_in.html")},
			Pattern: []string{},
		}
		tmpl, e := utils.RenderTemplate(config)
		if e != nil {
			return
		}
		mo := h.UC.GetViewModel("chris", "123456", "")
		e = tmpl.ExecuteTemplate(c.Writer, "logged_in.html", gin.H{
			"Title":    "首頁",
			"Menu":     mo.Menu,
			"Children": mo.MenuChildren,
		})
		if e != nil {
			return
		}
	} else {
		config := utils.TemplateConfig{
			Layout:  filepath.Join(dir, "layout", "share.html"),
			Page:    []string{filepath.Join(dir, "page", "index", "logged_out.html")},
			Pattern: []string{},
		}
		tmpl, e := utils.RenderTemplate(config)
		if e != nil {
			return
		}
		mo := h.UC.GetViewModel("", "", "")
		e = tmpl.ExecuteTemplate(c.Writer, "logged_out.html", gin.H{
			"Title":   "首頁",
			"Setting": mo,
		})
		if e != nil {
			return
		}
	}
}

func (h *IndexHandler) Register(c *gin.Context) {
	const msg = "register"
	logger := utils.NewFileLogger("./dist/logs/index/register", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()

	dir := filepath.Join("./web", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "index", "register.html")},
		Pattern: []string{},
	}
	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	e = tmpl.ExecuteTemplate(c.Writer, "register.html", gin.H{
		"Title": "首頁",
	})
	if e != nil {
		return
	}
}

func (h *IndexHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}
