package handler

import (
	"net/http"
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/utils"
	"penilaian-360/internal/app/model/evaluationModel"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type EvaluationHandler struct {
	HandlerOption
}

func (evaluation EvaluationHandler) Evaluation(c *gin.Context) {
	var (
		record  = loggers.StartRecord(c.Request)
		request = evaluationModel.EvaluationRequest{}
	)
	if err := c.Bind(&request); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	if err := request.Validate(); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	res, err := evaluation.EvaluationService.SaveEvaluation(record, request)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}

func (evaluation EvaluationHandler) GetEvaluation(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	paging := datapaging.Datapaging{
		Page:  cast.ToInt(c.Query("page_number")),
		Limit: cast.ToInt(c.Query("page_size")),
	}

	res, count, err := evaluation.EvaluationService.EvaluationList(record, paging)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "Data not found")
		return
	}
	data := map[string]interface{}{
		"page_number":        paging.Page,
		"page_size":          paging.Limit,
		"total_record_count": count,
		"records":            res,
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, data, "Success")
}

func (evaluation EvaluationHandler) EvaluationView(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	res, err := evaluation.EvaluationService.EvaluationView(record, cast.ToInt64(c.Param("id")))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}

func (evaluation EvaluationHandler) EvaluationDelete(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	err := evaluation.EvaluationService.EvaluationDelete(record, cast.ToInt64(c.Param("id")))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, nil, "Success")
}
