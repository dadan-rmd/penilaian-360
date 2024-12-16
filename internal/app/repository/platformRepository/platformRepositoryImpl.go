package platformRepository

import (
	"central-auth/internal/app/model/platformModel"

	"gorm.io/gorm"
)

type platformRepository struct {
	db *gorm.DB
}

func NewPlatformRepository(db *gorm.DB) IPlatformRepository {
	return &platformRepository{db}
}

func (d platformRepository) BulkInsert(tx *gorm.DB, data []platformModel.Platform) error {
	if tx != nil {
		return tx.Create(&data).Error
	} else {
		return d.db.Create(&data).Error
	}
}

func (d platformRepository) FindNameByID(id int64) (name []string, err error) {
	err = d.db.
		Model(platformModel.Platform{}).
		Where("id=?", id).
		Pluck("name", &name).
		Error
	return
}
