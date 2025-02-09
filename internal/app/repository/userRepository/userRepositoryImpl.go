package userRepository

import (
	"penilaian-360/internal/app/model/userModel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (d userRepository) FindUserByID(id int64) (*userModel.User, error) {
	var userData userModel.User
	db := d.db.Preload(clause.Associations)
	db.First(&userData, "id=? and deleted_at = ?", id, 0)
	return &userData, db.Error
}

func (d userRepository) FindUserByEmail(email string) (userData *userModel.User, err error) {
	err = d.db.First(&userData, "deleted_at = ? and email = ?", 0, email).Error
	return
}

func (d userRepository) Save(tx *gorm.DB, data *userModel.User) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d userRepository) Delete(userData userModel.User) error {
	db := d.db.Delete(&userData)
	return db.Error
}
