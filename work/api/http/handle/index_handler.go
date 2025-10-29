package handle

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/utils"
)

type IndexHandler struct {
	CommonHandler

	Debug bool
	UC    *usecase.IndexUsecase
}

func (h *IndexHandler) Entry(c *gin.Context) {
	const msg = "entry"

	var err error
	defer func() {
		if err != nil {
			h.Log.Error(msg, zap.Error(err))
			h.HandlePageException(c, err.Error())
		}
	}()

	var mo model.IndexView

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())

	dir := filepath.Join(template_dir)
	if h.Session.IsLogin() {
		config := utils.TemplateConfig{
			Layout:  filepath.Join(dir, "layout", "share.html"),
			Page:    []string{filepath.Join(dir, "page", "index", "logged_in.html")},
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

		list, err := h.UC.List(h.Session.GetAccount())
		if err != nil {
			return
		}
		mo = h.UC.GetViewModel("", "")

		e = tmpl.ExecuteTemplate(c.Writer, "logged_in.html", gin.H{
			"Title": "首頁",
			"Share": mo.LayoutShare,
			"List":  list,
		})
		if e != nil {
			return
		}
	} else {
		if h.Debug {
			mo = h.UC.GetViewModel("chris", "123456")
		} else {
			mo = h.UC.GetViewModel("", "")
		}

		config := utils.TemplateConfig{
			Layout:  filepath.Join(dir, "layout", "share.html"),
			Page:    []string{filepath.Join(dir, "page", "index", "logged_out.html")},
			Pattern: []string{},
		}
		tmpl, e := utils.RenderTemplate(config)
		if e != nil {
			return
		}

		e = tmpl.ExecuteTemplate(c.Writer, "logged_out.html", gin.H{
			"Title":    "首頁",
			"Account":  mo.Account,
			"Password": mo.Password,
		})
		if e != nil {
			return
		}
	}
}

func (h *IndexHandler) Register(c *gin.Context) {
	const msg = "register"

	var err error
	defer func() {
		if err != nil {
			h.Log.Error(msg, zap.Error(err))
			h.HandlePageException(c, err.Error())
		}
	}()

	dir := filepath.Join(template_dir)
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

func (h *IndexHandler) User(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}
