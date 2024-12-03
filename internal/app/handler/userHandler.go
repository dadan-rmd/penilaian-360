package handler

import (
	"central-auth/internal/app/commons/jsonHttpResponse"
	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/commons/requestvalidationerror"
	"central-auth/internal/app/commons/utils"
	"central-auth/internal/app/model/userModel"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	HandlerOption
}

func (userDelivery UserHandler) AddUser(c *gin.Context) {
	record := loggers.StartRecord(c.Request)

	var request userModel.CreateUserReq
	errBind := c.ShouldBind(&request)
	if errBind != nil {
		validations := requestvalidationerror.GetvalidationError(errBind)

		if len(validations) > 0 {
			loggers.EndRecord(record, errBind.Error(), http.StatusBadRequest)
			jsonHttpResponse.NewFailedMissingRequiredFieldResponse(c, validations)
			return
		}

		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, userModel.MsgFailedAddUser, errBind.Error())
		return
	}
	if err := request.Validate(); err != nil {
		loggers.EndRecord(record, err.Error(), http.StatusBadRequest)
		jsonHttpResponse.NewFailedMissingRequiredFieldResponse(c, err)
		return
	}

	user, statusCode, err := userDelivery.UserService.CreateUser(record, request)
	if err != nil || statusCode == http.StatusInternalServerError {
		utils.BasicResponse(record, c.Writer, false, statusCode, err.Error(), "")
		return
	}

	utils.BasicResponse(record, c.Writer, true, http.StatusOK, user, "Registrasi berhasil")
}
