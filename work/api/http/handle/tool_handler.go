package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"idv/chris/MemoNest/application/usecase"
)

type ToolHandler struct {
	UC *usecase.ToolUsecase
}

func (h *ToolHandler) GenUUID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Code": "OK",
		"uuid": uuid.NewString(),
	})
}
