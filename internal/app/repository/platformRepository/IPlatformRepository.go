package platformRepository

import (
	"central-auth/internal/app/model/platformModel"

	"gorm.io/gorm"
)

type IPlatformRepository interface {
	BulkInsert(tx *gorm.DB, data []platformModel.Platform) error
}
