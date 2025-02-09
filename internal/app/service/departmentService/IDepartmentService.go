package departmentService

import (
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/departmentModel"
)

type IDepartmentService interface {
	GetDepartmentAll(record *loggers.Data) (res []departmentModel.DepartmentResponse, err error)
}
