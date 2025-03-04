package evaluationModel

import (
	"fmt"
	"penilaian-360/internal/app/commons/constants"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	EvaluationAnswer struct {
		Id                  int64  `json:"id"`
		EvaluationId        int64  `json:"evaluation_id"`
		EvaluatorEmployeeId int64  `json:"evaluator_employee_id"`
		QuestionId          int64  `json:"question_id"`
		Answer              string `json:"answer"`
		FinalPoint          int    `json:"final_point"`
	}

	// DTO

	EvaluationAnswerResponse struct {
		Type           string `json:"type"`
		CompetencyType string `json:"competency_type"`
		TitleName      string `json:"title_name"`
		QuestionName   string `json:"question_name"`
		EvaluationAnswer
	}

	EvaluationAnswerRequest struct {
		Type string `json:"type"`
		EvaluationAnswer
	}

	EvaluationAnswerRequests struct {
		Data []EvaluationAnswerRequest
	}
)

func (e *EvaluationAnswerRequests) Validate() error {
	for i, v := range e.Data {
		if err := validation.ValidateStruct(&v,
			validation.Field(&v.Id, validation.Required, validation.Min(1)),                  // Id harus lebih dari 0
			validation.Field(&v.EvaluationId, validation.Required, validation.Min(1)),        // EvaluationId harus lebih dari 0
			validation.Field(&v.EvaluatorEmployeeId, validation.Required, validation.Min(1)), // EvaluatorEmployeeId harus lebih dari 0
			validation.Field(&v.QuestionId, validation.Required, validation.Min(1)),          // QuestionId harus lebih dari 0
			validation.Field(&v.Type, validation.Required, validation.In(
				string(constants.QuestionTypeRate),
				string(constants.QuestionTypeEssay),
			)),
			validation.Field(&v.FinalPoint, validation.Min(0), validation.Max(5)), // FinalPoint harus antara 0 - 5
		); err != nil {
			return fmt.Errorf("validasi gagal pada indeks %d: %w", i, err)
		}
	}
	return nil
}
