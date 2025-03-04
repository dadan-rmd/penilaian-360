package evaluatorEmployeesModel

type (
	EvaluatorEmployee struct {
		Id                  int64   `json:"id"`
		EvaluationId        int64   `json:"evaluation_id"`
		EvaluatedEmployeeId int64   `json:"evaluated_employee_id"`
		EmployeeId          int64   `json:"employee_id"`
		TotalFunctional     float64 `json:"total_functional"`
		TotalPersonal       float64 `json:"total_personal"`
		TotalAvg            float64 `json:"total_avg"`
		HasAssessed         bool    `json:"has_assessed"`
		RequiresAssessment  bool    `json:"requires_assessment"`
		EmailSent           string  `json:"email_sent"`
		Cc                  string  `json:"cc"`
	}

	// DTO
	EvaluatorEmployeeList struct {
		EvaluationId       int64   `json:"evaluation_id"`
		EvaluatedId        int64   `json:"evaluated_id"`
		EvaluatorId        int64   `json:"evaluator_id"`
		EmployeeName       string  `json:"employee_name" gorm:"column:Name"`
		Department         string  `json:"department" gorm:"column:Department"`
		Position           string  `json:"position" gorm:"column:Position"`
		TotalFunctional    float64 `json:"total_functional"`
		TotalPersonal      float64 `json:"total_personal"`
		TotalAvg           float64 `json:"total_avg"`
		HasAssessed        bool    `json:"has_assessed"`
		RequiresAssessment bool    `json:"requires_assessment"`
		Status             string  `json:"status"`
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
