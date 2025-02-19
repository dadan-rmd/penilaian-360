package evaluationRepository

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
)

type evaluationRepository struct {
	db *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) IEvaluationRepository {
	return &evaluationRepository{db}
}

func (d evaluationRepository) GetWithPaging(paging datapaging.Datapaging) (data []evaluationModel.FormHistoryList, count int64, err error) {
	db := d.db.Model(&evaluationModel.Evaluation{}).
		Select("evaluations.*, master_department.DepartmentName").
		Joins("JOIN master_department on master_department.id = evaluations.departement_id").
		Order("evaluations.created_at desc")

	db.Count(&count)

	if paging.Page != 0 {
		pg := datapaging.New(paging.Limit, paging.Page, []string{})
		db = db.Offset(pg.GetOffset()).Limit(paging.Limit)
	}

	err = db.Scan(&data).Error
	return
}

func (d evaluationRepository) FindByID(id int64) (entity *evaluationModel.Evaluation, err error) {
	err = d.db.First(&entity, id).Error
	return
}
func (d evaluationRepository) FindDepartmentNameByID(id int64) (DepartmentName string, err error) {
	err = d.db.Model(&evaluationModel.Evaluation{}).
		Select("master_department.DepartmentName").
		Joins("Join master_department on master_department.id = evaluations.departement_id").
		Pluck("DepartmentName", &DepartmentName).Error
	return
}

func (d evaluationRepository) Save(tx *gorm.DB, data *evaluationModel.Evaluation) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d evaluationRepository) Delete(tx *gorm.DB, id int64) error {
	if tx != nil {
		return tx.Delete(&evaluationModel.Evaluation{}, id).Error
	} else {
		return d.db.Delete(&evaluationModel.Evaluation{}, id).Error
	}
}
