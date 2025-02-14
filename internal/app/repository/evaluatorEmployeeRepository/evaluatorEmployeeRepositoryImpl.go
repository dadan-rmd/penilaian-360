package evaluatorEmployeeRepository

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"
	"time"

	"gorm.io/gorm"
)

type evaluatorEmployeeRepository struct {
	db *gorm.DB
}

func NewEvaluatorEmployeeRepository(db *gorm.DB) IEvaluatorEmployeeRepository {
	return &evaluatorEmployeeRepository{db}
}
func (d evaluatorEmployeeRepository) FindEmployeeIdByEvaluationId(evaluationId int64) (employeeId []int64, err error) {
	err = d.db.
		Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
		Where(evaluatorEmployeesModel.EvaluatorEmployee{
			EvaluationId: evaluationId,
		}).
		Pluck("employee_id", &employeeId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []int64{}, nil
		}
		return
	}
	return
}

func (d evaluatorEmployeeRepository) UpdateEmailSentByEvaluatedEmployeeIdAndEmployeeId(ids, evaluatedEmployeeId []int64) error {
	return d.db.Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
		Where("employee_id IN ? and evaluated_employee_id IN ?", ids, evaluatedEmployeeId).
		Update("email_sent", time.Now().Format("2006-01-02 15:04:05")).Error
}

func (d evaluatorEmployeeRepository) FindByEvaluatorId(paging datapaging.Datapaging, evaluationId, evaluatedEmployeeId int64) (entities []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {
	db := d.db.
		Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
		Select("evaluator_employees.employee_id,evaluator_employees.avg,master_karyawan.Name,master_karyawan.Department,master_karyawan.Position").
		Joins("Join master_karyawan on master_karyawan.id = evaluator_employees.employee_id").
		Joins("Join evaluated_employees on evaluated_employees.id = evaluator_employees.evaluated_employee_id").
		Where("evaluated_employees.employee_id = ?", evaluatedEmployeeId).
		Where(evaluatorEmployeesModel.EvaluatorEmployee{
			EvaluationId: evaluationId,
		}).
		Order("evaluator_employees.id desc").
		Count(&count)
	if paging.Page != 0 {
		pg := datapaging.New(paging.Limit, paging.Page, []string{})
		db.Offset(pg.GetOffset()).Limit(paging.Limit)
	}
	err = db.Find(&entities).Error
	if err != nil {
		return
	}
	return
}

func (d evaluatorEmployeeRepository) Save(tx *gorm.DB, data *[]evaluatorEmployeesModel.EvaluatorEmployee) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}
