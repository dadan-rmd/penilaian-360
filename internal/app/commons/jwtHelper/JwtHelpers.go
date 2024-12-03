package jwtHelper

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	ErrTokenExpired = errors.New("token already expired")
)

var (
	//jwtSignatureKey is a secret key to hash the JWT Token
	jwtSignatureKey                  string
	prospectRefreshTokenSignatureKey string
)

// CustomClaims - Represent object of claims. Encouraged all claims is referred to this struct
type CustomClaims struct {
	jwt.StandardClaims
	Id        int64  `json:"jti,omitempty"`
	Email     string `json:"email,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	TokenUID  string `json:"tid,omitempty"`
}

func init() {
	err := godotenv.Load("params/.env")
	if err != nil {
		log.Error().Msg("Error loading .env file")
	}

	jwtSignatureKey = os.Getenv("JWT_SECRET_KEY")
	prospectRefreshTokenSignatureKey = os.Getenv("JWT_REFRESH_SECRET_KEY")
}

// NewWithClaims will return token with custom claims
func NewWithClaims(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSignatureKey))

	if err != nil {
		return "", err
	}
	return ss, nil
}

func NewWithClaimsRefreshToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(prospectRefreshTokenSignatureKey))

	if err != nil {
		return "", err
	}
	return ss, nil
}

// VerifyTokenWithClaims will verify the validity of token and return the claims
func VerifyTokenWithClaims(token string) (*CustomClaims, error) {

	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignatureKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("Error retrieving claims")
	}

	timeNow := time.Now().Unix()
	if claims.ExpiresAt < timeNow {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

// ExtractClaims will retrieve token claim but not verify expiration
func ExtractClaims(token string) (*CustomClaims, error) {

	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignatureKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("Error retrieving claims")
	}

	return claims, nil
}

// VerifyTokenWithApplicantClaims will verify the validity of token and return the claims
func VerifyTokenWithApplicantClaims(token string) (*CustomClaims, error) {

	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignatureKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("error retrieving claims")
	}

	timeNow := time.Now().Unix()
	if claims.ExpiresAt < timeNow {
		return nil, ErrTokenExpired
	}

	return claims, nil
}
