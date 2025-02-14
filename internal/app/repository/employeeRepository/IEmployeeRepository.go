package employeeRepository

import (
	"penilaian-360/internal/app/model/employeeModel"
)

type IEmployeeRepository interface {
	FindByDepartement(departement string, ids []int64) (entities []employeeModel.Employee, err error)
	FindByEmailAndAccessToken(email, accessToken string) (entity employeeModel.Employee, err error)
	FindByIds(ids []int64) (entity []employeeModel.Employee, err error)
}
