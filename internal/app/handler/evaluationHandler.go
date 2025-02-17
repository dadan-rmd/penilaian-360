package handler

import (
	"net/http"
	datapaging "penilaian-360/internal/app/commons/dataPagingHelper"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/utils"
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
