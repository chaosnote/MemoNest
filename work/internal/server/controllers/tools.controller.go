package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"idv/chris/MemoNest/internal/service"
)

type ToolsHelper struct {
	db *sql.DB
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

func NewToolsController(rg *gin.RouterGroup, di service.DI) {
	c := &ToolsController{
		helper: &ToolsHelper{
			db: di.MariaDB,
		},
	}
	r := rg.Group("/tools")
	r.GET("/uuid", c.gen_uuid)
}
