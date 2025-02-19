package evaluatedEmployeeRepository

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/model/evaluatedEmployeesModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"

	"gorm.io/gorm"
)

type IEvaluatedEmployeeRepository interface {
	Save(tx *gorm.DB, data *[]evaluatedEmployeesModel.EvaluatedEmployee) error
	FindEmployeeIdByEvaluationId(evaluationId int64) (employeeId []int64, err error)
	RetrieveListWithPaging(paging datapaging.Datapaging, departement, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error)
	UpdateAvg(tx *gorm.DB, id int64, totalAvg float64) (err error)
}
