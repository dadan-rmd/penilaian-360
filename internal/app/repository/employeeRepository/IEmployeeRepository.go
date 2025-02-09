package employeeRepository

import (
	"penilaian-360/internal/app/model/employeeModel"
)

type IEmployeeRepository interface {
	FindByDepartement(departement string, ids []int64) (entities []employeeModel.Employee, err error)
}
