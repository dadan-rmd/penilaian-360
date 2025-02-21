package employeeModel

import (
	"penilaian-360/internal/app/commons/constants"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

/* Table Definition */
type Employee struct {
	Id          int64  `json:"id" gorm:"column:id"`
	EmployeeId  string `json:"employee_id" gorm:"column:EmployeeId"`
	Name        string `json:"name" gorm:"column:Name"`
	FirstName   string `json:"first_name" gorm:"column:FirstName"`
	MiddleName  string `json:"middle_name" gorm:"column:MiddleName"`
	LastName    string `json:"last_name" gorm:"column:LastName"`
	Department  string `json:"department" gorm:"column:Department"`
	Position    string `json:"position" gorm:"column:Position"`
	Email       string `json:"email" gorm:"column:Email"`
	AccessToken string `json:"access_token" gorm:"column:AccessToken"`
}

func (Employee) TableName() string {
	return "master_karyawan"
}

// DTO
type (
	EmployeeResponse struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	EmployeeParamas struct {
		EvaluationId int64  `form:"evaluation_id"`
		Type         string `form:"type"`
		Departement  string `form:"departement"`
	}

	EmployedEmployeeResponse struct {
		Name        string `json:"name" gorm:"column:Name"`
		EvaluatedId int64  `form:"evaluated_id"`
	}
)

func (entity Employee) ToEmployeeResponse() EmployeeResponse {
	return EmployeeResponse{
		Id:   entity.Id,
		Name: entity.Name,
	}
}

func (v *EmployeeParamas) Validate() error {
	return validation.ValidateStruct(v,
		validation.Field(&v.EvaluationId, validation.Required),
		validation.Field(&v.Type, validation.Required, validation.In(
			string(constants.EmployeeTypeEvaluated),
			string(constants.EmployeeTypeEvaluator),
		)),
	)
}
