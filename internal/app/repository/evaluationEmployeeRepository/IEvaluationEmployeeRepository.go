package evaluationEmployeeRepository

import (
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
)

type IEvaluationEmployeeRepository interface {
	FindByID(id int64) (entity *evaluationModel.EvaluationEmployee, err error)
	Save(tx *gorm.DB, data *evaluationModel.EvaluationEmployee) error
	Delete(evaluationEmployeeData evaluationModel.EvaluationEmployee) error
	FindByEvaluationId(evaluationId int64, typeEmployee string) (employeeId []int64, err error)
}
