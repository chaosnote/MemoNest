package usecase

import (
	"fmt"

	"github.com/google/uuid"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/domain/service"
)

type NodeUsecase struct {
	Repo repo.NodeRepository
	Tree service.NodeTree
	Menu service.MenuProvider
}

func (u *NodeUsecase) Add(account, parent_id, node_id, path_name string) (err error) {
	if parent_id == uuid.Nil.String() {
		_, err = u.Repo.AddParentNode(account, "", path_name)
	} else {
		_, err = u.Repo.AddChildNode(account, parent_id, "", path_name)
	}
	return err
}

func (u *NodeUsecase) Delete(account, node_id string) error {
	return u.Repo.Delete(account, node_id)
}

func (u *NodeUsecase) List(account string) ([]entity.Category, error) {
	return u.Repo.GetAllNode(account)
}

func (u *NodeUsecase) Edit(account, node_id, path_name string) error {
	return u.Repo.Edit(account, node_id, path_name)
}

func (u *NodeUsecase) Move(account, parent_id, node_id string) (err error) {
	var has_node = false
	if parent_id == uuid.Nil.String() {
		has_node = true
	} else {
		var parent_node entity.Category
		parent_node, err = u.Repo.GetNode(account, parent_id)
		if err != nil {
			return
		}
		if parent_node.RowID != 0 {
			has_node = true
		}
	}
	if !has_node {
		err = fmt.Errorf("無指定父節點")
		return
	}

	current_node, err := u.Repo.GetNode(account, node_id)
	if err != nil {
		return
	}
	if current_node.RowID != 0 {
		has_node = true
	}
	if !has_node {
		err = fmt.Errorf("無指定子節點")
		return
	}

	return u.Repo.Move(account, parent_id, node_id, current_node.PathName)
}

func (u *NodeUsecase) GetViewModel(account string, aes_key []byte, menu_id string) (mo model.NodeView, err error) {
	tmp_list, err := u.Repo.GetAllNode(account)
	if err != nil {
		return
	}

	node_list, node_map := u.Tree.GetInfo(tmp_list)
	for _, node := range node_list {
		u.Tree.Assign(node, aes_key)
	}

	mo.NodeList = node_list
	mo.NodeMap = node_map
	mo.Menu = u.Menu.GetList()
	mo.MenuChildren = u.Menu.GetList()[u.Menu.GetMap()[menu_id]].Children

	return
}

//-----------------------------------------------

func NewNodeUsecase(
	repo repo.NodeRepository,
	tree service.NodeTree,
	menu service.MenuProvider,
) *NodeUsecase {
	return &NodeUsecase{
		Repo: repo,
		Tree: tree,
		Menu: menu,
	}
}
