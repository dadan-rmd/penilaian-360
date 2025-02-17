package evaluationService

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"
)

type IEvaluationService interface {
	EvaluationList(record *loggers.Data, paging datapaging.Datapaging, employee employeeModel.Employee, params evaluatorEmployeesModel.EvaluatorEmployeeParams) (res []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error)
	EvaluationWithDepartementList(record *loggers.Data, paging datapaging.Datapaging, employee employeeModel.Employee) (res []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error)
}
