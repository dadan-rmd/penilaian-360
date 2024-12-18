package userService

import (
	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/commons/symmetricHash"
	"central-auth/internal/app/model/platformModel"
	"central-auth/internal/app/model/userModel"
	"central-auth/internal/app/repository/platformRepository"
	"central-auth/internal/app/repository/userRepository"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExist    = errors.New("email sudah terdaftar")
	ErrUserNotFound         = errors.New("user tidak ditemukan")
	ErrorDuplicatePhone     = errors.New("nomor telepon sudah terdaftar")
	ErrorInvalidFormatEmail = errors.New("format email tidak valid")
	ErrorInvalidFormatPhone = errors.New("format nomor telepon tidak valid")
)

type userService struct {
	userRepo     userRepository.IUserRepository
	db           *gorm.DB
	platformRepo platformRepository.IPlatformRepository
}

func NewUserService(userRepo userRepository.IUserRepository, db *gorm.DB, platformRepo platformRepository.IPlatformRepository) IUserService {
	return &userService{userRepo, db, platformRepo}
}

func (u userService) CreateUser(record *loggers.Data, req userModel.CreateUserReq) (user *userModel.User, statusCode int, err error) {
	var platform []platformModel.Platform
	user, err = u.userRepo.FindUserByEmail(req.Email)
	tx := u.db.Begin()
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &userModel.User{
				Username: req.Username,
				Email:    req.Email,
				Password: symmetricHash.ToBcrypt(req.Password),
			}
			err = u.userRepo.Save(tx, user)
			if err != nil {
				loggers.Logf(record, "[Error] Save User : ", err)
				statusCode = http.StatusInternalServerError
				return
			}
			for _, v := range req.Platform {
				platform = append(platform, platformModel.Platform{
					UserId: user.Id,
					Name:   strings.ToUpper(v),
				})
			}
			err = u.platformRepo.BulkInsert(tx, platform)
			if err != nil {
				loggers.Logf(record, "[Error] BulkInsert Platform : ", err)
				statusCode = http.StatusInternalServerError
				return
			}
			statusCode = http.StatusOK
			return
		}
		return nil, http.StatusInternalServerError, ErrEmailAlreadyExist
	}
	namePlatform, err := u.platformRepo.FindNameByUserID(user.Id)
	if err != nil {
		loggers.Logf(record, "[Error] FindNameByUserID : ", err)
		statusCode = http.StatusInternalServerError
		return
	}
	platformMap := make(map[string]bool)
	for _, data := range namePlatform {
		platformMap[data] = true
	}
	for _, v := range req.Platform {
		if !platformMap[v] {
			platform = append(platform, platformModel.Platform{
				UserId: user.Id,
				Name:   strings.ToUpper(v),
			})
		}
	}
	err = u.platformRepo.BulkInsert(tx, platform)
	if err != nil {
		loggers.Logf(record, "[Error] BulkInsert Platform : ", err)
		statusCode = http.StatusInternalServerError
		return
	}
	statusCode = http.StatusOK
	return
}
