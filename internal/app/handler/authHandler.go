package handler

import (
	"central-auth/internal/app/commons/jsonHttpResponse"
	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/commons/requestvalidationerror"
	"central-auth/internal/app/commons/utils"
	"central-auth/internal/app/model/authModel"
	"central-auth/internal/app/service/authService"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	HandlerOption
}

func (authDelivery AuthHandler) Login(c *gin.Context) {
	record := loggers.StartRecord(c.Request)
	var request authModel.LoginReq
	errBind := c.ShouldBind(&request)
	if errBind != nil {
		validations := requestvalidationerror.GetvalidationError(errBind)

		if len(validations) > 0 {
			jsonHttpResponse.NewFailedMissingRequiredFieldResponse(c, validations)
			return
		}
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, errBind.Error(), "")
		return
	}

	loginRes, err := authDelivery.AuthService.Login(record, request)
	if err != nil {
		if err == authService.ErrInvalidCredential {
			utils.BasicResponse(record, c.Writer, false, http.StatusUnauthorized, err.Error(), "")
			return
		}
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, loginRes, "Success")
}

func (authDelivery AuthHandler) ForgotPass(c *gin.Context) {
	record := loggers.StartRecord(c.Request)

	var request authModel.ChangePasswordFromForgotPassReq
	errBind := c.ShouldBind(&request)
	if errBind != nil {
		validations := requestvalidationerror.GetvalidationError(errBind)

		if len(validations) > 0 {
			jsonHttpResponse.NewFailedMissingRequiredFieldResponse(c, validations)
			return
		}
		utils.BasicResponse(record, c.Writer, false, http.StatusBadRequest, errBind.Error(), "")
		return
	}

	err := authDelivery.AuthService.ChangePasswordFromForgotPass(record, request)
	if err != nil {
		utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
		return
	}
	utils.BasicResponse(record, c.Writer, true, http.StatusOK, nil, "password changed successfully")
}
