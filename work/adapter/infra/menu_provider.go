package infra

import (
	"fmt"

	"idv/chris/MemoNest/model"
)

const (
	MP_ARTICLE = "article"
	MP_NODE    = "node"
)

type MenuProvider struct {
	menu_list []model.Menu
	menu_map  map[string]int
}

func NewMenuProvider() *MenuProvider {
	menu_list := []model.Menu{}
	menu_map := map[string]int{}

	menu_list = append(menu_list, model.Menu{
		MenuItem: model.MenuItem{
			Label: "文章",
			Path:  fmt.Sprintf("/api/v1/%s/list", MP_ARTICLE),
		},
		Children: []model.MenuItem{
			{
				Label: "清單",
				Path:  fmt.Sprintf("/api/v1/%s/list", MP_ARTICLE),
			},
			{
				Label: "新增",
				Path:  fmt.Sprintf("/api/v1/%s/fresh", MP_ARTICLE),
			},
		},
	})
	menu_map[MP_ARTICLE] = 0

	menu_list = append(menu_list, model.Menu{
		MenuItem: model.MenuItem{
			Label: "分類",
			Path:  fmt.Sprintf("/api/v1/%s/list", MP_NODE),
		},
		Children: []model.MenuItem{
			{
				Label: "清單",
				Path:  fmt.Sprintf("/api/v1/%s/list", MP_NODE),
			},
		},
	})
	menu_map[MP_NODE] = 1

	for k0 := range menu_list {
		menu_list[k0].Idx = k0
		for k1 := range menu_list[k0].Children {
			menu_list[k0].Children[k1].Idx = k1
		}
	}

	return &MenuProvider{
		menu_list: menu_list,
		menu_map:  menu_map,
	}
}

func (mp *MenuProvider) GetList() []model.Menu {
	return mp.menu_list
}

func (mp *MenuProvider) GetMap() map[string]int {
	return mp.menu_map
}
