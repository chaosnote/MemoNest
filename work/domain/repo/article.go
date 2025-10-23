package repo

import (
	"idv/chris/MemoNest/domain/entity"
)

type ArticleRepository interface {
	Add(node_id string) (int, error)
	Update(row_id int, title, content string) error
	Delete(id int) error
	Get(id int) ([]entity.Article, error)
	List() ([]entity.Article, error)
	Query(input string) ([]entity.Article, error)
	GetAllNode() ([]entity.Category, error)
}
