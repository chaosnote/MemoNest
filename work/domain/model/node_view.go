package model

type NodeViewModel struct {
	LayoutContext

	NodeList []*CategoryNode
	NodeMap  map[string]*CategoryNode
}
