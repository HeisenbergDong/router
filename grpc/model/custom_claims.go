package model

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UUID        string
	ID          uint64
	UserName    string
	NickName    string
	AuthorityId string
	BufferTime  int64
	jwt.StandardClaims
}
