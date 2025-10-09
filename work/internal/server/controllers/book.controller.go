package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

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

func (bh *BookHelper) init(pathName string) (*model.Category, error) {
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

// addCategory 插入一個新的分類節點
func (bh *BookHelper) addCategory(parentID, pathName string) (*model.Category, error) {
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

//-----------------------------------------------

type BookController struct {
	helper *BookHelper
}

// curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"test\"}" http://localhost:8080/api/v1/book/init
func (u *BookController) init() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		_, e = u.helper.init(params["name"].(string))
		if e != nil {
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": "", "message": fmt.Sprintf("增加 %s 成功", params["name"])})
	}
}

//-----------------------------------------------

func NewBookController(rg *gin.RouterGroup, di service.DI) {
	c := &BookController{
		helper: &BookHelper{
			db: di.MariaDB,
		},
	}
	r := rg.Group("/book")
	r.POST("/init", c.init())
}
