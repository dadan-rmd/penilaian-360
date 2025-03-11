package evaluationService

import (
	"fmt"
	"os"
	"penilaian-360/internal/app/commons/constants"
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
	var (
		accessRoleDepartement = strings.Split(os.Getenv("ACCESS_ROLE_DEPARTEMENT"), ",")
		whitelistUser         = strings.Split(os.Getenv("WHITELIST_USER"), ",")
	)
	if slices.Contains(accessRoleDepartement, employee.Department) || slices.Contains(whitelistUser, employee.Email) {
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
			if res[i].Action != "" {
				res[i].Action = res[i].Action + ",beri-penilaian"
			} else {
				res[i].Action = "beri-penilaian"
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
		if res[i].Action != "" {
			res[i].Action = res[i].Action + ",beri-penilaian"
		} else {
			res[i].Action = "beri-penilaian"
		}
	}
	return
}

func (s evaluationService) EvaluationNeeds(record *loggers.Data, paging datapaging.Datapaging, employee employeeModel.Employee, search string) (res []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {

	res, count, err = s.evaluatedEmployeeRepo.RetrieveNeedsWithPaging(paging, employee.Id, search)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluated RetrieveNeedsWithPaging %v", err))
		return
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
		totalFunctional, totalPersonal, maxFunctional, maxPersonal int64
		avgFunctional, avgPersonal, sumAvg, avg                    float64
		evaluationAnswer                                           []evaluationModel.EvaluationAnswer
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
		if v.Type == string(constants.QuestionTypeRate) && v.CompetencyType == string(constants.TypeOfCompetencyFunctional) {
			totalFunctional += int64(v.FinalPoint)
		} else if v.Type == string(constants.QuestionTypeRate) && v.CompetencyType == string(constants.TypeOfCompetencyPersonal) {
			totalPersonal += int64(v.FinalPoint)
		}
	}
	countRateFunctional, err := s.questionRepo.CountRateByEvaluationIdAndType(tx, req.Data[0].EvaluationId, string(constants.QuestionTypeRate), string(constants.TypeOfCompetencyFunctional))
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator FindByID %v", err))
		return
	}
	countRatePersonal, err := s.questionRepo.CountRateByEvaluationIdAndType(tx, req.Data[0].EvaluationId, string(constants.QuestionTypeRate), string(constants.TypeOfCompetencyPersonal))
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator FindByID %v", err))
		return
	}
	if countRateFunctional > 0 {
		maxFunctional = 5 * countRateFunctional
		avgFunctional = float64(totalFunctional) / float64(maxFunctional) * 100
	}

	if countRatePersonal > 0 {
		maxPersonal = 5 * countRatePersonal
		avgPersonal = float64(totalPersonal) / float64(maxPersonal) * 100
	}
	sumAvg = avgFunctional + avgFunctional
	if sumAvg > 0 {
		avg = sumAvg / 2
	}
	err = s.evaluatorEmployeeRepo.UpdateAvg(tx, req.Data[0].EvaluatorEmployeeId, avgFunctional, avgPersonal, avg)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator UpdateAvg %v", err))
		return
	}

	evaluatorEmployee, err := s.evaluatorEmployeeRepo.FindByID(tx, req.Data[0].EvaluatorEmployeeId)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator FindByID %v", err))
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

func (s evaluationService) EvaluationApprove(record *loggers.Data, evaluatorId int64) (err error) {

	err = s.evaluatorEmployeeRepo.ApproveStatusByEvaluatedEmployeeIdAndEmployeeId(evaluatorId)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, ApproveStatusByEvaluatedEmployeeIdAndEmployeeId %v", err))
		return
	}
	return
}
