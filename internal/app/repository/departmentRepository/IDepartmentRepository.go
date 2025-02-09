package departmentRepository

import "penilaian-360/internal/app/model/departmentModel"

type IDepartmentRepository interface {
	FindAll() (entities []departmentModel.Department, err error)
}
