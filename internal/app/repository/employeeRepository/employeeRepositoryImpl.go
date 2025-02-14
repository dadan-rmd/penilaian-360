package employeeRepository

import (
	"penilaian-360/internal/app/model/employeeModel"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) IEmployeeRepository {
	return &employeeRepository{db}
}

func (d employeeRepository) FindByDepartement(departement string, ids []int64) (entities []employeeModel.Employee, err error) {
	db := d.db.Model(&employeeModel.Employee{})
	if len(ids) > 0 {
		db.Where("id NOT IN ?", ids)
	}
	err = db.Where(employeeModel.Employee{Department: departement}).
		Find(&entities).Error
	return
}

func (d employeeRepository) FindByEmailAndAccessToken(email, accessToken string) (entity employeeModel.Employee, err error) {
	err = d.db.Where(employeeModel.Employee{Email: email, AccessToken: accessToken}).
		First(&entity).Error
	return
}

func (d employeeRepository) FindByIds(ids []int64) (entity []employeeModel.Employee, err error) {
	err = d.db.Find(&entity, ids).Error
	return
}
