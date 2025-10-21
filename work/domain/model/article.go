package model

type ArticleView struct {
	NodeList     []*CategoryNode
	NodeMap      map[string]*CategoryNode
	Menu         []Menu
	MenuChildren []MenuItem
}
