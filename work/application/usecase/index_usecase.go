package usecase

import (
	"idv/chris/MemoNest/adapter/infra"
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/domain/service"
)

type IndexUsecase struct {
	Menu service.MenuProvider
}

func (u *IndexUsecase) GetViewModel(account, password string) (mo model.IndexView) {
	mo.Account = account
	mo.Password = password
	mo.MainMenu = u.Menu.GetList()
	mo.MenuIdx = u.Menu.GetMap()[infra.MP_INDEX]
	mo.SubMenu = u.Menu.GetList()[mo.MenuIdx].Children
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
