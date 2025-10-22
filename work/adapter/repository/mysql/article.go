package mysql

import (
	"database/sql"
	"strings"
	"time"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
)

type ArticleRepo struct {
	db *sql.DB
}

func (r *ArticleRepo) GetAllNode() (categories []entity.Category, err error) {
	rows, err := r.db.Query("SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM categories ORDER BY LftIdx ASC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.RowID, &c.NodeID, &c.ParentID, &c.PathName, &c.LftIdx, &c.RftIdx); err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}

	return
}

func (r *ArticleRepo) Add(node_id string) (id int, e error) {
	t := time.Now().UTC()
	row := r.db.QueryRow("CALL insert_article(?, ?, ?, ?, ?) ;", "", "", t, t, node_id)
	e = row.Scan(&id)
	if e != nil {
		return
	}
	return
}

func (r *ArticleRepo) Delete(id int) (e error) {
	_, e = r.db.Exec(`DELETE from articles where RowID = ? ;`, id)
	if e != nil {
		return
	}
	return
}

func (r *ArticleRepo) Update(row_id int, title, content string) error {
	t := time.Now().UTC()
	query := `UPDATE articles SET Title = ?, Content = ?, UpdateDt =? WHERE RowID = ?;`
	_, e := r.db.Exec(query, title, content, t, row_id)
	if e != nil {
		return e
	}
	return nil
}

func (r *ArticleRepo) Get(id int) (articles []entity.Article, err error) {
	rows, err := r.db.Query(`
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
		var article entity.Article
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

func (r *ArticleRepo) List() (articles []entity.Article, err error) {
	rows, err := r.db.Query(`
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
		var article entity.Article
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

func (r *ArticleRepo) composit(input string) (query string, args []interface{}) {
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

func (r *ArticleRepo) Query(input string) (articles []entity.Article, err error) {
	cmd, args := r.composit(input)
	rows, err := r.db.Query(cmd, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article entity.Article
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
