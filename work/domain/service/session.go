package service

type Session interface {
	Init(account string)

	Clear()
	GetAESKey() string
	GetAccount() string
	IsLogin() bool
}
