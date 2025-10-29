package handle

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/utils"
)

type AssetHandler struct {
	CommonHandler

	UC *usecase.AssetUsecase
}

func (h *AssetHandler) Image(c *gin.Context) {
	const msg = "image"

	var e error
	defer func() {
		if e != nil {
			h.Log.Error(msg, zap.Error(e))
			return
		}
	}()

	id := c.Params.ByName("id")
	name := c.Params.ByName("name")
	h.Log.Info(msg, zap.String("id", id), zap.String("name", name))

	h.Session.Init(c)
	aes_key := []byte(h.Session.GetAESKey())
	account := h.Session.GetAccount()
	plain_text, e := utils.AesDecrypt(id, aes_key)
	if e != nil {
		return
	}

	file_path := h.UC.GetImageStoragePath(account, plain_text, name)
	if _, e := os.Stat(file_path); os.IsNotExist(e) {
		c.Status(http.StatusNotFound)
		return
	}

	mime := mime.TypeByExtension(filepath.Ext(file_path))
	c.Header("Content-Type", mime)
	c.File(file_path)
}
