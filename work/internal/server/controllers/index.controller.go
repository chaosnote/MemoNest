package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/service"
)

type IndexController struct{}

func (u *IndexController) entry(c *gin.Context) {
	s := sessions.Default(c)
	flag, ok := s.Get(model.SK_IS_LOGIN).(bool)
	if !ok || !flag {
		c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "未登入"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "已登入"})
}

func (u *IndexController) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "OK"})
}

func NewIndexController(engine *gin.Engine, di service.DI) {
	c := &IndexController{}

	engine.GET("/", c.entry)
	engine.GET("/health", c.health)
}
