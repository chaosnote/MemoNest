package http

import (
	"fmt"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
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
	conn  redigo.Conn
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

	s.conn.Do("SET", account, last_ip)
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

func (s *GinSession) Refresh() error {
	mem_ip, err := redigo.String(s.conn.Do("GET", s.GetAccount()))
	if err != nil {
		return err
	}
	my_ip := s.GetIP()
	if my_ip != mem_ip {
		return fmt.Errorf("來源IP[%s]不符", my_ip)
	}

	s.store.Set("_refresh", time.Now().Unix())
	s.store.Save()

	return nil
}

//-----------------------------------------------

func NewGinSession(redis_store redis.Store) service.Session {
	store, err := redis.GetRedisStore(redis_store)
	if err != nil {
		panic(err)
	}

	return &GinSession{
		conn: store.Pool.Get(),
	}
}
