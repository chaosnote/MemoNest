package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"
)

type BookHelper struct {
	db *sql.DB
}

func (bh *BookHelper) addParentNode(pathName string) (*model.Category, error) {
	tx, err := bh.db.Begin()
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
func (bh *BookHelper) addChildNode(parentID, pathName string) (*model.Category, error) {
	tx, err := bh.db.Begin()
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

// removeNode 移除指定的分類節點及其所有後代節點
func (bh *BookHelper) removeNode(nodeID string) error {
	tx, err := bh.db.Begin()
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

// removeParentNode 移除一個根節點 (或任何頂層節點)
func (bh *BookHelper) removeParentNode(nodeID string) error {
	return bh.removeNode(nodeID)
}

// removeChildNode 移除一個子節點
func (bh *BookHelper) removeChildNode(nodeID string) error {
	return bh.removeNode(nodeID)
}
func (th *BookHelper) getAllNode() (categories []model.Category) {
	// 從資料庫中讀取所有分類，並按 LftIdx 排序
	rows, err := th.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
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

type BookController struct {
	helper *BookHelper
}

// curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"test\"}" http://localhost:8080/api/v1/book/addParentNode
func (u *BookController) addParentNode(c *gin.Context) {
	logger := utils.NewFileLogger("./dist/book", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error("init", zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"error": e.Error()})
		}
	}()
	var params map[string]interface{}
	e = c.BindJSON(&params)
	if e != nil {
		return
	}
	logger.Info("init", zap.Any("params", params))
	_, e = u.helper.addParentNode(params["name"].(string))
	if e != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "", "message": fmt.Sprintf("增加 %s 成功", params["name"])})
}

// 顯示節點路徑
func (tc *BookController) getAllNode(c *gin.Context) {
	parent_id := uuid.Nil.String()
	allCategories := tc.helper.getAllNode()

	nodesMap := make(map[string]*model.CategoryNode)
	var rootNodes []*model.CategoryNode

	// 第一次遍歷：建立節點地圖
	for _, cat := range allCategories {
		nodesMap[cat.NodeID] = &model.CategoryNode{
			Category: cat,
		}
	}

	// 第二次遍歷：建立樹狀結構並生成路徑
	for _, cat := range allCategories {
		currentNode := nodesMap[cat.NodeID]

		// 建立完整路徑
		pathParts := []string{currentNode.PathName}
		tempNode := currentNode
		for {
			if tempNode.ParentID == parent_id {
				break
			}
			if parent, ok := nodesMap[tempNode.ParentID]; ok {
				pathParts = append([]string{parent.PathName}, pathParts...) // 將父節點名稱加到最前面
				tempNode = parent
			} else {
				break // 父節點不存在，停止回溯
			}
		}
		currentNode.Path = "/" + strings.Join(pathParts, "/")

		// 處理樹狀結構
		if cat.ParentID == parent_id {
			rootNodes = append(rootNodes, currentNode)
		} else {
			if parent, ok := nodesMap[cat.ParentID]; ok {
				parent.Children = append(parent.Children, currentNode)
			}
		}
	}

	data := gin.H{"Categories": rootNodes}

	templates := template.Must(template.ParseFiles(filepath.Join("./assets", "templates", "page", "book", "list.html")))
	err := templates.Execute(c.Writer, data)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (tc *BookController) tmp(c *gin.Context) {
	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout: filepath.Join(dir, "layout", "default.html"),
		Page:   []string{filepath.Join(dir, "page", "book", "tmp.html")},
	}
	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		fmt.Println(e)
		http.Error(c.Writer, e.Error(), http.StatusInternalServerError)
		return
	}
	e = tmpl.Execute(c.Writer, gin.H{"Title": "XXXX"})
	if e != nil {
		fmt.Println(e)
		http.Error(c.Writer, e.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("end")
}

//-----------------------------------------------

func NewBookController(rg *gin.RouterGroup, di service.DI) {
	c := &BookController{
		helper: &BookHelper{
			db: di.MariaDB,
		},
	}
	r := rg.Group("/book")
	r.GET("/list", c.getAllNode)
	r.GET("/tmp", c.tmp)
	r.POST("/init", c.addParentNode)

}
