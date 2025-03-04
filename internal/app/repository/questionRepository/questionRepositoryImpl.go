package questionRepository

import (
	"penilaian-360/internal/app/model/questionModel"

	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) IQuestionRepository {
	return &questionRepository{db}
}

func (d questionRepository) FindByID(id int64) (entity *questionModel.Question, err error) {
	err = d.db.First(&entity, "id=? ", id).Error
	return
}

func (d questionRepository) FindByEvaluationId(evaluationId int64) (entity []questionModel.Question, err error) {
	err = d.db.
		Where("evaluation_id = ?", evaluationId).
		Order("id asc").
		Find(&entity).Error
	return
}

func (d questionRepository) FindByEvaluationIdAndType(evaluationId int64, typeQuestion, competencyType string) (entity []questionModel.Question, err error) {
	db := d.db.Where("evaluation_id = ?", evaluationId)
	if typeQuestion != "" {
		db = db.Where("type=? ", typeQuestion)
	}
	if competencyType != "" {
		db = db.Where("competency_type=?", competencyType)
	}
	err = db.Order("id asc").Find(&entity).Error
	return
}

func (d questionRepository) Save(tx *gorm.DB, data *[]questionModel.Question) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d questionRepository) Delete(tx *gorm.DB, id []int64) error {
	if tx != nil {
		return tx.Delete(&questionModel.Question{}, id).Error
	} else {
		return d.db.Delete(&questionModel.Question{}, id).Error
	}
}

func (d questionRepository) DeleteEvaluationId(tx *gorm.DB, evaluationId int64) error {
	if tx != nil {
		return tx.Where("evaluation_id = ?", evaluationId).Delete(&questionModel.Question{}).Error
	} else {
		return d.db.Where("evaluation_id = ?", evaluationId).Delete(&questionModel.Question{}).Error
	}
}

func (d questionRepository) CountRate(tx *gorm.DB, ids []int64) (count int64, err error) {
	if tx != nil {
		err = tx.Model(&questionModel.Question{}).Where("id in (?) and type = 'rate'", ids).Count(&count).Error
		return
	} else {
		err = d.db.Model(&questionModel.Question{}).Where("id in (?) and type = 'rate'", ids).Count(&count).Error
		return
	}
}

func (d questionRepository) CountRateByEvaluationIdAndType(tx *gorm.DB, evaluationId int64, typeQuestion, competencyType string) (count int64, err error) {
	var db *gorm.DB
	if tx != nil {
		db = tx.Model(&questionModel.Question{})
	} else {
		db = d.db.Model(&questionModel.Question{})
	}
	db = db.Where("evaluation_id = ?", evaluationId)
	if typeQuestion != "" {
		db = db.Where("type=? ", typeQuestion)
	}
	if competencyType != "" {
		db = db.Where("competency_type=?", competencyType)
	}
	err = db.Count(&count).Error
	return
}
