package platformRepository

import (
	"central-auth/internal/app/model/platformModel"

	"gorm.io/gorm"
)

type IPlatformRepository interface {
	BulkInsert(tx *gorm.DB, data []platformModel.Platform) error
	FindNameByID(id int64) (name []string, err error)
	FindNameByUserID(id int64) (name []string, err error)
}
