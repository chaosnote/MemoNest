package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/server/controllers/share"
	"idv/chris/MemoNest/internal/server/middleware"
	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"
)

type ArticleHelper struct {
	db *sql.DB
}

func (ah *ArticleHelper) getAllNode() (categories []model.Category) {
	rows, err := ah.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return
		}
		categories = append(categories, c)
	}

	return
}

func (ah *ArticleHelper) add(node_id string) (id int, e error) {
	t := time.Now().UTC()
	row := ah.db.QueryRow("CALL insert_article(?, ?, ?, ?, ?) ;", "", "", t, t, node_id)
	e = row.Scan(&id)
	if e != nil {
		return
	}
	return
}

func (ah *ArticleHelper) del(id int) (e error) {
	_, e = ah.db.Exec(`DELETE from articles where RowID = ? ;`, id)
	if e != nil {
		return
	}
	return
}

func (ah *ArticleHelper) update(row_id int, title, content string) error {
	t := time.Now().UTC()
	query := `UPDATE articles SET Title = ?, Content = ?, UpdateDt =? WHERE RowID = ?;`
	_, e := ah.db.Exec(query, title, content, t, row_id)
	if e != nil {
		return e
	}
	return nil
}

func (ah *ArticleHelper) get(id int) (articles []model.Article) {
	rows, e := ah.db.Query(`
		SELECT 
			a.RowID AS ArticleRowID,
			a.Title,
			a.Content,
			a.NodeID,
			c.PathName,
			a.UpdateDt
		FROM articles as a
		JOIN categories as c ON a.NodeID = c.NodeID
		WHERE a.RowID = ? ;
	`, id)
	if e != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		e = rows.Scan(
			&article.RowID,
			&article.Title,
			&article.Content,
			&article.NodeID,
			&article.PathName,
			&article.UpdateDt,
		)
		if e != nil {
			return
		}
		articles = append(articles, article)
	}
	return
}

func (ah *ArticleHelper) list() (articles []model.Article) {
	rows, e := ah.db.Query(`
		SELECT 
			a.RowID AS ArticleRowID,
			a.Title,
			a.Content,
			a.NodeID,
			c.PathName,
			a.UpdateDt
		FROM articles as a
		JOIN categories as c ON a.NodeID = c.NodeID
		ORDER BY a.UpdateDt DESC ;
	`)
	if e != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		e = rows.Scan(
			&article.RowID,
			&article.Title,
			&article.Content,
			&article.NodeID,
			&article.PathName,
			&article.UpdateDt,
		)
		if e != nil {
			return
		}
		articles = append(articles, article)
	}
	return
}

func (ah *ArticleHelper) composit(input string) (query string, args []interface{}) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}

	andParts := strings.Split(input, "&")
	var conditions []string

	for _, part := range andParts {
		orParts := strings.Split(part, "+")
		orPartsClean := []string{}
		for _, kw := range orParts {
			kw = strings.TrimSpace(kw)
			if kw != "" {
				orPartsClean = append(orPartsClean, kw)
			}
		}

		if len(orPartsClean) == 1 {
			// 單一關鍵字：Title 或 Content 任一欄位包含
			pattern := "%" + orPartsClean[0] + "%"
			conditions = append(conditions, "(a.Title LIKE ? OR a.Content LIKE ? OR c.PathName LIKE ?)")
			args = append(args, pattern, pattern, pattern)
		} else if len(orPartsClean) > 1 {
			// 多關鍵字：Title 或 Content 任一欄位符合 REGEXP
			pattern := strings.Join(orPartsClean, "|")
			conditions = append(conditions, "(a.Title REGEXP ? OR a.Content REGEXP ? OR c.PathName REGEXP ?)")
			args = append(args, pattern, pattern, pattern)
		}
	}

	query = `
        SELECT 
            a.RowID AS ArticleRowID,
            a.Title,
            a.Content,
            a.NodeID,
            c.PathName,
            a.UpdateDt
        FROM articles AS a
        JOIN categories AS c ON a.NodeID = c.NodeID
        WHERE ` + strings.Join(conditions, " AND ") + `
        ORDER BY a.UpdateDt DESC
    `
	return
}

func (ah *ArticleHelper) query(input string) (articles []model.Article) {
	cmd, args := ah.composit(input)
	rows, e := ah.db.Query(cmd, args...)
	if e != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		e = rows.Scan(
			&article.RowID,
			&article.Title,
			&article.Content,
			&article.NodeID,
			&article.PathName,
			&article.UpdateDt,
		)
		if e != nil {
			return
		}
		articles = append(articles, article)
	}
	return
}

//-----------------------------------------------

type ArticleController struct {
	helper *ArticleHelper
}

func (ac *ArticleController) fresh(c *gin.Context) {
	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "article", "fresh.html")},
		Pattern: []string{},
	}
	tmpl, e := utils.RenderTemplate(config)

	defer func() {
		if e != nil {
			http.Error(c.Writer, e.Error(), http.StatusInternalServerError)
		}
	}()
	if e != nil {
		return
	}

	_, node_map := share.GenNodeInfo(ac.helper.getAllNode())
	menu_list, menu_map := share.GetMenu()
	e = tmpl.ExecuteTemplate(c.Writer, "fresh.html", gin.H{
		"Title":        "增加文章",
		"Menu":         menu_list,
		"Children":     menu_list[menu_map[share.MK_ARTICLE]].Children,
		"NodeMap":      node_map,
		"ArticleTitle": "請輸入文章標題",
	})
	if e != nil {
		return
	}
}

