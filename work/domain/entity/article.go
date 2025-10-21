package entity

import "time"

type Article struct {
	RowID    int
	Title    string
	Content  string
	NodeID   string
	PathName string
	UpdateDt time.Time
}
