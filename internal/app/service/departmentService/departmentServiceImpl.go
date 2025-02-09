package departmentService

import (
	"fmt"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/departmentModel"
	"penilaian-360/internal/app/repository/departmentRepository"
)

type departmentService struct {
	departmentRepo departmentRepository.IDepartmentRepository
}

func NewDepartmentService(
	departmentRepo departmentRepository.IDepartmentRepository,
) IDepartmentService {
	return &departmentService{departmentRepo}
}

func (s departmentService) GetDepartmentAll(record *loggers.Data) (res []departmentModel.DepartmentResponse, err error) {
	entities, err := s.departmentRepo.FindAll()
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindAll %v", err))
		return
	}
	for _, v := range entities {
		res = append(res, v.ToDepartmentResponse())
	}
	return
}
