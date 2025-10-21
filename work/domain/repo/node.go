package repo

import "idv/chris/MemoNest/model"

type NodeRepository interface {
	AddParentNode(nodeID, pathName string) (*model.Category, error)
	AddChildNode(parentID, nodeID, pathName string) (*model.Category, error)
	Delete(nodeID string) error
	Edit(nodeID, label string) error
	Move(parentID, nodeID, pathName string) error
	GetAllNode() ([]model.Category, error)
	GetNode(nodeID string) (model.Category, error)
}
