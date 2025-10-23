package repo

import "idv/chris/MemoNest/domain/entity"

type NodeRepository interface {
	AddParentNode(account, node_id, path_name string) (*entity.Category, error)
	AddChildNode(account, parent_id, node_id, path_name string) (*entity.Category, error)
	Delete(account, node_id string) error
	Edit(account, node_id, label string) error
	Move(account, parentID, node_id, path_name string) error
	GetAllNode(account string) ([]entity.Category, error)
	GetNode(account, node_id string) (entity.Category, error)
}
