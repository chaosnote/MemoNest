package service

import "github.com/gin-gonic/gin"

type Session interface {
	Init(ctx *gin.Context)

	Clear()
	IsLogin() bool
	GetAESKey() string
	SetAccount(string)
	GetAccount() string

	Refresh()
}
