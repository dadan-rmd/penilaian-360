package evaluationRepository

import (
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type evaluationRepository struct {
	db *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) IEvaluationRepository {
	return &evaluationRepository{db}
}

func (d evaluationRepository) FindByID(id int64) (entity *evaluationModel.Evaluation, err error) {
	db := d.db.Preload(clause.Associations)
	err = db.First(&entity, "id=? ", id).Error
	return
}

func (d evaluationRepository) Save(tx *gorm.DB, data *evaluationModel.Evaluation) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d evaluationRepository) Delete(evaluationData evaluationModel.Evaluation) error {
	db := d.db.Delete(&evaluationData)
	return db.Error
}
