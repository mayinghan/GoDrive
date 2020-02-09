package utils

import (
	"GoDrive/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// jwtKey is the key used to create the signature
var jwtKey = []byte("myhisaqt")

// Claims is a struct that is encoded to a jwt
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Gentoken generate a jwt token by using the user's username
func Gentoken(username string) (string, error) {
	expTime := time.Now().Add(config.JWTLife)
	claims := &Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			Issuer:    "godrivedev",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)

	return tokenStr, err
}

// ParseToken parse the given JWT and returns a Claim struct which contains username
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
