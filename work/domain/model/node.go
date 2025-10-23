package model

type NodeView struct {
	LayoutShare

	NodeList []*CategoryNode
	NodeMap  map[string]*CategoryNode
}
