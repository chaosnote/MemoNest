package model

import "idv/chris/MemoNest/domain/entity"

type NodeTreeViewModel struct {
	entity.Node

	Children  []*NodeTreeViewModel
	Path      string
	El_UID    string
	El_NodeID string
}
