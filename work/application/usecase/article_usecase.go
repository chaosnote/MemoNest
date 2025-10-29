package usecase

import (
	"fmt"
	"strconv"

	"idv/chris/MemoNest/adapter/infra"
	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/model"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
)

type ArticleUsecase struct {
	Repo repo.ArticleRepository
	Tree service.NodeTree
	Menu service.MenuProvider
	Img  service.ImageProcessor
}

func (u *ArticleUsecase) Add(account, node_id, article_title, article_content string) (err error) {
	var row_id int
	row_id, err = u.Repo.Add(account, node_id)
	if err != nil {
		err = utils.ParseSQLError(err, "新增文章失敗")
		return
	}

	article_id := fmt.Sprintf("%v", row_id)

	article_content = u.Img.ProcessBase64Images(account, article_id, article_content)
	err = u.Repo.Update(account, row_id, article_title, article_content)
	if err != nil {
		err = utils.ParseSQLError(err, "新增文章失敗")
		return
	}
	u.Img.CleanupUnusedImages(account, article_id, article_content)

	return
}

func (u *ArticleUsecase) Del(account, plain_text string) (err error) {
	var id int
	id, err = strconv.Atoi(plain_text)
	if err != nil {
		return
	}
	err = u.Repo.Delete(account, id)
	if err != nil {
		err = utils.ParseSQLError(err, "刪除文章失敗")
		return
	}
	u.Img.DelImageDir(account, fmt.Sprintf("%v", id))

	return
}

func (u *ArticleUsecase) Edit(account, plain_text string) (data entity.Article, err error) {
	var id int
	id, err = strconv.Atoi(plain_text)
	if err != nil {
		return
	}

	list, err := u.Repo.Get(account, id)
	if err != nil {
		err = utils.ParseSQLError(err, "編輯文章失敗")
		return
	}
	if len(list) == 0 {
		err = fmt.Errorf("無指定文章")
		return
	}
	data = list[0]
	return
}

func (u *ArticleUsecase) Renew(account, article_id, article_title, article_content string) (err error) {
	var row_id int
	row_id, err = strconv.Atoi(article_id)
	if err != nil {
		return
	}

	article_content = u.Img.ProcessBase64Images(account, article_id, article_content)
	err = u.Repo.Update(account, row_id, article_title, article_content)
	if err != nil {
		err = utils.ParseSQLError(err, "編輯文章失敗")
		return
	}
	u.Img.CleanupUnusedImages(account, article_id, article_content)

	return
}

func (u *ArticleUsecase) List(account, query string) (list []entity.Article, err error) {
	if len(query) > 0 {
		list, err = u.Repo.Query(account, query)
	} else {
		list, err = u.Repo.List(account)
	}

	if err != nil {
		err = utils.ParseSQLError(err, "查詢文章清單失敗")
	}
	return
}

func (u *ArticleUsecase) GetViewModel(account string, aes_key []byte) (mo model.ArticleView, err error) {
	tmp_list, err := u.Repo.GetAllNode(account)
	if err != nil {
		err = utils.ParseSQLError(err, "查詢節點失敗")
		return
	}

	node_list, node_map := u.Tree.GetInfo(tmp_list)
	for _, node := range node_list {
		u.Tree.Assign(node, aes_key)
	}

	mo.NodeList = node_list
	mo.NodeMap = node_map
	mo.MainMenu = u.Menu.GetList()
	mo.MenuIdx = u.Menu.GetMap()[infra.MP_ARTICLE]
	mo.SubMenu = u.Menu.GetList()[mo.MenuIdx].Children
	return
}

//-----------------------------------------------

func NewArticleUsecase(
	repo repo.ArticleRepository,
	tree service.NodeTree,
	menu service.MenuProvider,
	img service.ImageProcessor,
) *ArticleUsecase {
	return &ArticleUsecase{
		Repo: repo,
		Tree: tree,
		Menu: menu,
		Img:  img,
	}
}
