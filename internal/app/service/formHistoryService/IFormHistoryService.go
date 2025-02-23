package formHistoryService

import (
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/model/evaluationModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"
	"penilaian-360/internal/app/model/questionModel"
)

type IFormHistoryService interface {
	FormHistoryList(record *loggers.Data, paging datapaging.Datapaging) (res []evaluationModel.FormHistoryList, count int64, err error)
	FormHistoryView(record *loggers.Data, id int64) (res questionModel.QuestionWithEvaluation, err error)
	SaveFormHistory(record *loggers.Data, request evaluationModel.FormHistoryRequest) (res evaluationModel.FormHistoryResponse, err error)
	FormHistoryDelete(record *loggers.Data, id int64) (err error)
	FormHistoryAssignment(record *loggers.Data, request evaluationModel.AssignmentRequest) (err error)
	FormHistoryDetail(record *loggers.Data, paging datapaging.Datapaging, params evaluatorEmployeesModel.FormHistoryDetailParams) (res evaluatorEmployeesModel.FormHistoryDetailResponse, count int64, err error)
}
