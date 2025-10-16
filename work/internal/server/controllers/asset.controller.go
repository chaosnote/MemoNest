package controllers

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/server/controllers/share"
	"idv/chris/MemoNest/internal/server/middleware"
	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"
)

type AssetController struct {
}

func (ic *AssetController) image(c *gin.Context) {
	const msg = "image"
	logger := utils.NewFileLogger("./dist/logs/article/image", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			return
		}
	}()

	id := c.Params.ByName("id")
	name := c.Params.ByName("name")
	logger.Info(msg, zap.String("id", id), zap.String("name", name))

	s := sessions.Default(c)
	key := []byte(s.Get(model.SK_AES_KEY).(string))
	account := s.Get(model.SK_ACCOUNT).(string)
	output, _ := utils.AesDecrypt(id, key)

	file_path := share.GetImageStoragePath(account, output, name)

	if _, e := os.Stat(file_path); os.IsNotExist(e) {
		e = nil
		c.Status(http.StatusNotFound)
		return
	}

	mime := mime.TypeByExtension(filepath.Ext(file_path))
	c.Header("Content-Type", mime)
	c.File(file_path)
}

func NewAssetController(engine *gin.Engine, di service.DI) {
	c := &AssetController{}

	g := engine.Group("/asset/article")
	g.Use(middleware.MustLoginMiddleware(di))
	g.GET("/image/:id/:name", c.image)
}
