package evaluationService

import (
	"errors"
	"fmt"
	"penilaian-360/internal/app/commons/constants"
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/evaluationModel"
	"penilaian-360/internal/app/model/questionModel"
	"penilaian-360/internal/app/repository/evaluationRepository"
	"penilaian-360/internal/app/repository/questionRepository"
	"time"

	"gorm.io/gorm"
)

type evaluationService struct {
	db             *gorm.DB
	evaluationRepo evaluationRepository.IEvaluationRepository
	questionRepo   questionRepository.IQuestionRepository
}

func NewEvaluationService(
	db *gorm.DB,
	evaluationRepo evaluationRepository.IEvaluationRepository,
	questionRepo questionRepository.IQuestionRepository,
) IEvaluationService {
	return &evaluationService{db, evaluationRepo, questionRepo}
}

func (s evaluationService) SaveEvaluation(record *loggers.Data, request evaluationModel.EvaluationRequest) (res evaluationModel.EvaluationResponse, err error) {
	evaluation := evaluationModel.Evaluation{
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
	for _, v := range request.Question {
		questions = append(questions, questionModel.Question{
			EvaluationId: evaluation.Id,
			Question:     v.Question,
			Type:         constants.QuestionType(v.Type),
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
	res = evaluationModel.EvaluationResponse{
		Id: evaluation.Id,
	}
	return
}

func (s evaluationService) EvaluationList(record *loggers.Data, paging datapaging.Datapaging) (res []evaluationModel.EvaluationList, count int64, err error) {
	res, count, err = s.evaluationRepo.GetWithPaging(paging)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, GetWithPaging %v", err))
		return
	}
	return
}

func (s evaluationService) EvaluationView(record *loggers.Data, id int64) (res []questionModel.Question, err error) {
	res, err = s.questionRepo.FindByEvaluationId(id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationId %v", err))
		return
	}
	return
}

func (s evaluationService) EvaluationDelete(record *loggers.Data, id int64) (err error) {
	entity, err := s.evaluationRepo.FindByID(id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, FindByEvaluationId %v", err))
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
		loggers.Logf(record, fmt.Sprintf("Err, DeleteEvaluationId %v", err))
		return
	}
	err = s.evaluationRepo.Delete(tx, id)
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, Delete %v", err))
		return
	}
	return
}
