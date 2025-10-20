package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	xxx "idv/chris/MemoNest/adapter/http"
	"idv/chris/MemoNest/adapter/http/middleware"
	"idv/chris/MemoNest/adapter/infra"
	"idv/chris/MemoNest/domain/repo"
	zzz "idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/service"
	"idv/chris/MemoNest/utils"
)

type NodeController struct {
	repo repo.NodeRepository
	menu zzz.MenuProvider
	tree zzz.NodeTree
}

func (u *NodeController) add(c *gin.Context) {
	const msg = "add"
	logger := utils.NewFileLogger("./dist/logs/node/add", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
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
	logger.Info(msg, zap.Any("params", param))

	helper := xxx.NewGinSession(c)
	aes_key := []byte(helper.GetAESKey())
	parent_id, _ := utils.AesDecrypt(param.ID, aes_key)

	if parent_id == uuid.Nil.String() {
		_, e = u.repo.AddParentNode("", param.Label)
	} else {
		_, e = u.repo.AddChildNode(parent_id, "", param.Label)
	}
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": fmt.Sprintf("增加 %s 成功", param.Label)})
}

func (u *NodeController) del(c *gin.Context) {
	const msg = "del"
	logger := utils.NewFileLogger("./dist/logs/node/del", "console", 1)
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

	helper := xxx.NewGinSession(c)
	aes_key := []byte(helper.GetAESKey())
	node_id, _ := utils.AesDecrypt(param.ID, aes_key)

	e = u.repo.Delete(node_id)
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": fmt.Sprintf("刪除 %s 成功", param.ID)})
}

func (tc *NodeController) list(c *gin.Context) {
	dir := filepath.Join("./web", "templates")

	helper := xxx.NewGinSession(c)
	aes_key := []byte(helper.GetAESKey())

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
	tmpl, e := utils.RenderTemplate(config)
	defer func() {
		if e != nil {
			http.Error(c.Writer, e.Error(), http.StatusInternalServerError)
		}
	}()
	if e != nil {
		return
	}

	tmp_list, e := tc.repo.GetAllNode()
	if e != nil {
		return
	}
	node_list, node_map := tc.tree.GenInfo(tmp_list)
	for _, node := range node_list {
		tc.repo.AssignNode(node, aes_key)
	}

	e = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title":    "節點清單",
		"Menu":     tc.menu.GetList(),
		"Children": tc.menu.GetList()[tc.menu.GetMap()[infra.MP_NODE]].Children,
		"NodeMap":  node_map,
		"List":     node_list,
		"RootID":   uuid.Nil.String(),
	})
	if e != nil {
		return
	}
}

func (u *NodeController) edit(c *gin.Context) {
	const msg = "edit"
	logger := utils.NewFileLogger("./dist/logs/node/edit", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
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
	logger.Info(msg, zap.Any("params", param))
	if len(param.Label) == 0 {
		e = fmt.Errorf("label 長度為零")
		return
	}

	helper := xxx.NewGinSession(c)
	aes_key := []byte(helper.GetAESKey())
	node_id, _ := utils.AesDecrypt(param.ID, aes_key)

	e = u.repo.Edit(node_id, param.Label)
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "編輯成功"})
}

func (u *NodeController) move(c *gin.Context) {
	const msg = "move"
	logger := utils.NewFileLogger("./dist/logs/node/edit", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		ParentID  string `json:"parent_id"`
		CurrentID string `json:"current_id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	helper := xxx.NewGinSession(c)
	aes_key := []byte(helper.GetAESKey())
	parent_id, _ := utils.AesDecrypt(param.ParentID, aes_key)
	current_id, _ := utils.AesDecrypt(param.CurrentID, aes_key)

	var has_node = false
	if parent_id == uuid.Nil.String() {
		has_node = true
	} else {
		parent_node, e := u.repo.GetNode(parent_id)
		if e != nil {
			return
		}
		if parent_node.RowID != 0 {
			has_node = true
		}
	}
	if !has_node {
		e = fmt.Errorf("無指定父節點")
		return
	}

	current_node, e := u.repo.GetNode(current_id)
	if e != nil {
		return
	}
	if current_node.RowID != 0 {
		has_node = true
	}
	if !has_node {
		e = fmt.Errorf("無指定子節點")
		return
	}

	e = u.repo.Move(parent_id, current_id, current_node.PathName)
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "編輯成功"})
}

//-----------------------------------------------

func NewNodeController(rg *gin.RouterGroup, di service.DI, repo repo.NodeRepository, menu zzz.MenuProvider, tree zzz.NodeTree) {
	c := &NodeController{
		repo: repo,
		menu: menu,
		tree: tree,
	}
	r := rg.Group("/node")
	r.Use(middleware.Auth(di))
	r.GET("/list", c.list)
	r.POST("/add", c.add)
	r.POST("/del", c.del)
	r.POST("/edit", c.edit)
	r.POST("/move", c.move)
}
