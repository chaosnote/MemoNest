package service

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
)

// MariaDBImpl MariaDB 客戶端結構
type MariaDBImpl struct {
	db     *sql.DB
	logger *zap.Logger
}

func (mds *MariaDBImpl) AddRootNode(pathName string) (*model.Category, error) {
	tx, err := mds.db.Begin()
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

// AddCategory 插入一個新的分類節點
func (mds *MariaDBImpl) AddCategory(parentID, pathName string) (*model.Category, error) {
	tx, err := mds.db.Begin()
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

func (mds *MariaDBImpl) GetCategories() []model.Category {
	var categories []model.Category

	// 從資料庫中讀取所有分類，並按 LftIdx 排序
	rows, err := mds.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return categories
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return categories
		}
		categories = append(categories, c)
	}

	return categories
}

func (mds *MariaDBImpl) Close() error {
	return mds.db.Close()
}

// NewMariaDBImpl 建立 MariaDB 連線
func NewMariaDBImpl(cfg *model.APPConfig, logger *zap.Logger) (*MariaDBImpl, error) {
	db, e := sql.Open("mysql", cfg.Mariadb.DSN)
	if e != nil {
		return nil, e
	}

	if e = db.Ping(); e != nil {
		return nil, e
	}

	return &MariaDBImpl{db: db, logger: logger}, nil
}
