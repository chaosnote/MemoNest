package model

type NodeViewModel struct {
	LayoutContext

	NodeList []*NodeTreeViewModel
	NodeMap  map[string]*NodeTreeViewModel
}
