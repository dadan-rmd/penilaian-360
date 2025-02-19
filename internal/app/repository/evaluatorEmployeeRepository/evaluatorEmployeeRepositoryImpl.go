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

func (d evaluatorEmployeeRepository) FindByID(tx *gorm.DB, id int64) (entity *evaluatorEmployeesModel.EvaluatorEmployee, err error) {
	if tx != nil {
		err = tx.First(&entity, id).Error
		return
	} else {
		err = d.db.First(&entity, id).Error
		return
	}
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
		Select(`
			evaluator_employees.id,
			evaluator_employees.employee_id,
			evaluator_employees.avg,
			master_karyawan.Name,
			master_karyawan.Department,
			master_karyawan.Position,
			CASE 
				WHEN evaluator_employees.avg > 0 THEN 'Sudah Menilai' 
				ELSE 'Belum Menilai' 
			END AS status
		`).
		Joins("JOIN master_karyawan ON master_karyawan.id = evaluator_employees.employee_id").
		Joins("JOIN evaluated_employees ON evaluated_employees.id = evaluator_employees.evaluated_employee_id").
		Where("evaluated_employees.employee_id = ?", evaluatedEmployeeId).
		Where(evaluatorEmployeesModel.EvaluatorEmployee{
			EvaluationId: evaluationId,
		}).
		Order("evaluator_employees.id DESC").
		Count(&count)

	if paging.Page != 0 {
		pg := datapaging.New(paging.Limit, paging.Page, []string{})
		db = db.Offset(pg.GetOffset()).Limit(paging.Limit)
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

func (d evaluatorEmployeeRepository) RetrieveListWithPaging(paging datapaging.Datapaging, employeeId int64, email, notDepartement, departement, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {
	slqSelect := `
				evaluated_employees.evaluation_id,
				evaluated_employees.id as evaluated_id,
				evaluator_employees.id as evaluator_id,
				evaluated_employees.evaluation_id,
				evaluated_employees.total_avg,
				master_karyawan.Name, 
				master_karyawan.Department, 
				master_karyawan.Position`
	db := d.db.Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
		Where("evaluator_employees.employee_id = ?", employeeId).
		Joins("JOIN evaluated_employees ON evaluated_employees.id = evaluator_employees.evaluated_employee_id").
		Joins("JOIN master_karyawan on master_karyawan.id = evaluated_employees.employee_id").
		Order("evaluated_employees.id desc")
	if notDepartement != "" {
		if email != "" {
			db.Select(slqSelect + `,
				CASE 
					WHEN evaluator_employees.cc = "` + email + `" THEN 'lihat-penilaian' 
					ELSE NULL 
				END AS status
			`)
		} else {
			db.Select(slqSelect)
		}
		db.Where("master_karyawan.Department != ?", notDepartement)
	} else {
		db.Select(slqSelect)
	}
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

func (d evaluatorEmployeeRepository) RetrieveEvaluatorDetailWithPaging(paging datapaging.Datapaging, evaluatedId int64, departement, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {
	db := d.db.Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
		Select(`
			evaluated_employees.evaluation_id,
			evaluated_employees.id as evaluated_id,
			evaluator_employees.id as evaluator_id,
			evaluator_employees.evaluation_id,
			evaluator_employees.avg as total_avg,
			master_karyawan.Name, 
			master_karyawan.Department, 
			master_karyawan.Position,
			'baca-penilaian' as status
		`).
		Where("evaluator_employees.evaluated_employee_id = ?", evaluatedId).
		Joins("JOIN master_karyawan on master_karyawan.id = evaluator_employees.employee_id").
		Joins("JOIN evaluated_employees ON evaluated_employees.id = evaluator_employees.evaluated_employee_id").
		Order("evaluator_employees.id desc")
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

func (d evaluatorEmployeeRepository) TotalAvg(tx *gorm.DB, evaluatedEmployeeId int64) (totalAvg float64, err error) {
	if tx != nil {
		err = tx.Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
			Select("SUM(avg)/count(evaluation_id) as total_avg").
			Where("evaluated_employee_id = ?", evaluatedEmployeeId).
			Pluck("total_avg", &totalAvg).Error
		return
	} else {
		err = d.db.Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
			Select("SUM(avg)/count(evaluation_id) as total_avg").
			Where("evaluated_employee_id = ?", evaluatedEmployeeId).
			Pluck("total_avg", &totalAvg).Error
		return
	}
}

func (d evaluatorEmployeeRepository) UpdateAvg(tx *gorm.DB, id int64, avg float64) (err error) {
	if tx != nil {
		return tx.Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
			Where("id = ?", id).
			Update("avg", avg).Error
	} else {
		return d.db.Model(&evaluatorEmployeesModel.EvaluatorEmployee{}).
			Where("id = ?", id).
			Update("avg", avg).Error

	}
}
