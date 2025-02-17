package jwtHelper

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrInvalidToken = errors.New("invalid token")
)

var jwtSecret string

// Claims - Struct untuk menyimpan payload JWT
type Claims struct {
	Data struct {
		Email       string `json:"email"`
		AccessToken string `json:"accesstoken"`
	} `json:"data"`
	jwt.RegisteredClaims
}

func EncodeJWT(email, accessToken string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", ErrInvalidToken
	}

	claims := Claims{
		Data: struct {
			Email       string `json:"email"`
			AccessToken string `json:"accesstoken"`
		}{
			Email:       email,
			AccessToken: accessToken,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func DecodeJWT(tokenString string) (*Claims, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		fmt.Println("step 1->", err)
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		fmt.Println("step 2->", err)
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		fmt.Println("step 3->", err)
		return nil, ErrTokenExpired
	}

	return claims, nil
}
