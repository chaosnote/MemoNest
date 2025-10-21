package entity

type Category struct {
	RowID    int
	NodeID   string
	ParentID string
	PathName string
	LftIdx   int
	RftIdx   int
}
