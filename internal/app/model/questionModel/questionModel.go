package questionModel

import (
	"penilaian-360/internal/app/commons/constants"
	"penilaian-360/internal/app/model/evaluationModel"
)

type (
	Question struct {
		Id             int64                      `json:"id"`
		EvaluationId   int64                      `json:"evaluation_id"`
		Title          string                     `json:"title"`
		Question       string                     `json:"question"`
		Type           constants.QuestionType     `json:"type"`
		CompetencyType constants.TypeOfCompetency `json:"competency_type"`
	}

	DataQuestion struct {
		Id           int64  `json:"id"`
		EvaluationId int64  `json:"evaluation_id"`
		Title        string `json:"title"`
		Question     string `json:"question"`
	}

	QuestionWithEvaluation struct {
		evaluationModel.Evaluation
		Functional []DataQuestion `json:"functional"`
		Personal   []DataQuestion `json:"personal"`
		Essay      []DataQuestion `json:"essay"`
	}
)

func ToAssemblerQuestion(questions []Question) (entities []DataQuestion) {
	for _, v := range questions {
		entities = append(entities, DataQuestion{
			Id:           v.Id,
			EvaluationId: v.EvaluationId,
			Title:        v.Title,
			Question:     v.Question,
		})
	}
	return
}
