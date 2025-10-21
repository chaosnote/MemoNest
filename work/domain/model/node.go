package model

type NodeView struct {
	NodeList     []*CategoryNode
	NodeMap      map[string]*CategoryNode
	Menu         []Menu
	MenuChildren []MenuItem
}
