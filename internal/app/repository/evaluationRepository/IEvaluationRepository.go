package evaluationRepository

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/model/evaluationModel"

	"gorm.io/gorm"
)

type IEvaluationRepository interface {
	GetWithPaging(paging datapaging.Datapaging) (data []evaluationModel.FormHistoryList, count int64, err error)
	FindByID(id int64) (entity *evaluationModel.Evaluation, err error)
	Save(tx *gorm.DB, data *evaluationModel.Evaluation) error
	Delete(tx *gorm.DB, id int64) error
	FindDepartmentNameByID(id int64) (DepartmentName string, err error)
}
