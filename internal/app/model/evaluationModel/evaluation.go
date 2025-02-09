package evaluationModel

import (
	"time"
)

type Evaluation struct {
	Id            int64     `json:"id"`
	DepartementId int64     `json:"departement_id"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	Cc            string    `json:"cc"`
	DeadlineAt    string    `json:"deadline_at"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Evaluation) TableName() string {
	return "evaluatios"
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
		Question  []DataQuestion `json:"question"`
		IsDeleted []int64        `json:"is_deleted"`
	}

	AssignmentRequest struct {
		EvaluatedId            []int64 `json:"evaluated_id"`
		EvaluatedDepartementId []int64 `json:"evaluated_departement_id"`
		EvaluatorId            []int64 `json:"evaluator_id"`
		EvaluatorDepartementId []int64 `json:"evaluator_departement_id"`
		Cc                     string  `json:"cc"`
	}
)
