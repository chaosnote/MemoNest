package mysql

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
)

type NodeRepo struct {
	db                 *sql.DB
	node_formatter     string
	articles_formatter string
}

func (r *NodeRepo) AddParentNode(account, node_id, path_name string) (*entity.Node, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 找到最大的 RftIdx，作為新根節點的 LftIdx
	var maxRftIdx int
	query := fmt.Sprintf(
		"SELECT COALESCE(MAX(RftIdx), 0) FROM %s",
		fmt.Sprintf(r.node_formatter, account),
	)
	err = tx.QueryRow(query).Scan(&maxRftIdx)
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
	query = fmt.Sprintf(
		"INSERT INTO %s (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, ?, ?, ?, ?)",
		fmt.Sprintf(r.node_formatter, account),
	)
	result, err := tx.Exec(
		query,
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
	committed = true
	return &entity.Node{
		RowID:    int(rowID),
		NodeID:   node_id,
		ParentID: parent_id,
		PathName: path_name,
		LftIdx:   lftIdx,
		RftIdx:   rftIdx,
	}, nil
}

// AddChildNode 插入一個新的分類節點
func (r *NodeRepo) AddChildNode(account, parent_id, node_id, path_name string) (*entity.Node, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 1. 查詢父節點的 RftIdx
	query := fmt.Sprintf(
		"SELECT RftIdx FROM %s WHERE NodeID = ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	var parentRftIdx int
	err = tx.QueryRow(query, parent_id).Scan(&parentRftIdx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("parent node with NodeID '%s' not found[ERR]無指定父節點", parent_id)
		}
		return nil, err
	}

	// 2. 更新所有受影響的節點，為新節點騰出空間
	query = fmt.Sprintf(
		"UPDATE %s SET RftIdx = RftIdx + 2 WHERE RftIdx >= ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	_, err = tx.Exec(query, parentRftIdx)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf(
		"UPDATE %s SET LftIdx = LftIdx + 2 WHERE LftIdx >= ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	_, err = tx.Exec(query, parentRftIdx)
	if err != nil {
		return nil, err
	}

	// 3. 插入新節點
	if len(node_id) == 0 {
		node_id = uuid.New().String()
	}
	lftIdx := parentRftIdx
	rftIdx := parentRftIdx + 1

	query = fmt.Sprintf(
		"INSERT INTO %s (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, ?, ?, ?, ?)",
		fmt.Sprintf(r.node_formatter, account),
	)
	result, err := tx.Exec(
		query,
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
	committed = true
	return &entity.Node{
		RowID:    int(rowID),
		NodeID:   node_id,
		ParentID: parent_id,
		PathName: path_name,
		LftIdx:   lftIdx,
		RftIdx:   rftIdx,
	}, nil
}

// Delete 移除指定的分類節點及其所有後代節點
func (r *NodeRepo) Delete(account, node_id string) error {
	query := fmt.Sprintf(
		"SELECT COUNT(*) AS Total FROM %s WHERE NodeID = ?;",
		fmt.Sprintf(r.articles_formatter, account),
	)
	row := r.db.QueryRow(query, node_id)
	var total int
	e := row.Scan(&total)
	if e != nil {
		return e
	}
	if total != 0 {
		return fmt.Errorf("[ERR]該節點仍有文章(筆數: %v)", total)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 1. 查詢要刪除節點的 LftIdx 和 RftIdx
	query = fmt.Sprintf(
		"SELECT LftIdx, RftIdx FROM %s WHERE NodeID = ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	var lftIdx, rftIdx int
	err = tx.QueryRow(query, node_id).Scan(&lftIdx, &rftIdx)
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
	query = fmt.Sprintf(
		"DELETE FROM %s WHERE LftIdx >= ? AND RftIdx <= ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	_, err = tx.Exec(query, lftIdx, rftIdx)
	if err != nil {
		return err
	}

	// 3. 調整剩餘節點的索引，填補被刪除節點留下的空隙
	// 將所有 RftIdx > rftIdx 的右索引減去 width
	query = fmt.Sprintf(
		"UPDATE %s SET RftIdx = RftIdx - ? WHERE RftIdx > ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	_, err = tx.Exec(query, width, rftIdx)
	if err != nil {
		return err
	}

	// 將所有 LftIdx > rftIdx 的左索引減去 width
	query = fmt.Sprintf(
		"UPDATE %s SET LftIdx = LftIdx - ? WHERE LftIdx > ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	_, err = tx.Exec(query, width, rftIdx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	committed = true

	return err
}

// Edit 編輯
func (r *NodeRepo) Edit(account, node_id, label string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET PathName = ? WHERE NodeID = ?;",
		fmt.Sprintf(r.node_formatter, account),
	)
	_, err := r.db.Exec(query, label, node_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *NodeRepo) Move(account, parent_id, node_id, path_name string) error {
	e := r.Delete(account, node_id)
	if e != nil {
		return e
	}

	if parent_id == uuid.Nil.String() {
		_, e = r.AddParentNode(account, node_id, path_name)
		if e != nil {
			return e
		}
	} else {
		_, e = r.AddChildNode(account, parent_id, node_id, path_name)
		if e != nil {
			return e
		}
	}

	return nil
}

func (r *NodeRepo) GetAllNode(account string) (categories []entity.Node, err error) {
	// 從資料庫中讀取所有分類，並按 LftIdx 排序
	query := fmt.Sprintf(
		"SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM %s ORDER BY LftIdx ASC",
		fmt.Sprintf(r.node_formatter, account),
	)
	rows, err := r.db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c entity.Node
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}

	return
}

func (r *NodeRepo) GetNode(account, node_id string) (c entity.Node, e error) {
	query := fmt.Sprintf(
		"SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM %s WHERE NodeID = ?",
		fmt.Sprintf(r.node_formatter, account),
	)
	row := r.db.QueryRow(query, node_id)
	if err := row.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
		return
	}
	return
}

func NewNodeRepo(db *sql.DB) repo.NodeRepository {
	return &NodeRepo{
		db:                 db,
		node_formatter:     "node_%s",
		articles_formatter: "articles_%s",
	}
}
