package departmentRepository

import (
	"penilaian-360/internal/app/model/departmentModel"

	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) IDepartmentRepository {
	return &departmentRepository{db}
}

func (d departmentRepository) FindAll() (entities []departmentModel.Department, err error) {
	err = d.db.Find(&entities).Error
	return
}
