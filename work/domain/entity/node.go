package entity

type Node struct {
	RowID    int
	NodeID   string
	ParentID string
	PathName string
	LftIdx   int
	RftIdx   int
}
