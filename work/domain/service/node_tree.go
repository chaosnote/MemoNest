package service

import "idv/chris/MemoNest/model"

type NodeTree interface {
	GetInfo(source []model.Category) ([]*model.CategoryNode, map[string]*model.CategoryNode)
}
