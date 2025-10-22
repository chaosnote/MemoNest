package usecase

import (
	"fmt"
	"regexp"

	"idv/chris/MemoNest/domain/entity"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/utils"
)

type MemberUsecase struct {
	Repo repo.MemberRepository
}

func (u *MemberUsecase) check(account, password, last_ip string) (err error) {
	var account_regex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{4,9}$`)
	var password_regex = regexp.MustCompile(`^[a-zA-Z0-9]{6,10}$`)

	if !account_regex.MatchString(account) {
		err = fmt.Errorf("帳號格式錯誤")
		return
	}
	if !password_regex.MatchString(password) {
		err = fmt.Errorf("密碼格式錯誤")
		return
	}
	return
}

func (u *MemberUsecase) Login(account, password, last_ip string) (mo entity.Member, err error) {
	err = u.check(account, password, last_ip)
	if err != nil {
		return
	}

	aes_key := utils.GenAESKey([]byte(account))
	password, err = utils.AesEncrypt([]byte(password), aes_key)
	if err != nil {
		return
	}
	mo, err = u.Repo.Login(account, password, last_ip)
	if err != nil {
		return
	}

	return
}

func (u *MemberUsecase) Register(account, password, last_ip string) (mo entity.Member, err error) {
	err = u.check(account, password, last_ip)
	if err != nil {
		return
	}

	aes_key := utils.GenAESKey([]byte(account))
	password, err = utils.AesEncrypt([]byte(password), aes_key)
	if err != nil {
		return
	}
	mo, err = u.Repo.Register(account, password, last_ip)
	if err != nil {
		return
	}

	return
}

//-----------------------------------------------

func NewMemberUsecase(
	repo repo.MemberRepository,
) *MemberUsecase {
	return &MemberUsecase{
		Repo: repo,
	}
}
