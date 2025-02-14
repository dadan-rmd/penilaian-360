package evaluatedEmployeeRepository

import (
	"penilaian-360/internal/app/model/evaluatedEmployeesModel"

	"gorm.io/gorm"
)

type IEvaluatedEmployeeRepository interface {
	Save(tx *gorm.DB, data *[]evaluatedEmployeesModel.EvaluatedEmployee) error
	FindEmployeeIdByEvaluationId(evaluationId int64) (employeeId []int64, err error)
}
