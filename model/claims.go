package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.StandardClaims
}
