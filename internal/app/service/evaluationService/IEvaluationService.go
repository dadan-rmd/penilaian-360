package evaluationService

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/evaluationModel"
	"penilaian-360/internal/app/model/questionModel"
)

type IEvaluationService interface {
	EvaluationList(record *loggers.Data, paging datapaging.Datapaging) (res []evaluationModel.EvaluationList, count int64, err error)
	EvaluationView(record *loggers.Data, id int64) (res []questionModel.Question, err error)
	SaveEvaluation(record *loggers.Data, request evaluationModel.EvaluationRequest) (res evaluationModel.EvaluationResponse, err error)
	EvaluationDelete(record *loggers.Data, id int64) (err error)
}
