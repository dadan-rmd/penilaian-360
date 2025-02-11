package evaluationModel

import (
	"penilaian-360/internal/app/commons/constants"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Evaluation struct {
	Id            int64                      `json:"id"`
	DepartementId int64                      `json:"departement_id"`
	Title         string                     `json:"title"`
	Status        constants.EvaluationStatus `json:"status"`
	Cc            string                     `json:"cc"`
	DeadlineAt    string                     `json:"deadline_at"`
	CreatedAt     time.Time                  `json:"created_at" gorm:"autoCreateTime"`
}

func (Evaluation) TableName() string {
	return "evaluations"
}

// DTO
type (
	DataEvaluation struct {
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
	EvaluationRequest struct {
		DataEvaluation
		Question           []DataQuestion `json:"question"`
		IdToDeleteQuestion []int64        `json:"id_to_delete_question"`
	}
	EvaluationResponse struct {
		Id int64 `json:"id"`
	}
	EvaluationList struct {
		Id              int64     `json:"id"`
		DepartementName string    `json:"department_name" gorm:"column:DepartmentName"`
		Title           string    `json:"title"`
		Status          string    `json:"status"`
		DeadlineAt      string    `json:"deadline_at"`
		CreatedAt       time.Time `json:"created_at"`
	}

	AssignmentRequest struct {
		EvaluatedId            []int64 `json:"evaluated_id"`
		EvaluatedDepartementId []int64 `json:"evaluated_departement_id"`
		EvaluatorId            []int64 `json:"evaluator_id"`
		EvaluatorDepartementId []int64 `json:"evaluator_departement_id"`
		Cc                     string  `json:"cc"`
	}
)

func (v *EvaluationRequest) Validate() error {
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
