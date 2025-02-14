package evaluationAnswerRepository

import (
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type evaluationAnswerRepository struct {
	db *gorm.DB
}

func NewEvaluationAnswerRepository(db *gorm.DB) IEvaluationAnswerRepository {
	return &evaluationAnswerRepository{db}
}

func (d evaluationAnswerRepository) FindByID(id int64) (entity *evaluationModel.EvaluationAnswer, err error) {
	db := d.db.Preload(clause.Associations)
	err = db.First(&entity, "id=? ", id).Error
	return
}

func (d evaluationAnswerRepository) Save(tx *gorm.DB, data *[]evaluationModel.EvaluationAnswer) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d evaluationAnswerRepository) Delete(evaluationAnswerData evaluationModel.EvaluationAnswer) error {
	db := d.db.Delete(&evaluationAnswerData)
	return db.Error
}
