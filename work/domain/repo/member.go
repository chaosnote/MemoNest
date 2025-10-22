package repo

import "idv/chris/MemoNest/domain/entity"

type MemberRepository interface {
	Login(account, password, last_ip string) (mo entity.Member, err error)
	Register(account, password, last_ip string) (mo entity.Member, err error)
}
