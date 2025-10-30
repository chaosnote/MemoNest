package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
)

type ArticleRepo struct {
	db                 *sql.DB
	node_formatter     string
	articles_formatter string
}

func (r *ArticleRepo) GetAllNode(account string) (categories []entity.Category, err error) {
	query := `SELECT RowID, NodeID, ParentID, PathName, LftIdx, RftIdx FROM %s ORDER BY LftIdx ASC`
	query = fmt.Sprintf(query, fmt.Sprintf(r.node_formatter, account))

	rows, err := r.db.Query(query)
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

func (r *ArticleRepo) Add(account, node_id string) (id int, e error) {
	t := time.Now().UTC()

	row := r.db.QueryRow("CALL sp_add_article(?, ?, ?, ?, ?, ?) ;", account, "", "", t, t, node_id)
	e = row.Scan(&id)
	if e != nil {
		return
	}
	return
}

func (r *ArticleRepo) Delete(account string, id int) (e error) {
	query := `DELETE from %s where RowID = ? ;`
	query = fmt.Sprintf(query, fmt.Sprintf(r.articles_formatter, account))
	_, e = r.db.Exec(query, id)
	if e != nil {
		return
	}
	return
}

func (r *ArticleRepo) Update(account string, row_id int, title, content string) error {
	t := time.Now().UTC()
	query := `UPDATE %s SET Title = ?, Content = ?, UpdateDt =? WHERE RowID = ?;`
	query = fmt.Sprintf(query, fmt.Sprintf(r.articles_formatter, account))
	_, e := r.db.Exec(query, title, content, t, row_id)
	if e != nil {
		return e
	}
	return nil
}

func (r *ArticleRepo) Get(account string, id int) (articles []entity.Article, err error) {
	query := `
		SELECT 
			a.RowID AS ArticleRowID,
			a.Title,
			a.Content,
			a.NodeID,
			c.PathName,
			a.UpdateDt
		FROM %s as a
		JOIN %s as c ON a.NodeID = c.NodeID
		WHERE a.RowID = ? ;
	`
	query = fmt.Sprintf(query, fmt.Sprintf(r.articles_formatter, account), fmt.Sprintf(r.node_formatter, account))

	rows, err := r.db.Query(query, id)
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

func (r *ArticleRepo) List(account string) (articles []entity.Article, err error) {
	query := `
		SELECT 
			a.RowID AS ArticleRowID,
			a.Title,
			a.Content,
			a.NodeID,
			c.PathName,
			a.UpdateDt
		FROM %s as a
		JOIN %s as c ON a.NodeID = c.NodeID
		ORDER BY a.UpdateDt DESC
		LIMIT 10;
	`
	query = fmt.Sprintf(query, fmt.Sprintf(r.articles_formatter, account), fmt.Sprintf(r.node_formatter, account))

	rows, err := r.db.Query(query)
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

func (r *ArticleRepo) composit(account, input string) (query string, args []interface{}) {
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
        FROM %s AS a
        JOIN %s AS c ON a.NodeID = c.NodeID
        WHERE ` + strings.Join(conditions, " AND ") + `
        ORDER BY a.UpdateDt DESC
    `
	query = fmt.Sprintf(query, fmt.Sprintf(r.articles_formatter, account), fmt.Sprintf(r.node_formatter, account))
	return
}

func (r *ArticleRepo) Query(account, input string) (articles []entity.Article, err error) {
	cmd, args := r.composit(account, input)
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
	return &ArticleRepo{
		db:                 db,
		node_formatter:     "node_%s",
		articles_formatter: "articles_%s",
	}
}
