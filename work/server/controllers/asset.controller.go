package controllers

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	xxx "idv/chris/MemoNest/adapter/http"
	"idv/chris/MemoNest/adapter/http/middleware"
	zzz "idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/service"
	"idv/chris/MemoNest/utils"
)

type AssetController struct {
	img zzz.ImageProcessor
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

	helper := xxx.NewGinSession(c)
	account := helper.GetAccount()
	aes_key := []byte(helper.GetAESKey())
	plain_text, _ := utils.AesDecrypt(id, aes_key)

	file_path := ic.img.GetImageStoragePath(account, plain_text, name)

	if _, e := os.Stat(file_path); os.IsNotExist(e) {
		e = nil
		c.Status(http.StatusNotFound)
		return
	}

	mime := mime.TypeByExtension(filepath.Ext(file_path))
	c.Header("Content-Type", mime)
	c.File(file_path)
}

func NewAssetController(engine *gin.Engine, di service.DI, img zzz.ImageProcessor) {
	c := &AssetController{
		img: img,
	}

	g := engine.Group("/asset/article")
	g.Use(middleware.MustLoginMiddleware(di))
	g.GET("/image/:id/:name", c.image)
}
