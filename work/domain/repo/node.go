package repo

import "idv/chris/MemoNest/domain/entity"

type NodeRepository interface {
	AddParentNode(nodeID, pathName string) (*entity.Category, error)
	AddChildNode(parentID, nodeID, pathName string) (*entity.Category, error)
	Delete(nodeID string) error
	Edit(nodeID, label string) error
	Move(parentID, nodeID, pathName string) error
	GetAllNode() ([]entity.Category, error)
	GetNode(nodeID string) (entity.Category, error)
}
