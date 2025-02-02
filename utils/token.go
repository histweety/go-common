package utils

import (
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/histweety/go-common/errors"
	"github.com/histweety/go-common/types"
)

var jwtKey = os.Getenv("SECRET_KEY")
var refreshJwtKey = os.Getenv("SECRET_REFRESH_KEY")

func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &types.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &types.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshJwtKey)
}

func ParseToken(tokenString string) (*types.Claims, error) {
	claims := &types.Claims{}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.ErrUnauthorized
		}

		return nil, errors.ErrBadRequest
	}

	if !token.Valid {
		return nil, errors.ErrUnauthorized
	}

	return claims, nil
}
