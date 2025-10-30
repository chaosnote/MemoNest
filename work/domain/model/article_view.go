package model

type ArticleViewModel struct {
	LayoutContext

	NodeList []*NodeTreeViewModel
	NodeMap  map[string]*NodeTreeViewModel
}
