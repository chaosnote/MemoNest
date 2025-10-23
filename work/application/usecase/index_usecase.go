package usecase

import (
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/domain/service"
)

type IndexUsecase struct {
	Menu service.MenuProvider
}

func (u *IndexUsecase) GetViewModel(account, password, menu_id string) (mo model.IndexView) {
	mo.Account = account
	mo.Password = password
	mo.Menu = u.Menu.GetList()
	mo.MenuChildren = u.Menu.GetList()[u.Menu.GetMap()[menu_id]].Children
	return
}

//-----------------------------------------------

func NewIndexUsecase(
	menu service.MenuProvider,
) *IndexUsecase {
	return &IndexUsecase{
		Menu: menu,
	}
}
