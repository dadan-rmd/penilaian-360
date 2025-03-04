package evaluationModel

import (
	"penilaian-360/internal/app/commons/constants"
	"penilaian-360/internal/app/model/evaluatedEmployeesModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Evaluation struct {
	Id            int64                      `json:"id"`
	DepartementId int64                      `json:"departement_id"`
	Title         string                     `json:"title"`
	Status        constants.EvaluationStatus `json:"status"`
	DeadlineAt    string                     `json:"deadline_at"`
	CreatedAt     time.Time                  `json:"created_at" gorm:"autoCreateTime"`
}

func (Evaluation) TableName() string {
	return "evaluations"
}

// DTO
type (
	DataFormHistory struct {
		Id            int64  `json:"id"`
		DepartementId int64  `json:"departement_id"`
		Title         string `json:"title"`
		Status        string `json:"status"`
		DeadlineAt    string `json:"deadline_at"`
	}
	DataQuestion struct {
		Id       int64  `json:"id"`
		Title    string `json:"title"`
		Question string `json:"question"`
	}
	FormHistoryRequest struct {
		DataFormHistory
		Functional         []DataQuestion `json:"functional"`
		Personal           []DataQuestion `json:"personal"`
		Essay              []DataQuestion `json:"essay"`
		IdToDeleteQuestion []int64        `json:"id_to_delete_question"`
	}
	FormHistoryResponse struct {
		Id int64 `json:"id"`
	}
	FormHistoryList struct {
		Id              int64     `json:"id"`
		DepartementName string    `json:"department_name" gorm:"column:DepartmentName"`
		Title           string    `json:"title"`
		Status          string    `json:"status"`
		DeadlineAt      string    `json:"deadline_at"`
		CreatedAt       time.Time `json:"created_at"`
	}

	AssignmentRequest struct {
		Id          int64   `json:"id"`
		EvaluatedId []int64 `json:"evaluated_id"`
		EvaluatorId []int64 `json:"evaluator_id"`
		Cc          string  `json:"cc"`
	}
)

func (v *FormHistoryRequest) Validate() error {
	err := validation.ValidateStruct(v,
		validation.Field(&v.DeadlineAt, validation.Required, validation.Date(constants.DDMMYYYY)),
		validation.Field(&v.Title, validation.Required),
		validation.Field(&v.Status, validation.Required, validation.In(
			string(constants.EvaluationStatusDraft),
			string(constants.EvaluationStatusPublish),
		)),
	)
	if err != nil {
		return err
	}

	return nil
}

func (e *AssignmentRequest) ToEvaluatedEmployee() (entities []evaluatedEmployeesModel.EvaluatedEmployee) {
	for _, v := range e.EvaluatedId {
		entities = append(entities, evaluatedEmployeesModel.EvaluatedEmployee{
			EvaluationId: e.Id,
			EmployeeId:   v,
		})
	}
	return
}

func (e *AssignmentRequest) ToEvaluatorEmployee(evaluatedEmployeeId int64) (entities []evaluatorEmployeesModel.EvaluatorEmployee) {
	for _, v := range e.EvaluatorId {
		entities = append(entities, evaluatorEmployeesModel.EvaluatorEmployee{
			EvaluationId:        e.Id,
			EvaluatedEmployeeId: evaluatedEmployeeId,
			EmployeeId:          v,
			Cc:                  e.Cc,
		})
	}
	return
}
