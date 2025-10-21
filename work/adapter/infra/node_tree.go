package infra

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/utils"
)

type NodeTree struct {
}

func (nt *NodeTree) GetInfo(source []entity.Category) ([]*model.CategoryNode, map[string]*model.CategoryNode) {
	root_id := uuid.Nil.String()
	node_map := make(map[string]*model.CategoryNode)

	// 第一次遍歷：建立節點地圖
	for _, category := range source {
		node_map[category.NodeID] = &model.CategoryNode{
			Category: category,
		}
	}

	// 第二次遍歷：建立樹狀結構並生成路徑
	node_list := []*model.CategoryNode{}
	for _, category := range source {
		current_node := node_map[category.NodeID]

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
		current_node.Path = "/" + strings.Join(path_seg, "/")

		// 處理樹狀結構
		if category.ParentID == root_id {
			node_list = append(node_list, current_node)
		} else {
			if parent, ok := node_map[category.ParentID]; ok {
				parent.Children = append(parent.Children, current_node)
			}
		}
	}

	return node_list, node_map
}

func (nh *NodeTree) Assign(node *model.CategoryNode, aes_key []byte) {
	sUID, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", node.RowID)), aes_key)
	node.El_UID = sUID

	sNodeID, _ := utils.AesEncrypt([]byte(node.NodeID), aes_key)
	node.El_NodeID = sNodeID

	for _, child := range node.Children {
		nh.Assign(child, aes_key)
	}
}

func NewNodeTree() *NodeTree {
	return &NodeTree{}
}
