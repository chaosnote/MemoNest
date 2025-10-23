package model

type IndexView struct {
	Account      string
	Password     string
	Menu         []Menu
	MenuChildren []MenuItem
}
