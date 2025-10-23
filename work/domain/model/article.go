package model

type ArticleView struct {
	LayoutShare

	NodeList []*CategoryNode
	NodeMap  map[string]*CategoryNode
}
