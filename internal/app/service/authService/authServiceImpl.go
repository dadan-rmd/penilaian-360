package authService

import (
	"central-auth/internal/app/commons/jwtHelper"
	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/commons/symmetricHash"
	"central-auth/internal/app/model/authModel"
	"central-auth/internal/app/model/userModel"
	"central-auth/internal/app/repository/userRepository"
	"errors"
	"os"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
)

var (
	ErrInvalidCredential = errors.New("invalid credential")
	ErrUserNotFound      = errors.New("user not found")
)

const (
	MaxOTPVerificationLifetime = 30 * time.Minute
	OTPLength                  = 6
)

type authUseCase struct {
	userRepo userRepository.IUserRepository
}

func NewAuthService(
	userRepo userRepository.IUserRepository,
) IAuthService {
	return &authUseCase{userRepo}
}

func (a authUseCase) Login(record *loggers.Data, loginReq authModel.LoginReq) (loginRes userModel.ResLogin, err error) {
	userData, err := a.userRepo.FindUserByEmail(loginReq.Email)
	if err != nil {
		return loginRes, ErrUserNotFound
	}

	if !symmetricHash.CompareBcrypt(userData.Password, loginReq.Password) {
		return loginRes, ErrInvalidCredential
	}

	loginRes.User = *userData

	jwtExpirationDurationDayString := os.Getenv("JWT_EXPIRATION_DURATION_DAY")
	var jwtExpirationDurationDay int
	jwtExpirationDurationDay, err = strconv.Atoi(jwtExpirationDurationDayString)
	if err != nil {
		return loginRes, err
	}

	jwtExpiredAt := time.Now().Unix() + int64(jwtExpirationDurationDay*3600*24)
	tokenUID := uuid.NewV4().String() + "00" + cast.ToString(userData.Id)

	userClaims := jwtHelper.CustomClaims{Email: userData.Email, ExpiresAt: jwtExpiredAt, TokenUID: tokenUID}
	jwtToken, err := jwtHelper.NewWithClaims(userClaims)
	if err != nil {
		return loginRes, err
	}
	loginRes.Token = jwtToken
	return loginRes, nil
}
