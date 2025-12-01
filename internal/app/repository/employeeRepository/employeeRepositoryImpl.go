package employeeRepository

import (
	"fmt"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/model/evaluatedEmployeesModel"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) IEmployeeRepository {
	return &employeeRepository{db}
}

func (d employeeRepository) FindByDepartement(departement string, ids []int64, hasAssigned bool) (entities []employeeModel.Employee, err error) {
	db := d.db.Model(&employeeModel.Employee{})
	if len(ids) > 0 {
		if hasAssigned {
			db.Where("id IN ?", ids)
		} else {
			db.Where("id NOT IN ?", ids)
		}
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

func (d employeeRepository) FindNameAndEmployedIDByIds(ids []int64) (entities []employeeModel.EmployedEmployeeResponse, err error) {
	err = d.db.Model(evaluatedEmployeesModel.EvaluatedEmployee{}).
		Select("master_karyawan.Name,evaluated_employees.id as evaluated_id").
		Joins("JOIN master_karyawan on master_karyawan.id = evaluated_employees.employee_id").
		Where("master_karyawan.id in (?)", ids).
		Find(&entities).Error
	return
}

func (d employeeRepository) FindEmailsByKeyword(keyword string) (emails []string, err error) {
	if len(keyword) < 3 {
		return emails, fmt.Errorf("keyword must be at least 3 characters")
	}

	var results []employeeModel.Employee

	err = d.db.
		Select("Email").
		Where("Email LIKE ? OR Name LIKE ? OR FirstName LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&results).Error

	if err != nil {
		return emails, err
	}

	for _, r := range results {
		emails = append(emails, r.Email)
	}

	return emails, nil
}
