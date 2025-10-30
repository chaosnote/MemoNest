package model

type NodeViewModel struct {
	LayoutShare

	NodeList []*CategoryNode
	NodeMap  map[string]*CategoryNode
}
