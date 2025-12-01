package handler

import (
	"fmt"
	"net/http"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/utils"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/model/evaluatorEmployeesModel"

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

func (employee EmployeeHandler) CreateToken(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	res, err := employee.EmployeeService.CreateToken(record, c.Query("email"), c.Query("access_token"))
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}

func (employee EmployeeHandler) GetUser(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)
	res, err := employee.AuthMiddleware.GetEmployee(c)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}

	utils.BasicResponse(record, c.Writer, true, http.StatusOK, res, "Success")
}

func (employee EmployeeHandler) GetEmployeeEmails(c *gin.Context) {
	var (
		record = loggers.StartRecord(c.Request)
	)

	emails, err := employee.EmployeeService.GetEmployeeEmails(record, evaluatorEmployeesModel.EvaluatorEmployeeParams{
		Search: c.Query("search"),
	})
	if err != nil {
		loggers.Logf(record, fmt.Sprintf("Err, GetEmployeeEmails service: %v", err))
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}

	// success
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, emails, "Success")
}
