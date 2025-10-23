package repo

import (
	"idv/chris/MemoNest/domain/entity"
)

type ArticleRepository interface {
	Add(account, node_id string) (int, error)
	Update(account string, row_id int, title, content string) error
	Delete(account string, id int) error
	Get(account string, id int) ([]entity.Article, error)
	List(account string) ([]entity.Article, error)
	Query(account, input string) ([]entity.Article, error)
	GetAllNode(account string) ([]entity.Category, error)
}
