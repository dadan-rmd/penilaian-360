package userRepository

import (
	"penilaian-360/internal/app/model/userModel"

	"gorm.io/gorm"
)

type IUserRepository interface {
	FindUserByID(id int64) (data *userModel.User, err error)
	FindUserByEmail(email string) (userData *userModel.User, err error)
	Save(tx *gorm.DB, data *userModel.User) error
	Delete(userData userModel.User) error
}
