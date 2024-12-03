package authMiddleware

import (
	"central-auth/internal/app/commons/jsonHttpResponse"
	"central-auth/internal/app/repository/userRepository"
	"os"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	userRepo userRepository.IUserRepository
}

func NewAuthMiddleware(userRepo userRepository.IUserRepository) IAuthMiddleware {
	return &authMiddleware{userRepo}
}

func (auth *authMiddleware) BasicAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			res := jsonHttpResponse.FailedResponse{
				Status:       jsonHttpResponse.FailedStatus,
				ResponseCode: "00",
				Message:      "invalid basic auth credentials",
			}
			jsonHttpResponse.Unauthorized(c, res)
			c.Abort()
			return
		}

		isValid := (username == os.Getenv("AUTH_BASIC_USERNAME")) && (password == os.Getenv("AUTH_BASIC_PASSWORD"))
		if !isValid {
			res := jsonHttpResponse.FailedResponse{
				Status:       jsonHttpResponse.FailedStatus,
				ResponseCode: "00",
				Message:      "invalid basic auth credentials",
			}
			jsonHttpResponse.Unauthorized(c, res)
			c.Abort()
			return
		}
	}
}
