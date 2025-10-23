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
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
)

type ArticleHandler struct {
	UC      *usecase.ArticleUsecase
	Session service.Session
}

func (h *ArticleHandler) Fresh(c *gin.Context) {
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

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())

	mo, e := h.UC.GetViewModel(h.Session.GetAccount(), aes_key)
	if e != nil {
		return
	}

	e = tmpl.ExecuteTemplate(c.Writer, "fresh.html", gin.H{
		"Title":        "增加文章",
		"Share":        mo.LayoutShare,
		"NodeMap":      mo.NodeMap,
		"ArticleTitle": "請輸入文章標題",
	})
	if e != nil {
		return
	}
}

func (h *ArticleHandler) Add(c *gin.Context) {
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

	h.Session.Init(c)
	account := h.Session.GetAccount()

	e = h.UC.Add(account, param.NodeID, param.Title, param.Content)
	if e != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (h *ArticleHandler) Del(c *gin.Context) {
	const msg = "del"
	logger := utils.NewFileLogger("./dist/logs/article/del", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()
	var param struct {
		ID string `json:"id"`
	}
	err = c.BindJSON(&param)
	if err != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	h.Session.Init(c)
	account := h.Session.GetAccount()
	aes_key := []byte(h.Session.GetAESKey())
	plain_text, err := utils.AesDecrypt(param.ID, aes_key)

	h.UC.Del(account, plain_text)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (h *ArticleHandler) Edit(c *gin.Context) {
	const msg = "edit"
	logger := utils.NewFileLogger("./dist/logs/article/edit", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()

	id := c.Param("id")

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	plain_text, err := utils.AesDecrypt(id, aes_key)
	if err != nil {
		return
	}
	data, err := h.UC.Edit(h.Session.GetAccount(), plain_text)
	if err != nil {
		return
	}

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
	tmpl, err := utils.RenderTemplate(config)
	if err != nil {
		return
	}
	mo, err := h.UC.GetViewModel(h.Session.GetAccount(), aes_key)
	err = tmpl.ExecuteTemplate(c.Writer, "edit.html", gin.H{
		"Title":          "修改文章",
		"Share":          mo.LayoutShare,
		"ID":             id,
		"PathName":       data.PathName,
		"ArticleTitle":   data.Title,
		"ArticleContent": data.Content,
	})
	if err != nil {
		return
	}
}

func (h *ArticleHandler) Renew(c *gin.Context) {
	const msg = "renew"
	logger := utils.NewFileLogger("./dist/logs/article/renew", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()
	var param struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		ID      string `json:"id"`
	}
	err = c.BindJSON(&param)
	if err != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	account := h.Session.GetAccount()
	article_id, err := utils.AesDecrypt(param.ID, aes_key)
	if err != nil {
		return
	}

	err = h.UC.Renew(account, article_id, param.Title, param.Content)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (h *ArticleHandler) List(c *gin.Context) {
	const msg = "list"
	logger := utils.NewFileLogger("./dist/logs/article/list", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()

	Q := c.Query("q")
	list, err := h.UC.List(h.Session.GetAccount(), Q)
	if err != nil {
		return
	}

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())

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
	tmpl, err := utils.RenderTemplate(config)
	if err != nil {
		return
	}

	mo, err := h.UC.GetViewModel(h.Session.GetAccount(), aes_key)
	if err != nil {
		return
	}
	err = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title": "文章清單",
		"Share": mo.LayoutShare,
		"Q":     Q,
		"List":  list,
	})
	if err != nil {
		return
	}
}
