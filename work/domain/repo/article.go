package repo

import (
	"idv/chris/MemoNest/domain/entity"
)

type ArticleRepository interface {
	Add(nodeID string) (int, error)
	Update(rowID int, title, content string) error
	Delete(id int) error
	Get(id int) ([]entity.Article, error)
	List() ([]entity.Article, error)
	Query(input string) ([]entity.Article, error)
	GetAllNode() ([]entity.Category, error)
}
