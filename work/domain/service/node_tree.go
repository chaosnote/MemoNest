package service

import (
	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/model"
)

type NodeTree interface {
	GetInfo(source []entity.Node) ([]*model.CategoryNode, map[string]*model.CategoryNode)
	Assign(node *model.CategoryNode, aesKey []byte)
}
