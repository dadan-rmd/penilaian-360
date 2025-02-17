package evaluatorEmployeesModel

import (
	"time"
)

type (
	EvaluatorEmployee struct {
		Id                  int64     `json:"id"`
		EvaluationId        int64     `json:"evaluation_id"`
		EvaluatedEmployeeId int64     `json:"evaluated_employee_id"`
		EmployeeId          int64     `json:"employee_id"`
		Avg                 float64   `json:"avg"`
		EmailSent           time.Time `json:"email_sent"`
		Cc                  string    `json:"cc"`
	}

	// DTO
	EvaluatorEmployeeList struct {
		Id           int64   `json:"id"`
		EmployeeName string  `json:"employee_name" gorm:"column:Name"`
		Department   string  `json:"department" gorm:"column:Department"`
		Position     string  `json:"position" gorm:"column:Position"`
		Avg          float64 `json:"avg"`
		Status       string  `json:"status"`
	}
	EvaluatorEmployeeParams struct {
		Departement string
		Search      string
	}
	FormHistoryDetailResponse struct {
		Department string                  `json:"department" `
		Data       []EvaluatorEmployeeList `json:"data"`
	}

	FormHistoryDetailParams struct {
		ID           int64
		EmployeeIdID int64
	}
)
