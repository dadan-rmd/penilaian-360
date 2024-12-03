package userService

import (
	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/model/userModel"
)

type IUserService interface {
	CreateUser(record *loggers.Data, req userModel.CreateUserReq) (*userModel.User, int, error)
}
