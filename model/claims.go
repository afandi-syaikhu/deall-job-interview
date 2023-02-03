package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}
