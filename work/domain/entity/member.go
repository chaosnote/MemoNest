package entity

import "time"

type Member struct {
	RowID     int
	Account   string
	Password  string
	LastIP    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
