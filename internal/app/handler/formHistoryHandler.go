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

type FormHistoryHandler struct {
	HandlerOption
}

func (evaluation FormHistoryHandler) FormHistory(c *gin.Context) {
	var (
		record  = loggers.StartRecord(c.Request)
		request = evaluationModel.FormHistoryRequest{}
	)
	if err := c.Bind(&request); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	if err := request.Validate(); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	res, err := evaluation.FormHistoryService.SaveFormHistory(record, request)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}

func (evaluation FormHistoryHandler) GetFormHistory(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	paging := datapaging.Datapaging{
		Page:  cast.ToInt(c.Query("page_number")),
		Limit: cast.ToInt(c.Query("page_size")),
	}

	res, count, err := evaluation.FormHistoryService.FormHistoryList(record, paging)
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

func (evaluation FormHistoryHandler) FormHistoryView(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	res, err := evaluation.FormHistoryService.FormHistoryView(record, cast.ToInt64(c.Param("id")))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}

func (evaluation FormHistoryHandler) FormHistoryDelete(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	err := evaluation.FormHistoryService.FormHistoryDelete(record, cast.ToInt64(c.Param("id")))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, nil, "Success")
}

func (evaluation FormHistoryHandler) FormHistoryAssignment(c *gin.Context) {
	var (
		record  = loggers.StartRecord(c.Request)
		request = evaluationModel.AssignmentRequest{}
	)
	if err := c.Bind(&request); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	err := evaluation.FormHistoryService.FormHistoryAssignment(record, request)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, nil, "Success")
}

func (evaluation FormHistoryHandler) FormHistoryDetail(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	res, count, err := evaluation.FormHistoryService.FormHistoryDetail(record, datapaging.Datapaging{
		Page:  cast.ToInt(c.Query("page_number")),
		Limit: cast.ToInt(c.Query("page_size")),
	}, evaluatorEmployeesModel.FormHistoryDetailParams{
		ID:           cast.ToInt64(c.Query("form_id")),
		EmployeeIdID: cast.ToInt64(c.Query("employee_id")),
	})
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	data := map[string]interface{}{
		"page_number":        cast.ToInt(c.Query("page_number")),
		"page_size":          cast.ToInt(c.Query("page_size")),
		"total_record_count": count,
		"records":            res,
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, data, "Success")
}
