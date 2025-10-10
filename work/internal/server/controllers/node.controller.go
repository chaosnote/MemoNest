package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
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
	tx, err := nh.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 查詢要刪除節點的 LftIdx 和 RftIdx
	var lftIdx, rftIdx int
	err = tx.QueryRow("SELECT LftIdx, RftIdx FROM categories WHERE NodeID = ?", nodeID).Scan(&lftIdx, &rftIdx)
	if err != nil {
		if err == sql.ErrNoRows {
			// 節點不存在，視為成功
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
	logger := utils.NewFileLogger("./dist/node/add", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var params map[string]string
	e = c.BindJSON(&params)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", params))

	const (
		id    = "id"
		label = "label"
	)
	if params["id"] == uuid.Nil.String() {
		_, e = u.helper.addParentNode(params[label])
	} else {
		_, e = u.helper.addChildNode(params[id], params[label])
	}
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": fmt.Sprintf("增加 %s 成功", params[label])})
}

func (u *NodeController) del(c *gin.Context) {
	const msg = "del"
	logger := utils.NewFileLogger("./dist/node/del", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var params map[string]string
	e = c.BindJSON(&params)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", params))

	e = u.helper.del(params["id"])
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": fmt.Sprintf("刪除 %s 成功", params["id"])})
}

func (tc *NodeController) list(c *gin.Context) {
	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "node", "list.html")},
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

	root_id := uuid.Nil.String()
	source := tc.helper.getAll()
	node_map := make(map[string]*model.CategoryNode)

	root_node := []*model.CategoryNode{}
	// 第一次遍歷：建立節點地圖
	for _, cat := range source {
		node_map[cat.NodeID] = &model.CategoryNode{
			Category: cat,
		}
	}
	// 第二次遍歷：建立樹狀結構並生成路徑
	for _, cat := range source {
		current_node := node_map[cat.NodeID]

		// 建立完整路徑
		path_seg := []string{current_node.PathName}
		temp_node := current_node
		for {
			if temp_node.ParentID == root_id {
				break
			}
			parent, ok := node_map[temp_node.ParentID]
			if !ok {
				break
			}
			path_seg = append([]string{parent.PathName}, path_seg...) // 將父節點名稱加到最前面
			temp_node = parent
		}
		current_node.Path = strings.Join(path_seg, "/")

		// 處理樹狀結構
		if cat.ParentID == root_id {
			root_node = append(root_node, current_node)
		} else {
			if parent, ok := node_map[cat.ParentID]; ok {
				parent.Children = append(parent.Children, current_node)
			}
		}
	}

	e = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title":   "節點清單",
		"NodeMap": node_map,
		"List":    root_node,
		"RootID":  root_id,
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
	r.GET("/list", c.list)
	r.POST("/add", c.add)
	r.POST("/del", c.del)
}
