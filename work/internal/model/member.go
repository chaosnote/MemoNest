package model

import "time"

type Member struct {
	ID        uint
	Merchant  string
	Account   string
	Password  string
	LastIP    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Wallet    int64
}
