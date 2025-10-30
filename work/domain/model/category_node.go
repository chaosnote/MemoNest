package model

import "idv/chris/MemoNest/domain/entity"

type CategoryNode struct {
	entity.Node

	Children  []*CategoryNode
	Path      string
	El_UID    string
	El_NodeID string
}
