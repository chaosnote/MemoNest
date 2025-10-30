package usecase

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"idv/chris/MemoNest/adapter/infra"
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
)

type IndexUsecase struct {
	Menu service.MenuProvider
	Repo repo.ArticleRepository
}

func (u *IndexUsecase) List(account string, aes_key []byte) (list []model.ArticleListViewModel, err error) {
	articles, err := u.Repo.List(account)
	if err != nil {
		err = utils.ParseSQLError(err, "查詢文章清單失敗")
		return
	}

	loc, _ := time.LoadLocation("Asia/Taipei")
	for key, item := range articles {
		output := model.ArticleListViewModel{
			Article: item,
		}
		id, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", item.RowID)), aes_key)
		output.El_Idx = key + 1
		output.El_ID = id
		output.El_Time = item.UpdateDt.In(loc).Format("2006-01-02 15:04")
		output.El_Content = template.HTML(strings.ReplaceAll(item.Content, model.IMG_ENCRYPT, id))

		list = append(list, output)
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
