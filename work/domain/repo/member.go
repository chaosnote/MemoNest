package repo

import "idv/chris/MemoNest/domain/entity"

type MemberRepository interface {
	Register(*entity.Member) (*entity.Member, error)
}
