package model

import (
	"html/template"

	"idv/chris/MemoNest/domain/entity"
)

type ArticleListViewModel struct {
	entity.Article

	El_Idx     int
	El_ID      string
	El_Time    string
	El_Content template.HTML
}
