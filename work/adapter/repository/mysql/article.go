package mysql

import (
	"database/sql"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/model"
	"strings"
	"time"
)

type ArticleRepo struct {
	db *sql.DB
}

func (ah *ArticleRepo) GetAllNode() (categories []model.Category, err error) {
	rows, err := ah.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}

	return
}

func (ah *ArticleRepo) Add(node_id string) (id int, e error) {
	t := time.Now().UTC()
	row := ah.db.QueryRow("CALL insert_article(?, ?, ?, ?, ?) ;", "", "", t, t, node_id)
	e = row.Scan(&id)
	if e != nil {
		return
	}
	return
}

func (ah *ArticleRepo) Delete(id int) (e error) {
	_, e = ah.db.Exec(`DELETE from articles where RowID = ? ;`, id)
	if e != nil {
		return
	}
	return
}

func (ah *ArticleRepo) Update(row_id int, title, content string) error {
	t := time.Now().UTC()
	query := `UPDATE articles SET Title = ?, Content = ?, UpdateDt =? WHERE RowID = ?;`
	_, e := ah.db.Exec(query, title, content, t, row_id)
	if e != nil {
		return e
	}
	return nil
}

func (ah *ArticleRepo) Get(id int) (articles []model.Article, err error) {
	rows, err := ah.db.Query(`
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
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		err = rows.Scan(
			&article.RowID,
			&article.Title,
			&article.Content,
			&article.NodeID,
			&article.PathName,
			&article.UpdateDt,
		)
		if err != nil {
			return
		}
		articles = append(articles, article)
	}
	return
}

func (ah *ArticleRepo) List() (articles []model.Article, err error) {
	rows, err := ah.db.Query(`
		SELECT 
			a.RowID AS ArticleRowID,
			a.Title,
			a.Content,
			a.NodeID,
			c.PathName,
			a.UpdateDt
		FROM articles as a
		JOIN categories as c ON a.NodeID = c.NodeID
		ORDER BY a.UpdateDt DESC
		LIMIT 10;
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		err = rows.Scan(
			&article.RowID,
			&article.Title,
			&article.Content,
			&article.NodeID,
			&article.PathName,
			&article.UpdateDt,
		)
		if err != nil {
			return
		}
		articles = append(articles, article)
	}
	return
}

func (ah *ArticleRepo) composit(input string) (query string, args []interface{}) {
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

func (ah *ArticleRepo) Query(input string) (articles []model.Article, err error) {
	cmd, args := ah.composit(input)
	rows, err := ah.db.Query(cmd, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		err = rows.Scan(
			&article.RowID,
			&article.Title,
			&article.Content,
			&article.NodeID,
			&article.PathName,
			&article.UpdateDt,
		)
		if err != nil {
			return
		}
		articles = append(articles, article)
	}
	return
}

func NewArticleRepo(db *sql.DB) repo.ArticleRepository {
	return &ArticleRepo{db: db}
}
