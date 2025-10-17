package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/server/controllers/share"
	"idv/chris/MemoNest/internal/server/middleware"
	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"
)

type NodeHelper struct {
	db *sql.DB
}

func (nh *NodeHelper) addParentNode(pathName string) (*model.Category, error) {
	tx, err := nh.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // 如果有錯誤發生，確保交易回滾

	// 找到最大的 RftIdx，作為新根節點的 LftIdx
	var maxRftIdx int
	err = tx.QueryRow("SELECT COALESCE(MAX(RftIdx), 0) FROM categories").Scan(&maxRftIdx)
	if err != nil {
		return nil, err
	}

	newNodeID := uuid.New().String()
	parentID := uuid.Nil.String() // 根節點的 ParentID
	lftIdx := maxRftIdx + 1
	rftIdx := maxRftIdx + 2

	// 插入新節點
	result, err := tx.Exec(
		"INSERT INTO categories (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, ?, ?, ?, ?)",
		newNodeID, parentID, pathName, lftIdx, rftIdx,
	)
	if err != nil {
		return nil, err
	}

	rowID, _ := result.LastInsertId()
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &model.Category{
		RowID:    int(rowID),
		NodeID:   newNodeID,
		ParentID: parentID,
		PathName: pathName,
		LftIdx:   lftIdx,
		RftIdx:   rftIdx,
	}, nil
}

// addChildNode 插入一個新的分類節點
func (nh *NodeHelper) addChildNode(parentID, pathName string) (*model.Category, error) {
	tx, err := nh.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. 查詢父節點的 RftIdx
	var parentRftIdx int
	err = tx.QueryRow("SELECT RftIdx FROM categories WHERE NodeID = ?", parentID).Scan(&parentRftIdx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("parent node with NodeID '%s' not found", parentID)
		}
		return nil, err
	}

	// 2. 更新所有受影響的節點，為新節點騰出空間
	_, err = tx.Exec("UPDATE categories SET RftIdx = RftIdx + 2 WHERE RftIdx >= ?", parentRftIdx)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE categories SET LftIdx = LftIdx + 2 WHERE LftIdx >= ?", parentRftIdx)
	if err != nil {
		return nil, err
	}

	// 3. 插入新節點
	newNodeID := uuid.New().String()
	lftIdx := parentRftIdx
	rftIdx := parentRftIdx + 1

	result, err := tx.Exec(
		"INSERT INTO categories (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, ?, ?, ?, ?)",
		newNodeID, parentID, pathName, lftIdx, rftIdx,
	)
	if err != nil {
		return nil, err
	}

	rowID, _ := result.LastInsertId()
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &model.Category{
		RowID:    int(rowID),
		NodeID:   newNodeID,
		ParentID: parentID,
		PathName: pathName,
		LftIdx:   lftIdx,
		RftIdx:   rftIdx,
	}, nil
}

// ---------------------------------------------------------
// Nested Set 移除節點的核心邏輯
// ---------------------------------------------------------

// del 移除指定的分類節點及其所有後代節點
func (nh *NodeHelper) del(nodeID string) error {
	row := nh.db.QueryRow(`SELECT COUNT(*) AS Total FROM articles WHERE NodeID = ?;`, nodeID)
	var total int
	e := row.Scan(&total)
	if e != nil {
		return e
	}
	if total != 0 {
		return fmt.Errorf("該節點仍有文章(筆數: %v)", total)
	}

	tx, err := nh.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 查詢要刪除節點的 LftIdx 和 RftIdx
	var lftIdx, rftIdx int
	err = tx.QueryRow("SELECT LftIdx, RftIdx FROM categories WHERE NodeID = ?", nodeID).Scan(&lftIdx, &rftIdx)
	if err != nil {
		// 節點不存在，視為成功
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	// 刪除範圍的寬度
	width := rftIdx - lftIdx + 1

	// 2. 刪除節點及其所有後代節點 (LftIdx 介於 LftIdx 和 RftIdx 之間的)
	_, err = tx.Exec("DELETE FROM categories WHERE LftIdx >= ? AND RftIdx <= ?", lftIdx, rftIdx)
	if err != nil {
		return err
	}

	// 3. 調整剩餘節點的索引，填補被刪除節點留下的空隙
	// 將所有 RftIdx > rftIdx 的右索引減去 width
	_, err = tx.Exec("UPDATE categories SET RftIdx = RftIdx - ? WHERE RftIdx > ?", width, rftIdx)
	if err != nil {
		return err
	}

	// 將所有 LftIdx > rftIdx 的左索引減去 width
	_, err = tx.Exec("UPDATE categories SET LftIdx = LftIdx - ? WHERE LftIdx > ?", width, rftIdx)
	if err != nil {
		return err
	}

	// 提交事務
	return tx.Commit()
}

func (nh *NodeHelper) getAll() (categories []model.Category) {
	// 從資料庫中讀取所有分類，並按 LftIdx 排序
	rows, err := nh.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return
		}
		categories = append(categories, c)
	}

	return
}

//-----------------------------------------------

type NodeController struct {
	helper *NodeHelper
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

	helper := share.NewSessionHelper(c)
	aes_key := []byte(helper.GetAESKey())
	node_id, _ := utils.AesDecrypt(param.ID, aes_key)

	if node_id == uuid.Nil.String() {
		_, e = u.helper.addParentNode(param.Label)
	} else {
		_, e = u.helper.addChildNode(node_id, param.Label)
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

	helper := share.NewSessionHelper(c)
	aes_key := []byte(helper.GetAESKey())
	node_id, _ := utils.AesDecrypt(param.ID, aes_key)

	e = u.helper.del(node_id)
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": fmt.Sprintf("刪除 %s 成功", param.ID)})
}

func (tc *NodeController) list(c *gin.Context) {
	dir := filepath.Join("./assets", "templates")

	helper := share.NewSessionHelper(c)
	aes_key := []byte(helper.GetAESKey())

	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "node", "list.html")},
		Pattern: []string{},
		Funcs: map[string]any{
			"gen_id": func(id int) string {
				cipher_text, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", id)), aes_key)
				return string(cipher_text)
			},
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

	node_list, node_map := share.GenNodeInfo(tc.helper.getAll())
	menu_list, menu_map := share.GetMenu()
	e = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title":    "節點清單",
		"Menu":     menu_list,
		"Children": menu_list[menu_map[share.MK_NODE]].Children,
		"NodeMap":  node_map,
		"List":     node_list,
		"RootID":   uuid.Nil.String(),
	})
	if e != nil {
		return
	}
}

//-----------------------------------------------

func NewNodeController(rg *gin.RouterGroup, di service.DI) {
	c := &NodeController{
		helper: &NodeHelper{
			db: di.MariaDB,
		},
	}
	r := rg.Group("/node")
	r.Use(middleware.MustLoginMiddleware(di))
	r.GET("/list", c.list)
	r.POST("/add", c.add)
	r.POST("/del", c.del)
}
