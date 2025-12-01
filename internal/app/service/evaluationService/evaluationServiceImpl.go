package evaluationService

import (
	"fmt"
	"math"
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
	"strconv"
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

func (s evaluationService) EvaluationDetail(record *loggers.Data, paging datapaging.Datapaging, employeeId int64, params evaluatorEmployeesModel.EvaluatorEmployeeParams) (res []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {

	res, count, err = s.evaluatorEmployeeRepo.RetrieveEvaluatorDetailWithPaging(paging, employeeId, params.Departement, params.Search)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator RetrieveEvaluatorDetailWithPaging %v", err))
		return
	}
	return
}

func (s evaluationService) Score(record *loggers.Data, req evaluationModel.EvaluationAnswerRequests) (err error) {
	var (
		evaluationAnswer []evaluationModel.EvaluationAnswer
		functionalWeight float64
		personalWeight   float64
	)

	if len(req.Data) == 0 {
		loggers.Logf(record, "Warning: empty request data in Score()")
		return fmt.Errorf("empty request data")
	}

	// defaults
	functionalWeight = 0.7
	personalWeight = 0.3

	// read env weights (optional)
	if w := os.Getenv("FUNCTIONAL_WEIGHT"); w != "" {
		if fw, perr := strconv.ParseFloat(w, 64); perr == nil {
			functionalWeight = fw / 100
		} else {
			loggers.Logf(record, fmt.Sprintf("Warning: Invalid FUNCTIONAL_WEIGHT value: %v", perr))
		}
	}
	if w := os.Getenv("PERSONAL_WEIGHT"); w != "" {
		if pw, perr := strconv.ParseFloat(w, 64); perr == nil {
			personalWeight = pw / 100
		} else {
			loggers.Logf(record, fmt.Sprintf("Warning: Invalid PERSONAL_WEIGHT value: %v", perr))
		}
	}

	// normalize weights
	eps := 1e-9
	sumWeights := functionalWeight + personalWeight
	if math.Abs(sumWeights-1.0) > eps {
		if sumWeights > 0 {
			functionalWeight = functionalWeight / sumWeights
			personalWeight = personalWeight / sumWeights
			loggers.Logf(record, fmt.Sprintf("Info: weights normalized. Functional: %.2f%%, Personal: %.2f%%", functionalWeight*100, personalWeight*100))
		} else {
			functionalWeight = 0.7
			personalWeight = 0.3
			loggers.Logf(record, "Warning: invalid weights from env, fallback to defaults 70/30")
		}
	}

	// collect evaluationAnswer for saving (inner struct)
	for _, v := range req.Data {
		evaluationAnswer = append(evaluationAnswer, v.EvaluationAnswer)
	}

	// begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		loggers.Logf(record, fmt.Sprintf("Err: begin tx %v", tx.Error))
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			loggers.Logf(record, fmt.Sprintf("Panic recovered in Score: %v", r))
			err = fmt.Errorf("panic: %v", r)
			return
		}
		if err != nil {
			if rb := tx.Rollback(); rb != nil {
				loggers.Logf(record, fmt.Sprintf("Err, Rollback %v", rb))
			}
			return
		}
		if cm := tx.Commit().Error; cm != nil {
			loggers.Logf(record, fmt.Sprintf("Err, Commit %v", cm))
			err = cm
		}
	}()

	// save answers
	if err = s.evaluationAnswerRepo.Save(tx, &evaluationAnswer); err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, Save %v", err))
		return
	}

	// ---- NEW: accumulate totals & counts per evaluator using maps (use req.Data fields) ----
	totalsFunctional := make(map[int64]int64) // evaluatorID -> total points functional
	totalsPersonal := make(map[int64]int64)   // evaluatorID -> total points personal
	countsFunctional := make(map[int64]int64) // evaluatorID -> count functional answers
	countsPersonal := make(map[int64]int64)   // evaluatorID -> count personal answers
	evaluatorIDs := make(map[int64]struct{})  // set of evaluator IDs encountered

	for _, item := range req.Data {
		eid := item.EvaluationAnswer.EvaluatorEmployeeId
		evaluatorIDs[eid] = struct{}{}

		// Use the wrapper fields from item (not inner EvaluationAnswer) — this fixes the undefined field error
		if item.Type == string(constants.QuestionTypeRate) && item.CompetencyType == string(constants.TypeOfCompetencyFunctional) {
			totalsFunctional[eid] += int64(item.FinalPoint)
			countsFunctional[eid]++
		} else if item.Type == string(constants.QuestionTypeRate) && item.CompetencyType == string(constants.TypeOfCompetencyPersonal) {
			totalsPersonal[eid] += int64(item.FinalPoint)
			countsPersonal[eid]++
		}
	}

	// process every evaluator we saw
	for eid := range evaluatorIDs {
		totalFunc := totalsFunctional[eid]
		totalPers := totalsPersonal[eid]
		cntFunc := countsFunctional[eid]
		cntPers := countsPersonal[eid]

		// global total (based only on submitted answers)
		totalAll := totalFunc + totalPers
		cntAll := cntFunc + cntPers

		var percentAll float64
		if cntAll > 0 {
			maxAll := float64(5 * cntAll)
			percentAll = float64(totalAll) / maxAll * 100.0
		} else {
			percentAll = 0
		}

		avgFunctional := percentAll * functionalWeight
		avgPersonal := percentAll * personalWeight
		avg := percentAll // overall percent 0..100

		loggers.Logf(record, fmt.Sprintf(
			"Score computed (evaluator=%d): totalFunctional=%d totalPersonal=%d countFunc=%d countPers=%d percentAll=%.6f avgFunctional=%.6f avgPersonal=%.6f avgTotal=%.6f",
			eid, totalFunc, totalPers, cntFunc, cntPers, percentAll, avgFunctional, avgPersonal, avg,
		))

		// update evaluator avg
		if err = s.evaluatorEmployeeRepo.UpdateAvg(tx, eid, avgFunctional, avgPersonal, avg); err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluator UpdateAvg %v", err))
			return err
		}

		// recompute totalAvg for evaluated employee
		evaluatorEmployee, err2 := s.evaluatorEmployeeRepo.FindByID(tx, eid)
		if err2 != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluator FindByID %v", err2))
			return err2
		}

		totalAvg, err2 := s.evaluatorEmployeeRepo.TotalAvg(tx, evaluatorEmployee.EvaluatedEmployeeId)
		if err2 != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluator TotalAvg %v", err2))
			return err2
		}

		loggers.Logf(record, fmt.Sprintf("TotalAvg for evaluated_employee_id=%d is %.6f", evaluatorEmployee.EvaluatedEmployeeId, totalAvg))

		if err2 = s.evaluatedEmployeeRepo.UpdateAvg(tx, evaluatorEmployee.EvaluatedEmployeeId, totalAvg); err2 != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluated UpdateAvg %v", err2))
			return err2
		}
	}

	return nil
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
