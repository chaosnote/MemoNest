package usecase

import (
	"fmt"
	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/model"
	"strconv"
)

type ArticleUsecase struct {
	Repo repo.ArticleRepository
	Tree service.NodeTree
	Menu service.MenuProvider
	Img  service.ImageProcessor
}

func (u *ArticleUsecase) Add(account, node_id, article_title, article_content string) (err error) {
	var row_id int
	row_id, err = u.Repo.Add(node_id)
	if err != nil {
		return
	}

	article_id := fmt.Sprintf("%v", row_id)

	article_content = u.Img.ProcessBase64Images(account, article_id, article_content)
	err = u.Repo.Update(row_id, article_title, article_content)
	if err != nil {
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
	err = u.Repo.Delete(id)
	if err != nil {
		return
	}
	u.Img.DelImageDir(account, fmt.Sprintf("%v", id))

	return
}

func (u *ArticleUsecase) Edit(plain_text string) (data model.Article, err error) {
	var id int
	id, err = strconv.Atoi(plain_text)
	if err != nil {
		return
	}

	list, err := u.Repo.Get(id)
	if err != nil {
		return
	}
	if len(list) == 0 {
		err = fmt.Errorf("無指定資料")
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
	err = u.Repo.Update(row_id, article_title, article_content)
	if err != nil {
		return
	}
	u.Img.CleanupUnusedImages(account, article_id, article_content)

	return
}

func (u *ArticleUsecase) List(query string) (list []model.Article, err error) {
	if len(query) > 0 {
		list, err = u.Repo.Query(query)
	} else {
		list, err = u.Repo.List()
	}

	return
}

func (u *ArticleUsecase) GetViewModel(aes_key []byte, menu_id string) (mo entity.ArticleViewModel, err error) {
	tmp_list, err := u.Repo.GetAllNode()
	if err != nil {
		return
	}

	node_list, node_map := u.Tree.GetInfo(tmp_list)
	for _, node := range node_list {
		u.Tree.Assign(node, aes_key)
	}

	mo.NodeList = node_list
	mo.NodeMap = node_map
	mo.Menu = u.Menu.GetList()
	mo.MenuChildren = u.Menu.GetList()[u.Menu.GetMap()[menu_id]].Children

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
