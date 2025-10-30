package model

type ArticleViewModel struct {
	LayoutShare

	NodeList []*CategoryNode
	NodeMap  map[string]*CategoryNode
}
