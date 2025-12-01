package employeeService

import (
	"errors"
	"fmt"
	"penilaian-360/internal/app/commons/constants"
	"penilaian-360/internal/app/commons/jwtHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"
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
	if params.Type == string(constants.EmployeeTypeEvaluator) {
		ids, err = s.evaluatorEmployeeRepo.FindEmployeeIdByEvaluationId(params.FormId, params.EmployeeId)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluator FindEmployeeIdByEvaluationId %v", err))
			return
		}
	} else if params.Type == string(constants.EmployeeTypeEvaluated) {
		// ids, err = s.evaluatedEmployeeRepo.FindEmployeeIdByEvaluationId(params.FormId)
		// if err != nil {
		// 	loggers.Logf(record, fmt.Sprintf("Err, evaluated FindEmployeeIdByEvaluationId %v", err))
		// 	return
		// }
		ids = []int64{}
	} else {
		err = errors.New("Type not found")
		return
	}
	entities, err := s.employeeRepo.FindByDepartement(params.Departement, ids, params.HasAssigned)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByDepartement %v", err))
		return
	}
	for _, v := range entities {
		res = append(res, v.ToEmployeeResponse())
	}
	return
}

func (s employeeService) CreateToken(record *loggers.Data, email, accessToken string) (token string, err error) {
	token, err = jwtHelper.EncodeJWT(email, accessToken)
	return
}

func (s employeeService) GetEmployeeEmails(record *loggers.Data, params evaluatorEmployeesModel.EvaluatorEmployeeParams) (res []string, err error) {

	if strings.ToUpper(params.Departement) == "ALL" {
		params.Departement = ""
	}

	keyword := strings.TrimSpace(params.Search)
	if len(keyword) < 3 {
		err = fmt.Errorf("search keyword must be at least 3 characters")
		return
	}

	emails, err := s.employeeRepo.FindEmailsByKeyword(keyword)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindEmailsByKeyword %v", err))
		return
	}

	return emails, nil
}
