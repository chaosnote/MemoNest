package service

import (
	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/model"
)

type NodeTree interface {
	GetInfo(source []entity.Node) ([]*model.NodeTreeViewModel, map[string]*model.NodeTreeViewModel)
	Encrypt(node *model.NodeTreeViewModel, aesKey []byte)
}
