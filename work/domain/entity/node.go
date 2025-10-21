package entity

import "idv/chris/MemoNest/model"

type NodeViewModel struct {
	NodeList     []*model.CategoryNode
	NodeMap      map[string]*model.CategoryNode
	Menu         []model.Menu
	MenuChildren []model.MenuItem
}
