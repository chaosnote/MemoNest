package service

import "idv/chris/MemoNest/domain/model"

type MenuProvider interface {
	GetList() []model.Menu
	GetMap() map[string]int
}
