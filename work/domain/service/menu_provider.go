package service

import "idv/chris/MemoNest/model"

type MenuProvider interface {
	GetList() []model.Menu
	GetMap() map[string]int
}
