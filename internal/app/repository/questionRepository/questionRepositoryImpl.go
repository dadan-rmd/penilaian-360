package questionRepository

import (
	"penilaian-360/internal/app/model/questionModel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) IQuestionRepository {
	return &questionRepository{db}
}

func (d questionRepository) FindByID(id int64) (entity *questionModel.Question, err error) {
	db := d.db.Preload(clause.Associations)
	err = db.First(&entity, "id=? ", id).Error
	return
}

func (d questionRepository) Save(tx *gorm.DB, data *questionModel.Question) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d questionRepository) Delete(questionData questionModel.Question) error {
	db := d.db.Delete(&questionData)
	return db.Error
}
