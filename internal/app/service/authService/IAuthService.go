package authService

import (
	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/model/authModel"
	"central-auth/internal/app/model/userModel"
)

type IAuthService interface {
	Login(record *loggers.Data, loginReq authModel.LoginReq) (loginRes userModel.ResLogin, err error)
	ChangePasswordFromForgotPass(record *loggers.Data, request authModel.ChangePasswordFromForgotPassReq) (err error)
}
