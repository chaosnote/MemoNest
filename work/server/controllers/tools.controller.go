package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ToolsHelper struct {
}

//-----------------------------------------------

type ToolsController struct {
	helper *ToolsHelper
}

func (tc *ToolsController) gen_uuid(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Code": "OK",
		"uuid": uuid.NewString(),
	})
}

func NewToolsController(rg *gin.RouterGroup) {
	c := &ToolsController{
		helper: &ToolsHelper{},
	}
	r := rg.Group("/tools")
	r.GET("/uuid", c.gen_uuid)
}
