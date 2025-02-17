package evaluationService

import (
	"fmt"
	"os"
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"
	"penilaian-360/internal/app/repository/employeeRepository"
	"penilaian-360/internal/app/repository/evaluatedEmployeeRepository"
	"penilaian-360/internal/app/repository/evaluationAnswerRepository"
	"penilaian-360/internal/app/repository/evaluationRepository"
	"penilaian-360/internal/app/repository/evaluatorEmployeeRepository"
	"penilaian-360/internal/app/repository/questionRepository"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type evaluationService struct {
	db                    *gorm.DB
	evaluationRepo        evaluationRepository.IEvaluationRepository
	questionRepo          questionRepository.IQuestionRepository
	employeeRepo          employeeRepository.IEmployeeRepository
	evaluationAnswerRepo  evaluationAnswerRepository.IEvaluationAnswerRepository
	evaluatorEmployeeRepo evaluatorEmployeeRepository.IEvaluatorEmployeeRepository
	evaluatedEmployeeRepo evaluatedEmployeeRepository.IEvaluatedEmployeeRepository
}

func NewEvaluationService(
	db *gorm.DB,
	evaluationRepo evaluationRepository.IEvaluationRepository,
	questionRepo questionRepository.IQuestionRepository,
	employeeRepo employeeRepository.IEmployeeRepository,
	evaluationAnswerRepo evaluationAnswerRepository.IEvaluationAnswerRepository,
	evaluatorEmployeeRepo evaluatorEmployeeRepository.IEvaluatorEmployeeRepository,
	evaluatedEmployeeRepo evaluatedEmployeeRepository.IEvaluatedEmployeeRepository,
) IEvaluationService {
	return &evaluationService{
		db,
		evaluationRepo,
		questionRepo,
		employeeRepo,
		evaluationAnswerRepo,
		evaluatorEmployeeRepo,
		evaluatedEmployeeRepo,
	}
}

func (s evaluationService) EvaluationList(record *loggers.Data, paging datapaging.Datapaging, employee employeeModel.Employee, params evaluatorEmployeesModel.EvaluatorEmployeeParams) (res []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {
	accessRoleDepartement := strings.Split(os.Getenv("ACCESS_ROLE_DEPARTEMENT"), ",")
	if slices.Contains(accessRoleDepartement, employee.Department) {
		res, count, err = s.evaluatedEmployeeRepo.RetrieveListWithPaging(paging, params.Departement, params.Search)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluated RetrieveListWithPaging %v", err))
			return
		}
	} else {
		res, count, err = s.evaluatorEmployeeRepo.RetrieveListWithPaging(paging, employee.Id, "", employee.Department, params.Departement, params.Search)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluator RetrieveListWithPaging %v", err))
			return
		}
	}
	return
}

func (s evaluationService) EvaluationWithDepartementList(record *loggers.Data, paging datapaging.Datapaging, employee employeeModel.Employee) (res []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {

	res, count, err = s.evaluatorEmployeeRepo.RetrieveListWithPaging(paging, employee.Id, employee.Email, "", employee.Department, "")
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator RetrieveListWithPaging %v", err))
		return
	}
	return
}
