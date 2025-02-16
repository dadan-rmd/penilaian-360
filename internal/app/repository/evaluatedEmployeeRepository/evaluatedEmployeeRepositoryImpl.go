package evaluatedEmployeeRepository

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/model/evaluatedEmployeesModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"

	"gorm.io/gorm"
)

type evaluatedEmployeeRepository struct {
	db *gorm.DB
}

func NewEvaluatedEmployeeRepository(db *gorm.DB) IEvaluatedEmployeeRepository {
	return &evaluatedEmployeeRepository{db}
}

func (d evaluatedEmployeeRepository) Save(tx *gorm.DB, data *[]evaluatedEmployeesModel.EvaluatedEmployee) error {
	if tx != nil {
		return tx.Save(&data).Error
	} else {
		return d.db.Save(&data).Error
	}
}

func (d evaluatedEmployeeRepository) FindEmployeeIdByEvaluationId(evaluationId int64) (employeeId []int64, err error) {
	err = d.db.
		Model(&evaluatedEmployeesModel.EvaluatedEmployee{}).
		Where(evaluatedEmployeesModel.EvaluatedEmployee{
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

func (d evaluatedEmployeeRepository) RetrieveListWithPaging(paging datapaging.Datapaging, departement, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {
	db := d.db.Model(&evaluatedEmployeesModel.EvaluatedEmployee{}).
		Select(`
			evaluated_employees.id,
			evaluated_employees.evaluation_id,
			evaluated_employees.total_avg,
			master_karyawan.Name, 
			master_karyawan.Department, 
			master_karyawan.Position
		`).
		Joins("JOIN master_karyawan on master_karyawan.id = evaluated_employees.employee_id").
		Order("evaluated_employees.id desc")
	if departement != "" {
		db.Where("master_karyawan.Department = ?", departement)
	}
	if search != "" {
		db.Where("master_karyawan.Name like '%" + search + "%' or master_karyawan.Position like '%" + search + "%'")
	}
	db.Count(&count)

	if paging.Page != 0 {
		pg := datapaging.New(paging.Limit, paging.Page, []string{})
		db = db.Offset(pg.GetOffset()).Limit(paging.Limit)
	}

	err = db.Scan(&data).Error
	return
}
