package model

type MenuItem struct {
	Idx   int
	Label string
	Path  string
}

type Menu struct {
	MenuItem
	Children []MenuItem
}
