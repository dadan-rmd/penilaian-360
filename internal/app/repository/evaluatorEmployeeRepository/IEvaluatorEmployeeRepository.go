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
	RetrieveListWithPaging(paging datapaging.Datapaging, employeeId int64, email, notDepartement, departement, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error)
	RetrieveEvaluatorDetailWithPaging(paging datapaging.Datapaging, evaluatedId int64, departement, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error)
	FindByID(tx *gorm.DB, id int64) (entity *evaluatorEmployeesModel.EvaluatorEmployee, err error)
	TotalAvg(tx *gorm.DB, evaluatedEmployeeId int64) (totalAvg float64, err error)
	UpdateAvg(tx *gorm.DB, id int64, avg float64) (err error)
}
