package model

type LayoutShare struct {
	MenuIdx     int
	MainMenu    []Menu
	CurrentPath string
	SubMenu     []MenuItem
}
