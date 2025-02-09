package authService

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredential = errors.New("invalid credential")
	ErrUserNotFound      = errors.New("user not found")
	ErrPlatformNotFound  = errors.New("platform not found")
)

const (
	MaxOTPVerificationLifetime = 30 * time.Minute
	OTPLength                  = 6
)

type authService struct {
	db *gorm.DB
}

func NewAuthService(
	db *gorm.DB,
) IAuthService {
	return &authService{db}
}