func (ac *ArticleController) add(c *gin.Context) {
	const msg = "add"
	logger := utils.NewFileLogger("./dist/logs/article/add", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		NodeID  string `json:"node_id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))
	var row_id int
	row_id, e = ac.helper.add(param.NodeID)
	if e != nil {
		return
	}

	article_id := fmt.Sprintf("%v", row_id)

	helper := share.NewSessionHelper(c)
	account := helper.GetAccount()
	content := share.ProcessBase64Images(account, article_id, param.Content)
	e = ac.helper.update(row_id, param.Title, content)
	if e != nil {
		return
	}
	share.CleanupUnusedImages(account, article_id, content)

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (ac *ArticleController) del(c *gin.Context) {
	const msg = "del"
	logger := utils.NewFileLogger("./dist/logs/article/del", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		ID string `json:"id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	helper := share.NewSessionHelper(c)
	account := helper.GetAccount()
	aes_key := []byte(helper.GetAESKey())
	plain_text, _ := utils.AesDecrypt(param.ID, aes_key)

	var id int
	id, e = strconv.Atoi(plain_text)
	if e != nil {
		return
	}
	e = ac.helper.del(id)
	if e != nil {
		return
	}

	share.DelImageDir(account, fmt.Sprintf("%v", id))

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (ac *ArticleController) edit(c *gin.Context) {
	const msg = "edit"
	logger := utils.NewFileLogger("./dist/logs/article/edit", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()

	helper := share.NewSessionHelper(c)
	aes_key := []byte(helper.GetAESKey())
	plain_text, e := utils.AesDecrypt(c.Param("id"), aes_key)
	if e != nil {
		return
	}

	var id int
	id, e = strconv.Atoi(plain_text)
	if e != nil {
		return
	}

	list := ac.helper.get(id)
	if len(list) == 0 {
		e = fmt.Errorf("無指定資料")
		return
	}
	data := list[0]

	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "article", "edit.html")},
		Pattern: []string{},
		Funcs: map[string]any{
			"trans": func(id, data string) string {
				return strings.ReplaceAll(data, model.IMG_ENCRYPT, id)
			},
		},
	}

	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	menu_list, menu_map := share.GetMenu()
	e = tmpl.ExecuteTemplate(c.Writer, "edit.html", gin.H{
		"Title":          "修改文章",
		"Menu":           menu_list,
		"Children":       menu_list[menu_map[share.MK_ARTICLE]].Children,
		"ID":             c.Param("id"),
		"PathName":       data.PathName,
		"ArticleTitle":   data.Title,
		"ArticleContent": data.Content,
	})
	if e != nil {
		return
	}
}

func (ac *ArticleController) renew(c *gin.Context) {
	const msg = "renew"
	logger := utils.NewFileLogger("./dist/logs/article/renew", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()
	var param struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		ID      string `json:"id"`
	}
	e = c.BindJSON(&param)
	if e != nil {
		return
	}
	logger.Info(msg, zap.Any("params", param))

	helper := share.NewSessionHelper(c)
	account := helper.GetAccount()
	aes_key := []byte(helper.GetAESKey())
	article_id, _ := utils.AesDecrypt(param.ID, aes_key)

	var row_id int
	row_id, e = strconv.Atoi(article_id)
	if e != nil {
		return
	}

	content := share.ProcessBase64Images(account, article_id, param.Content)
	e = ac.helper.update(row_id, param.Title, content)
	if e != nil {
		return
	}
	share.CleanupUnusedImages(account, article_id, content)

	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}

func (ac *ArticleController) list(c *gin.Context) {
	const msg = "list"
	logger := utils.NewFileLogger("./dist/logs/article/list", "console", 1)
	var e error
	defer func() {
		if e != nil {
			logger.Error(msg, zap.Error(e))
			c.JSON(http.StatusOK, gin.H{"Code": e.Error()})
		}
	}()

	q := c.Query("q")

	var list []model.Article
	if len(q) > 0 {
		list = ac.helper.query(q)
	} else {
		list = ac.helper.list()
	}

	helper := share.NewSessionHelper(c)
	aes_key := []byte(helper.GetAESKey())

	dir := filepath.Join("./assets", "templates")
	config := utils.TemplateConfig{
		Layout:  filepath.Join(dir, "layout", "share.html"),
		Page:    []string{filepath.Join(dir, "page", "article", "list.html")},
		Pattern: []string{},
		Funcs: map[string]any{
			"format": func(t time.Time) string {
				loc, _ := time.LoadLocation("Asia/Taipei")
				return t.In(loc).Format("2006-01-02 15:04")
			},
			"encrypt": func(id int) string {
				cipher_text, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", id)), aes_key)
				return string(cipher_text)
			},
			"trans": func(id int, data string) template.HTML {
				output, _ := utils.AesEncrypt([]byte(fmt.Sprintf("%v", id)), aes_key)
				return template.HTML(strings.ReplaceAll(data, model.IMG_ENCRYPT, output))
			},
		},
	}
	tmpl, e := utils.RenderTemplate(config)
	if e != nil {
		return
	}
	menu_list, menu_map := share.GetMenu()
	e = tmpl.ExecuteTemplate(c.Writer, "list.html", gin.H{
		"Title":    "文章清單",
		"Menu":     menu_list,
		"Children": menu_list[menu_map[share.MK_ARTICLE]].Children,
		"List":     list,
	})
	if e != nil {
		return
	}
}

func NewArticleController(rg *gin.RouterGroup, di service.DI) {
	c := &ArticleController{
		helper: &ArticleHelper{
			db: di.MariaDB,
		},
	}
	r := rg.Group("/article")
	r.Use(middleware.MustLoginMiddleware(di))
	r.GET("/fresh", c.fresh)
	r.GET("/list", c.list)
	r.POST("/add", c.add)
	r.POST("/del", c.del)
	r.POST("/renew", c.renew)
	r.GET("/edit/:id", c.edit)
}
