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
		Question string `json:"question"`
		Type     string `json:"type"`
	}
	FormHistoryRequest struct {
		DataFormHistory
		Question           []DataQuestion `json:"question"`
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

	DetailForm struct {
		EmployeeName string  `json:"employee_name"`
		Departement  string  `json:"department"`
		Position     string  `json:"position"`
		Avg          float64 `json:"avg"`
		Status       float64 `json:"status"`
	}
	DetailFormResponse struct {
		DepartementName string `json:"department_name" gorm:"column:DepartmentName"`
		EmployeeId      int64  `json:"employee_id"`
		Data            int64  `json:"data"`
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
		validation.Field(&v.Question, validation.Required),
	)
	if err != nil {
		return err
	}

	for _, q := range v.Question {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (q *DataQuestion) Validate() error {
	return validation.ValidateStruct(q,
		validation.Field(&q.Question, validation.Required),
		validation.Field(&q.Type, validation.Required, validation.In(
			string(constants.QuestionTypeRate),
			string(constants.QuestionTypeEssay),
		)),
	)
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
