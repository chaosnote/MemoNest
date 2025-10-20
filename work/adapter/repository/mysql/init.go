package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"idv/chris/MemoNest/config"
)

// NewMariaDB 建立 MariaDB 連線
func NewMariaDB(cfg *config.APPConfig) (*sql.DB, error) {
	db, e := sql.Open("mysql", cfg.Mariadb.DSN)
	if e != nil {
		return nil, e
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if e = db.Ping(); e != nil {
		return nil, e
	}

	return db, nil
}
