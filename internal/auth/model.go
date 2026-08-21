package auth

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id       string
	Name     string
	Password string
}

type UserInfo struct {
	Name string
}

type UserClaims struct {
	jwt.RegisteredClaims
	TokenType string
	UserInfo
}
