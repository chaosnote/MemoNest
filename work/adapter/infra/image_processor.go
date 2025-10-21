package infra

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"idv/chris/MemoNest/domain/model"
)

type ImageProcessor struct {
}

func (ip *ImageProcessor) ProcessBase64Images(account, article_id, content string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return content
	}

	idx := time.Now().UnixNano()
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if !strings.HasPrefix(src, "data:image/") {
			re := regexp.MustCompile(`/asset/article/image/([a-f0-9]{32})/`)
			s.SetAttr("src", re.ReplaceAllString(src, fmt.Sprintf("/asset/article/image/%s/", model.IMG_ENCRYPT)))
		} else {
			commaIdx := strings.Index(src, ",")
			raw := src[commaIdx+1:]
			data, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				return
			}

			filename := fmt.Sprintf("img_%d.png", idx)
			idx++
			path := ip.GetImageStoragePath(account, article_id, filename)
			os.MkdirAll(filepath.Dir(path), 0755)
			os.WriteFile(path, data, 0644)

			s.SetAttr("src", fmt.Sprintf("%s/%s/%s", model.IMG_SRC, model.IMG_ENCRYPT, filename)) // 留意:這邊需對應 ExtractImageFilenamesFromHTML 正則
		}
	})

	output, _ := doc.Find("body").Html()
	return output
}

func (ip *ImageProcessor) GetImageStoragePath(account, article_id, filename string) string {
	// exePath, _ := os.Executable()
	// fmt.Println("DIR", exePath)
	// baseDir := filepath.Dir(exePath)
	baseDir := "./dist"
	return filepath.Join(baseDir, "uploads", account, "article", article_id, filename)
}

func (ip *ImageProcessor) CleanupUnusedImages(account, article_id, html string) {
	used := ip.ExtractImageFilenamesFromHTML(html)
	dir := filepath.Dir(ip.GetImageStoragePath(account, article_id, "dummy"))
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if !slices.Contains(used, f.Name()) {
			os.Remove(filepath.Join(dir, f.Name()))
		}
	}
}

func (ip *ImageProcessor) ExtractImageFilenamesFromHTML(content string) []string {
	re := regexp.MustCompile(fmt.Sprintf(`<img[^>]+src="%s(?:/[^/]+)?/([^"]+)"`, model.IMG_SRC))
	matches := re.FindAllStringSubmatch(content, -1)
	var filenames []string
	for _, m := range matches {
		filenames = append(filenames, m[1])
	}
	return filenames
}

func (ip *ImageProcessor) DelImageDir(account, article_id string) {
	baseDir := "./dist"
	dir := filepath.Join(baseDir, "uploads", account, "article", article_id)
	os.RemoveAll(dir)
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
