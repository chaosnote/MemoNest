package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/service"
)

type IndexController struct{}

func (u *IndexController) entry(service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"error": "", "message": "首頁"})
	}
}

func (u *IndexController) health(service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"error": "", "message": "OK"})
	}
}

func NewIndexController(engine *gin.Engine, di service.DI) {
	c := &IndexController{}

	engine.GET("/", c.entry(di))
	engine.GET("/health", c.health(di))
}
