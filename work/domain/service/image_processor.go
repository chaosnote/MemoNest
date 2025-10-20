package service

type ImageProcessor interface {
	ProcessBase64Images(account, articleID, content string) string
	CleanupUnusedImages(account, articleID, content string)
	DelImageDir(account, articleID string)
	GetImageStoragePath(account, article_id, filename string) string
}
