package handler

import (
	"net/http"
	"penilaian-360/internal/app/commons/jsonHttpResponse"
	"penilaian-360/internal/app/commons/loggers"
	"penilaian-360/internal/app/commons/requestvalidationerror"
	"penilaian-360/internal/app/commons/utils"
	"penilaian-360/internal/app/model/authModel"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	HandlerOption
}

func (auth AuthHandler) Login(c *gin.Context) {
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

	// loginRes, err := auth.AuthService.Login(record, request)
	// if err != nil {
	// 	if err == authService.ErrInvalidCredential {
	// 		utils.BasicResponse(record, c.Writer, false, http.StatusUnauthorized, err.Error(), "")
	// 		return
	// 	}
	// 	utils.BasicResponse(record, c.Writer, false, http.StatusInternalServerError, err.Error(), "")
	// 	return
	// }
	// utils.BasicResponse(record, c.Writer, true, http.StatusOK, loginRes, "Success")
}
