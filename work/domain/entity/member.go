package entity

import "time"

type Member struct {
	RowID     int
	Account   string
	Password  string
	Level     int
	LastIP    string
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
