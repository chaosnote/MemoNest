package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/service"
)

type IndexController struct{}

func (u *IndexController) entry(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/health")
}

func (u *IndexController) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": "OK"})
}

func NewIndexController(engine *gin.Engine, di service.DI) {
	c := &IndexController{}

	engine.GET("/", c.entry)
	engine.GET("/health", c.health)
}
