package usecase

import (
	"idv/chris/MemoNest/adapter/infra"
	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
)

type IndexUsecase struct {
	Menu service.MenuProvider
	Repo repo.ArticleRepository
}

func (u *IndexUsecase) List(account string) (list []entity.Article, err error) {
	list, err = u.Repo.List(account)
	if err != nil {
		err = utils.ParseSQLError(err, "查詢文章清單失敗")
	}
	return
}

func (u *IndexUsecase) GetViewModel(account, password string) (mo model.IndexViewModel) {
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
	repo repo.ArticleRepository,
) *IndexUsecase {
	return &IndexUsecase{
		Menu: menu,
		Repo: repo,
	}
}
