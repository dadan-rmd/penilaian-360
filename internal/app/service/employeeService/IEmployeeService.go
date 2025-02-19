package employeeService

import (
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/employeeModel"
)

type IEmployeeService interface {
	GetEmployeeAll(record *loggers.Data, params employeeModel.EmployeeParamas) (res []employeeModel.EmployeeResponse, err error)
	CreateToken(record *loggers.Data, email, accessToken string) (token string, err error)
}
