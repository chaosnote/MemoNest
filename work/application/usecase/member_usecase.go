package usecase

type MemberUsecase struct {
}

func (u *MemberUsecase) Login(account, password string) bool {
	return true
}

func (u *MemberUsecase) Register(account, password string) bool {
	return true
}

//-----------------------------------------------

func NewMemberUsecase() *MemberUsecase {
	return &MemberUsecase{}
}
