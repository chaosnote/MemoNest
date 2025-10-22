package mysql

import (
	"database/sql"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
)

type MemberRepo struct {
	db *sql.DB
}

func (r *MemberRepo) Register(src *entity.Member) (mo *entity.Member, err error) {
	return
}

//-----------------------------------------------

func NewMemberRepo(db *sql.DB) repo.MemberRepository {
	return &MemberRepo{db: db}
}
