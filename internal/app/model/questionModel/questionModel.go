package questionModel

import (
	"penilaian-360/internal/app/commons/constants"
	"penilaian-360/internal/app/model/evaluationModel"
)

type (
	Question struct {
		Id           int64                  `json:"id"`
		EvaluationId int64                  `json:"evaluation_id"`
		Question     string                 `json:"question"`
		Type         constants.QuestionType `json:"type"`
	}

	QuestionWithEvaluation struct {
		evaluationModel.Evaluation
		Questions []Question `json:"questions"`
	}
)
