package evaluatedEmployeesModel

type (
	EvaluatedEmployee struct {
		Id           int64   `json:"id"`
		EvaluationId int64   `json:"evaluation_id"`
		EmployeeId   int64   `json:"employee_id"`
		TotalAvg     float64 `json:"total_avg"`
	}
)

func (EvaluatedEmployee) TableName() string {
	return "evaluated_employees"
}
