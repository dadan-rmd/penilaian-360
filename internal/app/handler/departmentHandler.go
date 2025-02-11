package handler

import (
	"net/http"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/utils"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	HandlerOption
}

func (department DepartmentHandler) GetDepartment(c *gin.Context) {
	record := loggers.StartRecord(c.Request)
	res, err := department.DepartmentService.GetDepartmentAll(record)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}
