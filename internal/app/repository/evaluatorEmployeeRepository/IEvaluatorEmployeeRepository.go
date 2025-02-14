package evaluatorEmployeeRepository

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"

	"gorm.io/gorm"
)

type IEvaluatorEmployeeRepository interface {
	FindEmployeeIdByEvaluationId(evaluationId int64) (employeeId []int64, err error)
	UpdateEmailSentByEvaluatedEmployeeIdAndEmployeeId(ids, evaluatedEmployeeId []int64) error
	FindByEvaluatorId(paging datapaging.Datapaging, evaluationId, evaluatedEmployeeId int64) (entities []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error)
	Save(tx *gorm.DB, data *[]evaluatorEmployeesModel.EvaluatorEmployee) error
}
