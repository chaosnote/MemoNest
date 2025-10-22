package mysql

import (
	"database/sql"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
)

type MemberRepo struct {
	db *sql.DB
}

func (r *MemberRepo) Login(account, password string, last_ip string) (mo entity.Member, err error) {
	row := r.db.QueryRow(
		"CALL `sp_login`(?,?,?)",
		account,
		password,
		last_ip,
	)
	err = row.Scan(
		&mo.RowID,
		&mo.Account,
		&mo.Password,
		&mo.LastIP,
		&mo.IsEnabled,
		&mo.CreatedAt,
		&mo.UpdatedAt,
	)
	if err != nil {
		return
	}
	return
}

func (r *MemberRepo) Register(account, password string, last_ip string) (mo entity.Member, err error) {
	row := r.db.QueryRow(
		"CALL `sp_add_member`(?,?,?)",
		account,
		password,
		last_ip,
	)
	err = row.Scan(
		&mo.RowID,
		&mo.Account,
		&mo.Password,
		&mo.LastIP,
		&mo.IsEnabled,
		&mo.CreatedAt,
		&mo.UpdatedAt,
	)
	if err != nil {
		return
	}

	return
}

//-----------------------------------------------

func NewMemberRepo(db *sql.DB) repo.MemberRepository {
	return &MemberRepo{db: db}
}
