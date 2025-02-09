package employeeService

import (
	"fmt"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/repository/employeeRepository"
	"penilaian-360/internal/app/repository/evaluationEmployeeRepository"
	"strings"
)

type employeeService struct {
	employeeRepo           employeeRepository.IEmployeeRepository
	evaluationEmployeeRepo evaluationEmployeeRepository.IEvaluationEmployeeRepository
}

func NewEmployeeService(
	employeeRepo employeeRepository.IEmployeeRepository,
	evaluationEmployeeRepo evaluationEmployeeRepository.IEvaluationEmployeeRepository,
) IEmployeeService {
	return &employeeService{employeeRepo, evaluationEmployeeRepo}
}

func (s employeeService) GetEmployeeAll(record *loggers.Data, params employeeModel.EmployeeParamas) (res []employeeModel.EmployeeResponse, err error) {
	if strings.ToUpper(params.Departement) == "ALL" {
		params.Departement = ""
	}
	employeeId, err := s.evaluationEmployeeRepo.FindByEvaluationId(params.EvaluationId, params.Type)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationId %v", err))
		return
	}
	entities, err := s.employeeRepo.FindByDepartement(params.Departement, employeeId)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByDepartement %v", err))
		return
	}
	for _, v := range entities {
		res = append(res, v.ToEmployeeResponse())
	}
	return
}
