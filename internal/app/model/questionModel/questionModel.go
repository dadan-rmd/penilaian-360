package questionModel

import "penilaian-360/internal/app/commons/constants"

type (
	Question struct {
		Id           int64                  `json:"id"`
		EvaluationId int64                  `json:"evaluation_id"`
		Question     string                 `json:"question"`
		Type         constants.QuestionType `json:"type"`
	}

	QuestionWithDepartement struct {
		Id            int64                  `json:"id"`
		EvaluationId  int64                  `json:"evaluation_id"`
		DepartementId int64                  `json:"departement_id"`
		Question      string                 `json:"question"`
		Type          constants.QuestionType `json:"type"`
	}
)
