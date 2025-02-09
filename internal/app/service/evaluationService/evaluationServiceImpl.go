package evaluationService

import (
	"penilaian-360/internal/app/commons/loggers"

	"gorm.io/gorm"
)

type evaluationService struct {
	db *gorm.DB
}

func NewEvaluationService(
	db *gorm.DB,
) IEvaluationService {
	return &evaluationService{db}
}

func (a evaluationService) SaveEvaluation(record *loggers.Data) (err error) {

	return
}
