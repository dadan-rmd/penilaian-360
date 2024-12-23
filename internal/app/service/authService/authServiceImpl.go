package authService

import (
	"central-auth/internal/app/commons/jwtHelper"
	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/commons/symmetricHash"
	"central-auth/internal/app/model/authModel"
	"central-auth/internal/app/model/userModel"
	"central-auth/internal/app/repository/platformRepository"
	"central-auth/internal/app/repository/userRepository"
	"errors"
	"os"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
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

type authUseCase struct {
	userRepo     userRepository.IUserRepository
	platformRepo platformRepository.IPlatformRepository
	db           *gorm.DB
}

func NewAuthService(
	userRepo userRepository.IUserRepository,
	platformRepo platformRepository.IPlatformRepository,
	db *gorm.DB,
) IAuthService {
	return &authUseCase{userRepo, platformRepo, db}
}

func (a authUseCase) Login(record *loggers.Data, loginReq authModel.LoginReq) (loginRes userModel.ResLogin, err error) {
	userData, err := a.userRepo.FindUserByEmail(loginReq.Email)
	if err != nil {
		return loginRes, ErrUserNotFound
	}

	if !symmetricHash.CompareBcrypt(userData.Password, loginReq.Password) {
		return loginRes, ErrInvalidCredential
	}
	platformName, err := a.platformRepo.FindNameByID(userData.Id)
	if err != nil {
		return loginRes, ErrPlatformNotFound
	}

	loginRes.User = *userData
	loginRes.Platform = platformName

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

func (a authUseCase) ChangePasswordFromForgotPass(record *loggers.Data, request authModel.ChangePasswordFromForgotPassReq) (err error) {
	loggers.Logf(record, "Info, VerifyForgotPassword")
	userData, err := a.userRepo.FindUserByEmail(request.Email)
	if err != nil {
		loggers.Logf(record, "Error, FindUserByEmail")
		return ErrUserNotFound
	}
	tx := a.db.Begin()
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				loggers.Logf(record, "[Error] Rollback: ", errRollback.Error)
			}
			return
		}
		if err = tx.Commit().Error; err != nil {
			loggers.Logf(record, "[Error] Commit : ", err)
		}
	}()
	userData.Password = symmetricHash.ToBcrypt(request.NewPassword)
	err = a.userRepo.Save(tx, userData)
	if err != nil {
		loggers.Logf(record, "Error, UpdateUser")
		return err
	}

	return nil
}
