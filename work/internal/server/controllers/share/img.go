package share

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
)

func ProcessBase64Images(userID, articleID, html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if !strings.HasPrefix(src, "data:image/") {
			return
		}
		commaIdx := strings.Index(src, ",")
		raw := src[commaIdx+1:]
		data, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return
		}

		filename := fmt.Sprintf("img_%d.png", time.Now().UnixNano())
		path := GetImageStoragePath(userID, articleID, filename)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, data, 0644)

		s.SetAttr("src", fmt.Sprintf("/article/%s", filename))
	})

	htmlOut, _ := doc.Find("body").Html()
	return htmlOut
}

func GetImageStoragePath(userID, articleID, filename string) string {
	// exePath, _ := os.Executable()
	// fmt.Println("DIR", exePath)
	// baseDir := filepath.Dir(exePath)
	baseDir := "./dist"
	return filepath.Join(baseDir, "uploads", userID, "article", articleID, filename)
}

func CleanupUnusedImages(userID, articleID, html string) {
	used := ExtractImageFilenamesFromHTML(html)
	dir := filepath.Dir(GetImageStoragePath(userID, articleID, "dummy"))
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if !slices.Contains(used, f.Name()) {
			os.Remove(filepath.Join(dir, f.Name()))
		}
	}
}

func ExtractImageFilenamesFromHTML(html string) []string {
	re := regexp.MustCompile(`<img[^>]+src="/article/([^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)
	var filenames []string
	for _, m := range matches {
		filenames = append(filenames, m[1])
	}
	return filenames
}
