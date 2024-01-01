package auth

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   string `json:"ID"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID string `json:"ID"`
	jwt.RegisteredClaims
}
