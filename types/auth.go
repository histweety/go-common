package types

import "github.com/dgrijalva/jwt-go"

type ConfigAuth struct {
	SecretKey        string
	SecretRefreshKey string
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
