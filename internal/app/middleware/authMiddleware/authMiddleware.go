package authMiddleware

import (
	"github.com/gin-gonic/gin"
)

type IAuthMiddleware interface {
	BasicAuthenticate() gin.HandlerFunc
}
