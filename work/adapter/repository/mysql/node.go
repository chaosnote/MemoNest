package mysql

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/model"
	"idv/chris/MemoNest/utils"
)

type NodeRepo struct {
	db *sql.DB
}

func (nh *NodeRepo) AddParentNode(node_id, path_name string) (*model.Category, error) {
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

	if len(node_id) == 0 {
		node_id = uuid.New().String()
	}
	parent_id := uuid.Nil.String() // 根節點的 ParentID
	lftIdx := maxRftIdx + 1
	rftIdx := maxRftIdx + 2

	// 插入新節點
	result, err := tx.Exec(
		"INSERT INTO categories (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, ?, ?, ?, ?)",
		node_id, parent_id, path_name, lftIdx, rftIdx,
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
		NodeID:   node_id,
		ParentID: parent_id,
		PathName: path_name,
		LftIdx:   lftIdx,
		RftIdx:   rftIdx,
	}, nil
}

// AddChildNode 插入一個新的分類節點
func (nh *NodeRepo) AddChildNode(parent_id, node_id, path_name string) (*model.Category, error) {
	tx, err := nh.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. 查詢父節點的 RftIdx
	var parentRftIdx int
	err = tx.QueryRow("SELECT RftIdx FROM categories WHERE NodeID = ?", parent_id).Scan(&parentRftIdx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("parent node with NodeID '%s' not found", parent_id)
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
	if len(node_id) == 0 {
		node_id = uuid.New().String()
	}
	lftIdx := parentRftIdx
	rftIdx := parentRftIdx + 1

	result, err := tx.Exec(
		"INSERT INTO categories (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, ?, ?, ?, ?)",
		node_id, parent_id, path_name, lftIdx, rftIdx,
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
		NodeID:   node_id,
		ParentID: parent_id,
		PathName: path_name,
		LftIdx:   lftIdx,
		RftIdx:   rftIdx,
	}, nil
}

// Delete 移除指定的分類節點及其所有後代節點
func (nh *NodeRepo) Delete(node_id string) error {
	row := nh.db.QueryRow(`SELECT COUNT(*) AS Total FROM articles WHERE NodeID = ?;`, node_id)
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
	err = tx.QueryRow("SELECT LftIdx, RftIdx FROM categories WHERE NodeID = ?", node_id).Scan(&lftIdx, &rftIdx)
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

// Edit 編輯
func (nh *NodeRepo) Edit(node_id, label string) error {
	_, err := nh.db.Exec(`UPDATE categories SET PathName = ? WHERE NodeID = ?;`, label, node_id)
	if err != nil {
		return err
	}
	return nil
}

func (nh *NodeRepo) Move(parent_id, node_id, path_name string) error {
	e := nh.Delete(node_id)
	if e != nil {
		return e
	}

	if parent_id == uuid.Nil.String() {
		_, e = nh.AddParentNode(node_id, path_name)
		if e != nil {
			return e
		}
	} else {
		_, e = nh.AddChildNode(parent_id, node_id, path_name)
		if e != nil {
			return e
		}
	}

	return nil
}

func (nh *NodeRepo) GetAllNode() (categories []model.Category, err error) {
	// 從資料庫中讀取所有分類，並按 LftIdx 排序
	rows, err := nh.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}

	return
}

func (nh *NodeRepo) GetNode(node_id string) (c model.Category, e error) {
	row := nh.db.QueryRow("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories WHERE NodeID = ?", node_id)
	if err := row.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
		return
	}
	return
}

func (nh *NodeRepo) AssignNode(node *model.CategoryNode, aes_key []byte) {
	sUID, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", node.RowID)), aes_key)
	node.El_UID = sUID

	sNodeID, _ := utils.AesEncrypt([]byte(node.NodeID), aes_key)
	node.El_NodeID = sNodeID

	for _, child := range node.Children {
		nh.AssignNode(child, aes_key)
	}
}

func NewNodeRepo(db *sql.DB) repo.NodeRepository {
	return &NodeRepo{db: db}
}
