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

func (d questionRepository) FindWithDepartementByEvaluationId(evaluationId int64) (entity []questionModel.QuestionWithDepartement, err error) {
	err = d.db.Model(&questionModel.Question{}).
		Select("questions.*,evaluations.departement_id").
		Joins("JOIN evaluations on evaluations.id = questions.evaluation_id").
		Where("questions.evaluation_id = ?", evaluationId).
		Order("questions.id asc").
		Find(&entity).Error
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
