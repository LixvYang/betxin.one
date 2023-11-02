package jwt

import (
	"fmt"
	"os"
	"testing"
)

func Test_JWTToken(t *testing.T) {
	s, err := GenToken("zhangsan")
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%s\n", s)

	// 解析jwt
	claims, err := ParseJwt(s)
	if err != nil {
		fmt.Println("parse jwt failed, ", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", claims)
}
