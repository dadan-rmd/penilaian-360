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

func (d evaluatedEmployeeRepository) FindByEvaluationIdAndEmployeeId(evaluationId, employeeId int64) (entity *evaluatedEmployeesModel.EvaluatedEmployee, err error) {
	err = d.db.Where(evaluatedEmployeesModel.EvaluatedEmployee{
		EvaluationId: evaluationId,
		EmployeeId:   employeeId,
	}).
		Find(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return
	}
	return
}

func (d evaluatedEmployeeRepository) RetrieveListWithPaging(paging datapaging.Datapaging, departement, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {
	slqSelect := `
		evaluated_employees.evaluation_id,
		evaluated_employees.id as evaluated_id,
		evaluator_employees.id as evaluator_id,
		evaluated_employees.evaluation_id,
		evaluated_employees.employee_id,
		evaluator_employees.total_functional,
		evaluator_employees.total_personal,
		evaluated_employees.total_avg,
		evaluator_employees.has_assessed,
		evaluator_employees.requires_assessment,
		evaluator_employees.status,
		master_karyawan.Name, 
		master_karyawan.Department, 
		master_karyawan.Position,
		CASE WHEN evaluator_employees.has_assessed THEN msc.label ELSE NULL END AS classification,
		'lihat-penilaian' as action
	`

	db := d.db.Model(&evaluatedEmployeesModel.EvaluatedEmployee{}).
		Select(slqSelect).
		Joins("JOIN master_karyawan on master_karyawan.id = evaluated_employees.employee_id").
		Joins("JOIN evaluator_employees ON evaluated_employees.id = evaluator_employees.evaluated_employee_id").
		// LEFT JOIN ke master_score_classifications dengan kondisi range fleksibel (MySQL)
		Joins(`LEFT JOIN master_score_classifications msc ON (
			(msc.min_score IS NULL AND msc.max_score >= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED))
			OR (msc.max_score IS NULL AND msc.min_score <= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED))
			OR (msc.min_score <= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED) AND msc.max_score >= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED))
		)`).
		Order("evaluated_employees.id desc").
		Group("evaluated_employees.employee_id")

	// departement filter (parameterized)
	if departement != "" {
		db = db.Where("master_karyawan.Department = ?", departement)
	}

	// parameterized search
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("master_karyawan.Name LIKE ? OR master_karyawan.Position LIKE ?", like, like)
	}

	// hitung total (perhatikan: Count setelah Group akan menghitung jumlah group)
	if err = db.Count(&count).Error; err != nil {
		return
	}

	// paging
	if paging.Page != 0 {
		pg := datapaging.New(paging.Limit, paging.Page, []string{})
		db = db.Offset(pg.GetOffset()).Limit(paging.Limit)
	}

	// scan ke struct (pastikan struct EvaluatorEmployeeList punya field Classification)
	err = db.Scan(&data).Error
	return
}

func (d evaluatedEmployeeRepository) RetrieveNeedsWithPaging(paging datapaging.Datapaging, employeeId int64, search string) (data []evaluatorEmployeesModel.EvaluatorEmployeeList, count int64, err error) {
	slqSelect := `
		evaluated_employees.evaluation_id,
		evaluated_employees.id as evaluated_id,
		evaluator_employees.id as evaluator_id,
		evaluated_employees.evaluation_id,
		evaluated_employees.employee_id,
		evaluator_employees.total_functional,
		evaluator_employees.total_personal,
		evaluated_employees.total_avg,
		evaluator_employees.has_assessed,
		evaluator_employees.requires_assessment,
		evaluator_employees.status,
		master_karyawan.Name, 
		master_karyawan.Department, 
		master_karyawan.Position,
		CASE WHEN evaluator_employees.has_assessed THEN msc.label ELSE NULL END AS classification,
		'beri-penilaian' as action
	`

	db := d.db.Model(&evaluatedEmployeesModel.EvaluatedEmployee{}).
		Select(slqSelect).
		Joins("JOIN evaluator_employees on evaluator_employees.evaluated_employee_id = evaluated_employees.id").
		Joins("JOIN master_karyawan on master_karyawan.id = evaluated_employees.employee_id").
		Joins(`LEFT JOIN master_score_classifications msc ON (
			(msc.min_score IS NULL AND msc.max_score >= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED))
			OR (msc.max_score IS NULL AND msc.min_score <= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED))
			OR (msc.min_score <= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED) AND msc.max_score >= CAST(ROUND(evaluated_employees.total_avg) AS SIGNED))
		)`).
		Where("evaluator_employees.employee_id = ?", employeeId).
		Order("evaluated_employees.id desc")

	if search != "" {
		like := "%" + search + "%"
		db = db.Where("master_karyawan.Name LIKE ? OR master_karyawan.Position LIKE ?", like, like)
	}

	// count total before limit
	if err = db.Count(&count).Error; err != nil {
		return
	}

	// paging
	if paging.Page != 0 {
		pg := datapaging.New(paging.Limit, paging.Page, []string{})
		db = db.Offset(pg.GetOffset()).Limit(paging.Limit)
	}

	err = db.Scan(&data).Error
	return
}

func (d evaluatedEmployeeRepository) UpdateAvg(tx *gorm.DB, id int64, totalAvg float64) (err error) {
	if tx != nil {
		return tx.Model(&evaluatedEmployeesModel.EvaluatedEmployee{}).
			Where("id = ?", id).
			Update("total_avg", totalAvg).Error
	} else {
		return d.db.Model(&evaluatedEmployeesModel.EvaluatedEmployee{}).
			Where("id = ?", id).
			Update("total_avg", totalAvg).Error

	}
}
