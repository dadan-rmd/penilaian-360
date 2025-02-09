package employeeModel

import (
	"penilaian-360/internal/app/commons/constants"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

/* Table Definition */
type Employee struct {
	Id         int64  `json:"id" gorm:"id"`
	EmployeeId string `json:"employee_id" gorm:"EmployeeId"`
	Name       string `json:"name" gorm:"Name"`
	FirstName  string `json:"first_name" gorm:"FirstName"`
	MiddleName string `json:"middle_name" gorm:"MiddleName"`
	LastName   string `json:"last_name" gorm:"LastName"`
	Department string `json:"department" gorm:"Department"`
	Position   string `json:"position" gorm:"Position"`
}

func (Employee) TableName() string {
	return "master_karyawan"
}

// DTO
type (
	EmployeeResponse struct {
		Id   int64  `json:"id" gorm:"id"`
		Name string `json:"name"`
	}
	EmployeeParamas struct {
		EvaluationId int64  `form:"evaluation_id"`
		Type         string `form:"type"`
		Departement  string `form:"departement"`
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
		validation.Field(&v.Departement, validation.Required),
		validation.Field(&v.EvaluationId, validation.Required),
		validation.Field(&v.Type, validation.Required, validation.In(
			string(constants.EmployeeTypeEvaluated),
			string(constants.EmployeeTypeEvaluator),
		)),
	)
}
