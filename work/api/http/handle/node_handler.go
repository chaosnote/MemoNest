package handle

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/utils"
)

type NodeHandler struct {
	CommonHandler

	UC *usecase.NodeUsecase
}

func (h *NodeHandler) Add(c *gin.Context) {
	const msg = "add"

	var e error
	defer func() {
		if e != nil {
			h.Log.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	h.Log.Info(msg, zap.Any("params", param))

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	parent_id, e := utils.AesDecrypt(param.ID, aes_key)
	if e != nil {
		return
	}

	e = h.UC.Add(h.Session.GetAccount(), parent_id, "", param.Label)
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": fmt.Sprintf("增加 %s 成功", param.Label)})
}

func (h *NodeHandler) Del(c *gin.Context) {
	const msg = "del"

	var e error
	defer func() {
		if e != nil {
			h.Log.Error(msg, zap.Error(e))
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
	h.Log.Info(msg, zap.Any("params", param))

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	node_id, e := utils.AesDecrypt(param.ID, aes_key)
	if e != nil {
		return
	}

	e = h.UC.Delete(h.Session.GetAccount(), node_id)
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": fmt.Sprintf("刪除 %s 成功", param.ID)})
}

func (h *NodeHandler) List(c *gin.Context) {
	const msg = "list"

	var err error
	defer func() {
		if err != nil {
			h.Log.Error(msg, zap.Error(err))
			h.HandlePageException(c, err.Error())
		}
	}()

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	mo, err := h.UC.GetViewModel(h.Session.GetAccount(), aes_key)
	if err != nil {
		return
	}
	mo.Url = h.Session.GetURL()

	dir := filepath.Join(template_dir)
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "node", "list.html")},
		Pattern: []string{},
		Funcs: map[string]any{
			"encrypt": func(id string) string {
				cipher_text, _ := utils.AesEncrypt([]byte(id), aes_key)
				return string(cipher_text)
			},
		},
	}

	tmpl, err := utils.RenderTemplate(config)
	if err != nil {
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title":   "節點清單",
		"Share":   mo.LayoutContext,
		"NodeMap": mo.NodeMap,
		"List":    mo.NodeList,
		"RootID":  uuid.Nil.String(),
	})
	if err != nil {
		return
	}
}

func (h *NodeHandler) Edit(c *gin.Context) {
	const msg = "edit"

	var e error
	defer func() {
		if e != nil {
			h.Log.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	h.Log.Info(msg, zap.Any("params", param))
	if len(param.Label) == 0 {
		e = fmt.Errorf("label 長度為零")
		return
	}

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	node_id, _ := utils.AesDecrypt(param.ID, aes_key)

	e = h.UC.Edit(h.Session.GetAccount(), node_id, param.Label)
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "編輯成功"})
}

func (h *NodeHandler) Move(c *gin.Context) {
	const msg = "move"

	var err error
	defer func() {
		if err != nil {
			h.Log.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()
	var param struct {
		ParentID  string `json:"parent_id"`
		CurrentID string `json:"current_id"`
	}
	err = c.BindJSON(&param)
	if err != nil {
		return
	}
	h.Log.Info(msg, zap.Any("params", param))

	if param.ParentID == param.CurrentID {
		err = fmt.Errorf("搬移節點相同")
		return
	}

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	parent_id, err := utils.AesDecrypt(param.ParentID, aes_key)
	if err != nil {
		return
	}
	current_id, err := utils.AesDecrypt(param.CurrentID, aes_key)
	if err != nil {
		return
	}

	err = h.UC.Move(h.Session.GetAccount(), parent_id, current_id)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "編輯成功"})
}
