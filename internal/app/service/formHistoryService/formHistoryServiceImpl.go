package formHistoryService

import (
	"errors"
	"fmt"
	"penilaian-360/internal/app/commons/constants"
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/mail"
	"penilaian-360/internal/app/model/evaluationModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"
	"penilaian-360/internal/app/model/questionModel"
	"penilaian-360/internal/app/repository/employeeRepository"
	"penilaian-360/internal/app/repository/evaluatedEmployeeRepository"
	"penilaian-360/internal/app/repository/evaluationAnswerRepository"
	"penilaian-360/internal/app/repository/evaluationRepository"
	"penilaian-360/internal/app/repository/evaluatorEmployeeRepository"
	"penilaian-360/internal/app/repository/questionRepository"
	"strings"
	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type formHistoryService struct {
	db                    *gorm.DB
	evaluationRepo        evaluationRepository.IEvaluationRepository
	questionRepo          questionRepository.IQuestionRepository
	employeeRepo          employeeRepository.IEmployeeRepository
	evaluationAnswerRepo  evaluationAnswerRepository.IEvaluationAnswerRepository
	evaluatorEmployeeRepo evaluatorEmployeeRepository.IEvaluatorEmployeeRepository
	evaluatedEmployeeRepo evaluatedEmployeeRepository.IEvaluatedEmployeeRepository
}

func NewFormHistoryService(
	db *gorm.DB,
	evaluationRepo evaluationRepository.IEvaluationRepository,
	questionRepo questionRepository.IQuestionRepository,
	employeeRepo employeeRepository.IEmployeeRepository,
	evaluationAnswerRepo evaluationAnswerRepository.IEvaluationAnswerRepository,
	evaluatorEmployeeRepo evaluatorEmployeeRepository.IEvaluatorEmployeeRepository,
	evaluatedEmployeeRepo evaluatedEmployeeRepository.IEvaluatedEmployeeRepository,
) IFormHistoryService {
	return &formHistoryService{
		db,
		evaluationRepo,
		questionRepo,
		employeeRepo,
		evaluationAnswerRepo,
		evaluatorEmployeeRepo,
		evaluatedEmployeeRepo,
	}
}

func (s formHistoryService) SaveFormHistory(record *loggers.Data, request evaluationModel.FormHistoryRequest) (res evaluationModel.FormHistoryResponse, err error) {
	evaluation := evaluationModel.Evaluation{
		Id:            request.Id,
		DepartementId: request.DepartementId,
		Title:         request.Title,
		Status:        constants.EvaluationStatus(request.Status),
		DeadlineAt:    request.DeadlineAt,
		CreatedAt:     time.Now(),
	}
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

	err = s.evaluationRepo.Save(tx, &evaluation)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluation Save %v", err))
		return
	}
	var questions []questionModel.Question
	for _, v := range request.Functional {
		questions = append(questions, questionModel.Question{
			Id:             v.Id,
			EvaluationId:   evaluation.Id,
			Title:          v.Title,
			Question:       v.Question,
			Type:           constants.QuestionTypeRate,
			CompetencyType: constants.TypeOfCompetencyFunctional,
		})
	}
	for _, v := range request.Personal {
		questions = append(questions, questionModel.Question{
			Id:             v.Id,
			EvaluationId:   evaluation.Id,
			Title:          v.Title,
			Question:       v.Question,
			Type:           constants.QuestionTypeRate,
			CompetencyType: constants.TypeOfCompetencyPersonal,
		})
	}
	for _, v := range request.Essay {
		questions = append(questions, questionModel.Question{
			Id:           v.Id,
			EvaluationId: evaluation.Id,
			Title:        v.Title,
			Question:     v.Question,
			Type:         constants.QuestionTypeEssay,
		})
	}
	err = s.questionRepo.Save(tx, &questions)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, question Save %v", err))
		return
	}
	if len(request.IdToDeleteQuestion) > 0 {
		err = s.questionRepo.Delete(tx, request.IdToDeleteQuestion)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, question Delete %v", err))
			return
		}
	}
	res = evaluationModel.FormHistoryResponse{
		Id: evaluation.Id,
	}
	return
}

