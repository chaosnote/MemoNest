package handle

import (
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const template_dir = "./web/templates"

type CommonHandler struct {
	Log     *zap.Logger
	Session service.Session
}

func (h *CommonHandler) PageException(c *gin.Context, message string) {
	h.Session.Init(c)
	h.Session.Clear()

	dir := filepath.Join(template_dir)
	config := utils.TemplateConfig{
		Layout:  "",
		Page:    []string{filepath.Join(dir, "error", "exception.html")},
		Pattern: []string{},
	}

	const msg = "page_exception"

	var err error
	defer func() {
		if err != nil {
			h.Log.Error(msg, zap.Error(err))
		}
	}()

	tmpl, err := utils.RenderTemplate(config)
	if err != nil {
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "exception.html", gin.H{
		"Code": message,
	})
	if err != nil {
		return
	}
}
