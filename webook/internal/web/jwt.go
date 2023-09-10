package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Token interface {
	SetJWTToken(context *gin.Context, uid int64) error
}

type JWTToken struct {
}

func NewJwtToken() Token {
	return &JWTToken{}
}
func (j JWTToken) SetJWTToken(context *gin.Context, uid int64) error {
	key := []byte("4B2B17F68975BA8C806846C5CC898070F3F62423BEC9D87F2DBE844B4C14F137")
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &UserCliams{
		Id: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UserAgent: context.Request.UserAgent(),
	})
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return err
	}
	context.Header("x-jwt-token", tokenStr)
	return nil
}
