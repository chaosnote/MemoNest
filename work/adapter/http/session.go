package http

import (
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SK_ACCOUNT  = "account"
	SK_AES_KEY  = "aes_key"
	SK_IS_LOGIN = "is_login"
	SK_IP       = "ip"
	SK_URL      = "url"
)

type GinSession struct {
	store sessions.Session
}

func (s *GinSession) Init(ctx *gin.Context) {
	s.store = sessions.Default(ctx)
}

func (s *GinSession) Clear() {
	s.store.Clear()
	s.store.Save()
}

func (s *GinSession) SetAccount(account, last_ip string) {
	s.store.Set(SK_ACCOUNT, account)
	s.store.Set(SK_AES_KEY, utils.MD5Hash(account)) // 可改為 row_id
	s.store.Set(SK_IS_LOGIN, true)
	s.store.Set(SK_IP, last_ip)
	s.store.Save()
}

func (s *GinSession) GetAccount() string {
	val := s.store.Get(SK_ACCOUNT)
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (s *GinSession) SetURL(path string) {
	s.store.Set(SK_URL, path)
	s.store.Save()
}

func (s *GinSession) GetURL() string {
	val := s.store.Get(SK_URL)
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (s *GinSession) GetAESKey() string {
	val := s.store.Get(SK_AES_KEY)
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (s *GinSession) GetIP() string {
	val := s.store.Get(SK_IP)
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (s *GinSession) IsLogin() bool {
	val := s.store.Get(SK_IS_LOGIN)
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

func (s *GinSession) Refresh() {
	s.store.Set("_refresh", time.Now().Unix())
	s.store.Save()
}

//-----------------------------------------------

func NewGinSession() service.Session {
	return &GinSession{}
}
