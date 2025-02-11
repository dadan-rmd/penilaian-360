package handler

import (
	"net/http"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/utils"
	"penilaian-360/internal/app/model/employeeModel"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	HandlerOption
}

func (employee EmployeeHandler) GetEmployee(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
		params = employeeModel.EmployeeParamas{}
	)
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	if err := params.Validate(); err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, err.Error(), "")
		return
	}
	res, err := employee.EmployeeService.GetEmployeeAll(record, params)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}
