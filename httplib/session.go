package httplib

/**
session
使用方法：
中间件
   Session{IdKey:"TOT-SESSION",MaxAge:300,Writer:http.ResponseWriter,Request:*http.Request}.Use()
使用时：
   Session{Request:*http.Request}.Get();
*/
import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"go-nest/cache"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	SESSION_ID_KEY     string = "PIA-TOTANT"
	SESSION_REDIS_HOST string = "localhost:6379" //redis服务器
	SESSION_REDIS_PSWD string = ""               //redis密码
	SESSION_REDIS_DB   int    = 1                //redis数据库
)

type Session struct {
	IdKey      string
	MaxAge     int64  //生存时间，单位：秒
	RdsHost    string //redis服务器
	RdsPaswrod string //redis密码
	RdsDB      int    //redis数据库编号
	cache      *cache.Cache
	Writer     http.ResponseWriter
	Request    *http.Request
}

//创建session时使用
func (self Session) Use() *cache.Cache {
	o := &self
	sid := o.getCookieSessionId()
	if sid == "" {
		sid = o.createSessionId()
	}
	if o.IdKey != "" {
		SESSION_ID_KEY = o.IdKey
	}
	if o.RdsHost != "" {
		SESSION_REDIS_HOST = o.RdsHost
	}
	if o.RdsPaswrod != "" {
		SESSION_REDIS_PSWD = o.RdsPaswrod
	}
	if o.RdsDB != 0 {
		SESSION_REDIS_DB = o.RdsDB
	}
	o.initCache(sid)
	o.cache.Set(SESSION_ID_KEY, sid)
	o.cache.Expire(o.MaxAge)
	return o.cache
}

//使用session时使用
func (self Session) Get(r *http.Request) *cache.Cache {
	o := &self
	o.Request = r
	sid := o.getCookieSessionId()
	o.initCache(sid)
	return o.cache
}

//从cookie里面读取session id
func (this *Session) getCookieSessionId() string {
	cookie, err := this.Request.Cookie(SESSION_ID_KEY)
	if err == nil && cookie.Value != "" {
		sid, e := url.QueryUnescape(cookie.Value)
		if e == nil && sid != "" {
			return sid
		}
	}
	return ""
}

//初始化缓存数据库
func (this *Session) initCache(sid string) {
	this.cache = cache.New(cache.Redis(SESSION_REDIS_HOST, SESSION_REDIS_PSWD, SESSION_REDIS_DB, sid))
}

//创建session id
func (this *Session) createSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	key := time.Now().String() + string(b) //session 的格式= 当前时间+随机数
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(key))
	cipherStr := md5Ctx.Sum(nil)
	sid := hex.EncodeToString(cipherStr)
	cookie := http.Cookie{Name: SESSION_ID_KEY, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(this.MaxAge)}
	http.SetCookie(this.Writer, &cookie)
	return sid
}