func (s formHistoryService) FormHistoryList(record *loggers.Data, paging datapaging.Datapaging) (res []evaluationModel.FormHistoryList, count int64, err error) {
	res, count, err = s.evaluationRepo.GetWithPaging(paging)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, GetWithPaging %v", err))
		return
	}
	return
}

func (s formHistoryService) FormHistoryView(record *loggers.Data, id int64) (res questionModel.QuestionWithEvaluation, err error) {
	evaluation, err := s.evaluationRepo.FindByID(id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByFormHistoryId %v", err))
		return
	}
	functional, err := s.questionRepo.FindByEvaluationIdAndType(id, string(constants.QuestionTypeRate), string(constants.TypeOfCompetencyFunctional))
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationIdAndType functional %v", err))
		return
	}
	personal, err := s.questionRepo.FindByEvaluationIdAndType(id, string(constants.QuestionTypeRate), string(constants.TypeOfCompetencyPersonal))
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationIdAndType functional %v", err))
		return
	}
	essay, err := s.questionRepo.FindByEvaluationIdAndType(id, string(constants.QuestionTypeEssay), "")
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationIdAndType functional %v", err))
		return
	}
	res = questionModel.QuestionWithEvaluation{
		Evaluation: *evaluation,
		Functional: questionModel.ToAssemblerQuestion(functional),
		Personal:   questionModel.ToAssemblerQuestion(personal),
		Essay:      questionModel.ToAssemblerQuestion(essay),
	}
	return
}

func (s formHistoryService) FormHistoryDelete(record *loggers.Data, id int64) (err error) {
	entity, err := s.evaluationRepo.FindByID(id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByFormHistoryId %v", err))
		return
	}
	if entity.Status == constants.EvaluationStatusPublish {
		err = errors.New("Published data cannot be deleted.")
		return
	}
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
	err = s.questionRepo.DeleteEvaluationId(tx, id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, DeleteFormHistoryId %v", err))
		return
	}
	err = s.evaluationRepo.Delete(tx, id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, Delete %v", err))
		return
	}
	return
}

func (s formHistoryService) FormHistoryAssignment(record *loggers.Data, request evaluationModel.AssignmentRequest) (err error) {
	var (
		employeeID          []int64
		evaluatedEmployeeID []int64
		evaluationAnswer    []evaluationModel.EvaluationAnswer
		evaluators          []evaluatorEmployeesModel.EvaluatorEmployee
		requiresAssessment  bool
	)
	entity, err := s.evaluationRepo.FindByID(request.Id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByFormHistoryId %v", err))
		return
	}
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

		employees, err := s.employeeRepo.FindByIds(request.EvaluatorId)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, employees FindByIds %v", err))
			return
		}

		evaluatedEmployees, err := s.employeeRepo.FindNameAndEmployedIDByIds(request.EvaluatedId)
		if err != nil {
			loggers.Logf(record, fmt.Sprintf("Err, employees FindByIds %v", err))
			return
		}
		parsedTime, _ := time.Parse(constants.DDMMYYYY, entity.DeadlineAt)
		deadline := parsedTime.Format("02 January 2006")
		cc := strings.Split(request.Cc, ",")
		go func() {
			for _, evaluatedEmployee := range evaluatedEmployees {
				for _, v := range employees {
					err = mail.SendEvaluation([]string{v.Email}, cc, cast.ToString(evaluatedEmployee.Name), v.Name, deadline)
					if err == nil {
						employeeID = append(employeeID, v.Id)
					}
				}
				evaluatedEmployeeID = append(evaluatedEmployeeID, evaluatedEmployee.EvaluatedId)
			}
			err = s.evaluatorEmployeeRepo.UpdateEmailSentByEvaluatedEmployeeIdAndEmployeeId(employeeID, evaluatedEmployeeID)
			if err != nil {
				fmt.Println("error UpdateEmailSentById : %v", err)
			}
		}()
	}()
	evaluateds := request.ToEvaluatedEmployee()
	for i, v := range evaluateds {
		entity, errEntity := s.evaluatedEmployeeRepo.FindByEvaluationIdAndEmployeeId(v.EvaluationId, v.EmployeeId)
		if errEntity != nil {
			loggers.Logf(record, fmt.Sprintf("Err, evaluated FindByEvaluationIdAndEmployeeId Save %v", err))
		}
		evaluateds[i].Id = entity.Id
		if entity.Id == 0 {
			fmt.Println("masuk if")
			err = s.evaluatedEmployeeRepo.Save(tx, &evaluateds)
			if err != nil {
				loggers.Logf(record, fmt.Sprintf("Err, evaluated employee Save %v", err))
				return
			}
		}
	}

	countRate, err := s.questionRepo.CountRateByEvaluationIdAndType(tx, request.Id, string(constants.QuestionTypeRate), "")
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, CountRateByEvaluationIdAndType Functional %v", err))
		return
	}
	if countRate > 0 {
		requiresAssessment = true
	}
	for _, v := range evaluateds {
		evaluators = append(evaluators, request.ToEvaluatorEmployee(v.Id, requiresAssessment)...)
	}
	err = s.evaluatorEmployeeRepo.Save(tx, &evaluators)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluator employee Save %v", err))
		return
	}

	questions, err := s.questionRepo.FindByEvaluationId(request.Id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, question FindByFormHistoryId %v", err))
		return
	}
	for _, evaluator := range evaluators {
		for _, question := range questions {
			evaluationAnswer = append(evaluationAnswer, evaluationModel.EvaluationAnswer{
				EvaluationId:        request.Id,
				EvaluatorEmployeeId: evaluator.Id,
				QuestionId:          question.Id,
			})
		}
	}
	err = s.evaluationAnswerRepo.Save(tx, &evaluationAnswer)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluation Answer Save %v", err))
		return
	}

	return
}

