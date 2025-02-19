package evaluationService

import (
	"fmt"
	"os"
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/model/evaluationModel"
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
		res, count, err = s.evaluatorEmployeeRepo.RetrieveListWithPaging(paging, employee.Id, employee.Email, employee.Department, params.Departement, params.Search)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluator RetrieveListWithPaging %v", err))
			return
		}
		for i := 0; i < len(res); i++ {
			if res[i].Status != "" {
				res[i].Status = res[i].Status + ",beri-penilaian"
			} else {
				res[i].Status = "beri-penilaian"
			}
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
	for i := 0; i < len(res); i++ {
		if res[i].Status != "" {
			res[i].Status = res[i].Status + ",beri-penilaian"
		} else {
			res[i].Status = "beri-penilaian"
		}
	}
	return
}

func (s evaluationService) EvaluationDetail(record *loggers.Data, paging datapaging.Datapaging, evaluatedId int64, params evaluatorEmployeesModel.EvaluatorEmployeeParams) (res []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {

	res, count, err = s.evaluatorEmployeeRepo.RetrieveEvaluatorDetailWithPaging(paging, evaluatedId, params.Departement, params.Search)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator RetrieveEvaluatorDetailWithPaging %v", err))
		return
	}
	return
}

func (s evaluationService) Score(record *loggers.Data, req evaluationModel.EvaluationAnswerRequests) (err error) {
	var (
		idQuestions      []int64
		totalValue       int64
		maxValue         int64
		avg              float64
		evaluationAnswer []evaluationModel.EvaluationAnswer
	)
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				loggers.Logf(record, fmt.Sprintf("Err, Rollback %v", errRollback))
			}
			return
		}
		if err = tx.Commit().Error; err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, Commit %v", err))
		}
	}()
	for _, v := range req.Data {
		evaluationAnswer = append(evaluationAnswer, v.EvaluationAnswer)
	}
	err = s.evaluationAnswerRepo.Save(tx, &evaluationAnswer)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, Save %v", err))
		return
	}
	for _, v := range req.Data {
		idQuestions = append(idQuestions, v.Id)
		totalValue += int64(v.FinalPoint)
	}
	countRate, err := s.questionRepo.CountRate(tx, idQuestions)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator FindByID %v", err))
		return
	}
	maxValue = 5 * countRate
	evaluatorEmployee, err := s.evaluatorEmployeeRepo.FindByID(tx, req.Data[0].EvaluatorEmployeeId)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator FindByID %v", err))
		return
	}

	avg = float64(totalValue) / float64(maxValue) * 100
	err = s.evaluatorEmployeeRepo.UpdateAvg(tx, req.Data[0].EvaluatorEmployeeId, avg)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator UpdateAvg %v", err))
		return
	}
	totalAvg, err := s.evaluatorEmployeeRepo.TotalAvg(tx, evaluatorEmployee.EvaluatedEmployeeId)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator TotalAvg %v", err))
		return
	}
	err = s.evaluatedEmployeeRepo.UpdateAvg(tx, evaluatorEmployee.EvaluatedEmployeeId, totalAvg)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluated UpdateAvg %v", err))
		return
	}

	return
}

func (s evaluationService) ScoreDetail(record *loggers.Data, evaluationId, evaluatorEmployeeId int64) (res *[]evaluationModel.EvaluationAnswerResponse, err error) {
	res, err = s.evaluationAnswerRepo.FindByEvaluationAndevaluatorID(evaluationId, evaluatorEmployeeId)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationAndevaluatorID %v", err))
		return
	}
	return
}
