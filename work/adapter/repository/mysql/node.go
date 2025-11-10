package mysql

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
)

type NodeRepo struct {
	db *sql.DB
}

func (r *NodeRepo) AddParentNode(account, parent_id, node_id, path_name string) (n entity.Node, e error) {
	query := "CALL `sp_node_add_parent`(?, ?, ?, ?)"
	row := r.db.QueryRow(query, account, parent_id, node_id, path_name)

	n = entity.Node{}
	e = row.Scan(
		&n.RowID,
		&n.NodeID,
		&n.ParentID,
		&n.PathName,
		&n.LftIdx,
		&n.RftIdx,
	)
	if e != nil {
		return
	}

	return
}

// AddChildNode 插入一個新的分類節點
func (r *NodeRepo) AddChildNode(account, parent_id, node_id, path_name string) (n entity.Node, e error) {
	query := "CALL `sp_node_add_child`(?, ?, ?, ?)"
	row := r.db.QueryRow(query, account, parent_id, node_id, path_name)

	e = row.Err()
	if e != nil {
		return
	}

	e = row.Scan(
		&n.RowID,
		&n.NodeID,
		&n.ParentID,
		&n.PathName,
		&n.LftIdx,
		&n.RftIdx,
	)
	if e != nil {
		return
	}

	return
}

// Delete 移除指定的分類節點及其所有後代節點
func (r *NodeRepo) Delete(account, node_id string) error {
	query := "CALL `sp_node_del`(?, ?)"
	_, e := r.db.Exec(query, account, node_id)
	return e
}

// Edit 編輯
func (r *NodeRepo) Edit(account, node_id, path_name string) error {
	query := "CALL `sp_node_edit`(?, ?, ?)"
	_, e := r.db.Exec(query, account, node_id, path_name)
	return e
}

func (r *NodeRepo) Move(account, parent_id, node_id, path_name string) error {
	e := r.Delete(account, node_id)
	if e != nil {
		return e
	}

	if parent_id == uuid.Nil.String() {
		fmt.Println("AddParentNode", account, parent_id, node_id, path_name)
		_, e = r.AddParentNode(account, parent_id, node_id, path_name)
		if e != nil {
			return e
		}
	} else {
		fmt.Println("AddChildNode", account, parent_id, node_id, path_name)
		_, e = r.AddChildNode(account, parent_id, node_id, path_name)
		if e != nil {
			return e
		}
	}

	return nil
}

func (r *NodeRepo) GetAllNode(account string) (list []entity.Node, err error) {
	// 從資料庫中讀取所有分類，並按 LftIdx 排序
	query := "CALL sp_node_list(?)"

	rows, err := r.db.Query(query, account)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c entity.Node
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return list, err
		}
		list = append(list, c)
	}

	return
}

func (r *NodeRepo) GetNode(account, node_id string) (n entity.Node, e error) {
	query := "CALL sp_node_get(?,?)"

	row := r.db.QueryRow(query, account, node_id)
	e = row.Scan(&n.RowID, &n.NodeID, &n.ParentID, &n.PathName, &n.LftIdx, &n.RftIdx)
	if e != nil {
		return
	}
	return
}

func NewNodeRepo(db *sql.DB) repo.NodeRepository {
	return &NodeRepo{
		db: db,
	}
}