func (s formHistoryService) FormHistoryDetail(record *loggers.Data, paging datapaging.Datapaging, params evaluatorEmployeesModel.FormHistoryDetailParams) (res evaluatorEmployeesModel.FormHistoryDetailResponse, count int64, err error) {
	departmentName, err := s.evaluationRepo.FindDepartmentNameByID(params.ID)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindEvaluated %v", err))
		return
	}
	evaluator, count, err := s.evaluatorEmployeeRepo.FindByEvaluatorId(paging, params.ID, params.EmployeeIdID)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindEvaluated %v", err))
		return
	}
	res = evaluatorEmployeesModel.FormHistoryDetailResponse{
		Department: departmentName,
		Data:       evaluator,
	}
	return
}

func (s formHistoryService) CopyFormHistory(record *loggers.Data, id int64) (res evaluationModel.FormHistoryResponse, err error) {

	evaluation, err := s.evaluationRepo.FindByID(id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, evaluation FindByID %v", err))
		return
	}
	functional, err := s.questionRepo.FindByEvaluationIdAndType(id, string(constants.QuestionTypeRate), string(constants.TypeOfCompetencyFunctional))
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationIdAndType functional %v", err))
		return
	}
	personal, err := s.questionRepo.FindByEvaluationIdAndType(id, string(constants.QuestionTypeRate), string(constants.TypeOfCompetencyPersonal))
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationIdAndType functional %v", err))
		return
	}
	essay, err := s.questionRepo.FindByEvaluationIdAndType(id, string(constants.QuestionTypeEssay), "")
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationIdAndType functional %v", err))
		return
	}
	request := evaluationModel.FormHistoryRequest{
		DataFormHistory: evaluationModel.DataFormHistory{
			DepartementId: evaluation.DepartementId,
			Title:         evaluation.Title,
			Status:        string(constants.EvaluationStatusDraft),
			DeadlineAt:    evaluation.DeadlineAt,
		},
		Functional: questionModel.ToAssemblerQuestionV2(functional),
		Personal:   questionModel.ToAssemblerQuestionV2(personal),
		Essay:      questionModel.ToAssemblerQuestionV2(essay),
	}
	res, err = s.SaveFormHistory(record, request)
	return
}
