package model

import "time"

type Article struct {
	ArticleRowID int
	Title        string
	Content      string
	NodeID       string
	PathName     string
	UpdateDt     time.Time
}
