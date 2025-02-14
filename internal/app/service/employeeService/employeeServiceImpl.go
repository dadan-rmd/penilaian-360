package employeeService

import (
	"errors"
	"fmt"
	"penilaian-360/internal/app/commons/constants"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/repository/employeeRepository"
	"penilaian-360/internal/app/repository/evaluatedEmployeeRepository"
	"penilaian-360/internal/app/repository/evaluatorEmployeeRepository"
	"strings"
)

type employeeService struct {
	employeeRepo          employeeRepository.IEmployeeRepository
	evaluatorEmployeeRepo evaluatorEmployeeRepository.IEvaluatorEmployeeRepository
	evaluatedEmployeeRepo evaluatedEmployeeRepository.IEvaluatedEmployeeRepository
}

func NewEmployeeService(
	employeeRepo employeeRepository.IEmployeeRepository,
	evaluatorEmployeeRepo evaluatorEmployeeRepository.IEvaluatorEmployeeRepository,
	evaluatedEmployeeRepo evaluatedEmployeeRepository.IEvaluatedEmployeeRepository,
) IEmployeeService {
	return &employeeService{employeeRepo, evaluatorEmployeeRepo, evaluatedEmployeeRepo}
}

func (s employeeService) GetEmployeeAll(record *loggers.Data, params employeeModel.EmployeeParamas) (res []employeeModel.EmployeeResponse, err error) {
	var (
		ids []int64
	)
	if strings.ToUpper(params.Departement) == "ALL" {
		params.Departement = ""
	}
	if strings.ToUpper(params.Type) == string(constants.EmployeeTypeEvaluator) {
		ids, err = s.evaluatorEmployeeRepo.FindEmployeeIdByEvaluationId(params.EvaluationId)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluator FindEmployeeIdByEvaluationId %v", err))
			return
		}
	} else if strings.ToUpper(params.Type) == string(constants.EmployeeTypeEvaluated) {
		ids, err = s.evaluatedEmployeeRepo.FindEmployeeIdByEvaluationId(params.EvaluationId)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluated FindEmployeeIdByEvaluationId %v", err))
			return
		}
	} else {
		err = errors.New("Type not found")
		return
	}
	entities, err := s.employeeRepo.FindByDepartement(params.Departement, ids)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByDepartement %v", err))
		return
	}
	for _, v := range entities {
		res = append(res, v.ToEmployeeResponse())
	}
	return
}
