package evaluatedEmployeeRepository

import (
	"penilaian-360/internal/app/model/evaluatedEmployeesModel"

	"gorm.io/gorm"
)

type evaluatedEmployeeRepository struct {
	db *gorm.DB
}

func NewEvaluatedEmployeeRepository(db *gorm.DB) IEvaluatedEmployeeRepository {
	return &evaluatedEmployeeRepository{db}
}

func (d evaluatedEmployeeRepository) Save(tx *gorm.DB, data *[]evaluatedEmployeesModel.EvaluatedEmployee) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d evaluatedEmployeeRepository) FindEmployeeIdByEvaluationId(evaluationId int64) (employeeId []int64, err error) {
	err = d.db.
		Model(&evaluatedEmployeesModel.EvaluatedEmployee{}).
		Where(evaluatedEmployeesModel.EvaluatedEmployee{
			EvaluationId: evaluationId,
		}).
		Pluck("employee_id", &employeeId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []int64{}, nil
		}
		return
	}
	return
}
