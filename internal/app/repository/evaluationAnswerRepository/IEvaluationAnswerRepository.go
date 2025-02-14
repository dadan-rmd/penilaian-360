package evaluationAnswerRepository

import (
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
)

type IEvaluationAnswerRepository interface {
	FindByID(id int64) (entity *evaluationModel.EvaluationAnswer, err error)
	Save(tx *gorm.DB, data *[]evaluationModel.EvaluationAnswer) error
	Delete(evaluationAnswerData evaluationModel.EvaluationAnswer) error
}
