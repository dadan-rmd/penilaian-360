package evaluationEmployeeRepository

import (
	"penilaian-360/internal/app/commons/constants"
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type evaluationEmployeeRepository struct {
	db *gorm.DB
}

func NewEvaluationEmployeeRepository(db *gorm.DB) IEvaluationEmployeeRepository {
	return &evaluationEmployeeRepository{db}
}

func (d evaluationEmployeeRepository) FindByID(id int64) (entity *evaluationModel.EvaluationEmployee, err error) {
	db := d.db.Preload(clause.Associations)
	err = db.First(&entity, "id=? ", id).Error
	return
}

func (d evaluationEmployeeRepository) Save(tx *gorm.DB, data *evaluationModel.EvaluationEmployee) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d evaluationEmployeeRepository) Delete(evaluationEmployeeData evaluationModel.EvaluationEmployee) error {
	db := d.db.Delete(&evaluationEmployeeData)
	return db.Error
}

func (d evaluationEmployeeRepository) FindByEvaluationId(evaluationId int64, typeEmployee string) (employeeId []int64, err error) {
	err = d.db.
		Model(&evaluationModel.EvaluationEmployee{}).
		Where(evaluationModel.EvaluationEmployee{
			EvaluationId: evaluationId,
			Type:         constants.EmployeeType(typeEmployee),
		}).
		First(&employeeId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []int64{}, nil
		}
		return
	}
	return
}
