package evaluationRepository

import (
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
)

type IEvaluationRepository interface {
	FindByID(id int64) (entity *evaluationModel.Evaluation, err error)
	Save(tx *gorm.DB, data *evaluationModel.Evaluation) error
	Delete(evaluationData evaluationModel.Evaluation) error
}
