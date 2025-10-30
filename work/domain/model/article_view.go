package model

type ArticleViewModel struct {
	LayoutContext

	NodeList []*CategoryNode
	NodeMap  map[string]*CategoryNode
}
