package middleware

import (
	"boot-camp-homework/homework/week2/webhook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 不需要登录校验
		for _, path := range l.paths {
			if context.Request.URL.Path == path {
				return
			}
		}
		// 我现在用JWT来校验
		tokenHeader := context.GetHeader("x-jwt-token")
		if tokenHeader == "" {
			//没登录
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 1 {
			//有人搞事情
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[0]
		claims := &web.UserCliams{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("4B2B17F68975BA8C806846C5CC898070F3F62423BEC9D87F2DBE844B4C14F137"), nil
		})
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//err为nil, token 不为nil
		if token == nil || !token.Valid || claims.Id == 0 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != context.Request.UserAgent() {
			//严重的安全问题
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		key := []byte("4B2B17F68975BA8C806846C5CC898070F3F62423BEC9D87F2DBE844B4C14F137")
		if claims.RegisteredClaims.ExpiresAt.Sub(now) < 50*time.Minute {
			claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Hour))
			tokenStr, err = token.SignedString(key)
			if err != nil {
				log.Panicln("续约失败", err)
			}
			context.Header("x-jwt-token", tokenStr)
		}
		context.Set("cliams", claims)
	}
}
