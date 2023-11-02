package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("123456")

type MyClaims struct {
	Uid                  string `json:"uid"`
	jwt.RegisteredClaims        // v5版本新加的方法
}

func GenToken(uid string) (string, error) {
	claims := MyClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 生效时间
		},
	}
	// 使用HS256签名算法
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(mySecret)

	return s, err
}

func ParseJwt(tokenstring string) (*MyClaims, error) {
	t, err := jwt.ParseWithClaims(tokenstring, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecret), nil
	})

	if claims, ok := t.Claims.(*MyClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}