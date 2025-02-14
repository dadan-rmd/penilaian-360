package authMiddleware

import (
	"penilaian-360/internal/app/model/employeeModel"

	"github.com/gin-gonic/gin"
)

type IAuthMiddleware interface {
	AuthorizeEmployee() gin.HandlerFunc
	GetEmployee(c *gin.Context) (entity employeeModel.Employee, err error)
}
