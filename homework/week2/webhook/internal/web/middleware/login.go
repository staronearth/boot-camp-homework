package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(context *gin.Context) {
		// 不需要登录校验
		for _, path := range l.paths {
			if context.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(context)
		id := sess.Get("UserId")
		if id == nil {
			//没有登录
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//updateTime := sess.Get("update_time")
		//sess.Set("userId", id)
		//sess.Options(sessions.Options{
		//	MaxAge: 60,
		//})
		//now := time.Now()
		//// 说明还没有刷新过，刚登陆，还没刷新过
		//if updateTime == nil {
		//	sess.Set("update_time", now)
		//	if err := sess.Save(); err != nil {
		//		panic(err)
		//	}
		//}
		//// updateTime 是有的
		//updateTimeVal, _ := updateTime.(time.Time)
		//if now.Sub(updateTimeVal) > time.Second*10 {
		//	sess.Set("update_time", now)
		//	if err := sess.Save(); err != nil {
		//		panic(err)
		//	}
		//}

		currenttime := time.Now()
		updatetime := sess.Get("updatetime")
		sess.Set("UserId", id)
		sess.Options(sessions.Options{
			MaxAge: 180,
		})
		if updatetime == nil {
			sess.Set("updatetime", currenttime)
			sess.Save()
			return
		}
		if currenttime.Sub(updatetime.(time.Time)) > time.Minute {
			sess.Set("updatetime", currenttime)
			sess.Save()
		}
	}
}
