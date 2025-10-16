package share

import (
	"fmt"
	"idv/chris/MemoNest/internal/model"
)

const (
	MK_ARTICLE = "article"
	MK_NODE    = "node"
)

func GetMenu() (menu_list []model.Menu, menu_map map[string]int) {
	menu_map = map[string]int{}

	menu_list = append(menu_list, model.Menu{
		MenuItem: model.MenuItem{
			Label: "文章",
			Path:  fmt.Sprintf("/api/v1/%s/list", MK_ARTICLE),
		},
		Children: []model.MenuItem{
			{
				Label: "清單",
				Path:  fmt.Sprintf("/api/v1/%s/list", MK_ARTICLE),
			},
			{
				Label: "新增",
				Path:  fmt.Sprintf("/api/v1/%s/fresh", MK_ARTICLE),
			},
		},
	})
	menu_map[MK_ARTICLE] = 0

	menu_list = append(menu_list, model.Menu{
		MenuItem: model.MenuItem{
			Label: "分類",
			Path:  fmt.Sprintf("/api/v1/%s/list", MK_NODE),
		},
		Children: []model.MenuItem{
			{
				Label: "清單",
				Path:  fmt.Sprintf("/api/v1/%s/list", MK_NODE),
			},
		},
	})
	menu_map[MK_NODE] = 1

	for k0 := range menu_list {
		menu_list[k0].Idx = k0
		for k1 := range menu_list[k0].Children {
			menu_list[k0].Children[k1].Idx = k1
		}
	}
	return
}
