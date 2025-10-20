package repo

import "idv/chris/MemoNest/model"

type ArticleRepository interface {
	Add(nodeID string) (int, error)
	Update(rowID int, title, content string) error
	Delete(id int) error
	Get(id int) ([]model.Article, error)
	List() ([]model.Article, error)
	Query(input string) ([]model.Article, error)
	GetAllNode() ([]model.Category, error)
}
