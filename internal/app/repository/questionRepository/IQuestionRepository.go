package questionRepository

import (
	"penilaian-360/internal/app/model/questionModel"

	"gorm.io/gorm"
)

type IQuestionRepository interface {
	FindByID(id int64) (entity *questionModel.Question, err error)
	Save(tx *gorm.DB, data *questionModel.Question) error
	Delete(questionData questionModel.Question) error
}
