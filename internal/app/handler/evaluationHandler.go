package handler

import (
	"net/http"
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/utils"
	"penilaian-360/internal/app/model/evaluationModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type EvaluationHandler struct {
	HandlerOption
}

func (evaluation EvaluationHandler) EvaluationList(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	paging := datapaging.Datapaging{
		Page:  cast.ToInt(c.Query("page_number")),
		Limit: cast.ToInt(c.Query("page_size")),
	}
	employee, err := evaluation.AuthMiddleware.GetEmployee(c)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}

	res, count, err := evaluation.EvaluationService.EvaluationList(record, paging, employee, evaluatorEmployeesModel.EvaluatorEmployeeParams{
		Departement: c.Query("departement"),
		Search:      c.Query("search"),
	})
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
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

func (evaluation EvaluationHandler) EvaluationDepartementList(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	paging := datapaging.Datapaging{
		Page:  cast.ToInt(c.Query("page_number")),
		Limit: cast.ToInt(c.Query("page_size")),
	}
	employee, err := evaluation.AuthMiddleware.GetEmployee(c)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}

	res, count, err := evaluation.EvaluationService.EvaluationWithDepartementList(record, paging, employee)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
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

func (evaluation EvaluationHandler) EvaluationNeeds(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	paging := datapaging.Datapaging{
		Page:  cast.ToInt(c.Query("page_number")),
		Limit: cast.ToInt(c.Query("page_size")),
	}
	employee, err := evaluation.AuthMiddleware.GetEmployee(c)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}

	res, count, err := evaluation.EvaluationService.EvaluationNeeds(record, paging, employee, c.Query("search"))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
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

func (evaluation EvaluationHandler) EvaluationDetail(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	paging := datapaging.Datapaging{
		Page:  cast.ToInt(c.Query("page_number")),
		Limit: cast.ToInt(c.Query("page_size")),
	}

	res, count, err := evaluation.EvaluationService.EvaluationDetail(record, paging, cast.ToInt64(c.Param("evaluated_id")), evaluatorEmployeesModel.EvaluatorEmployeeParams{
		Departement: c.Query("departement"),
		Search:      c.Query("search"),
	})
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
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

func (evaluation EvaluationHandler) ScoreDetail(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)

	res, err := evaluation.EvaluationService.ScoreDetail(record, cast.ToInt64(c.Query("evaluation_id")), cast.ToInt64(c.Query("evaluator_employee_id")))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}

func (evaluation EvaluationHandler) Score(c *gin.Context) {
	var (
		record  = loggers.StartRecord(c.Request)
		request = evaluationModel.EvaluationAnswerRequests{}
	)
	if err := c.Bind(&request); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	if err := request.Validate(); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}

	err := evaluation.EvaluationService.Score(record, request)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, nil, "Success")
}

func (evaluation EvaluationHandler) Approve(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)

	err := evaluation.EvaluationService.EvaluationApprove(record, cast.ToInt64(c.Param("evaluator_id")))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, nil, "Success")
}
