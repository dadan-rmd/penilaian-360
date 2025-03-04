package questionRepository

import (
	"penilaian-360/internal/app/model/questionModel"

	"gorm.io/gorm"
)

type IQuestionRepository interface {
	FindByID(id int64) (entity *questionModel.Question, err error)
	FindByEvaluationId(evaluationId int64) (entity []questionModel.Question, err error)
	FindByEvaluationIdAndType(evaluationId int64, typeQuestion, competencyType string) (entity []questionModel.Question, err error)
	Save(tx *gorm.DB, data *[]questionModel.Question) error
	Delete(tx *gorm.DB, id []int64) error
	DeleteEvaluationId(tx *gorm.DB, evaluationId int64) error
	CountRate(tx *gorm.DB, ids []int64) (count int64, err error)
	CountRateByEvaluationIdAndType(tx *gorm.DB, evaluationId int64, typeQuestion, competencyType string) (count int64, err error)
}
