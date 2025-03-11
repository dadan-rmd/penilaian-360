package authMiddleware

import (
	"errors"
	"os"
	"penilaian-360/internal/app/commons/jsonHttpResponse"
	"penilaian-360/internal/app/commons/jwtHelper"
	"penilaian-360/internal/app/model/employeeModel"
	"penilaian-360/internal/app/repository/employeeRepository"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrUserNotFound = errors.New("user not found")
	ErrTokenRevoked = errors.New("token revoked")
)

type authMiddleware struct {
	employeeRepo employeeRepository.IEmployeeRepository
}

func NewAuthMiddleware(employeeRepo employeeRepository.IEmployeeRepository) IAuthMiddleware {
	return &authMiddleware{employeeRepo}
}

func (auth *authMiddleware) AuthorizeEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		bearerTokenSplit := strings.Split(bearerToken, " ")

		if len(bearerTokenSplit) < 2 {
			res := jsonHttpResponse.FailedResponse{
				Status:       jsonHttpResponse.FailedStatus,
				ResponseCode: "00",
				Message:      "invalid token",
			}
			jsonHttpResponse.Unauthorized(c, res)
			c.Abort()
			return
		}

		jwtToken := bearerTokenSplit[1]
		employee, err := auth.getUserFromJWTWithRoleValidation(jwtToken)
		if err != nil {
			if err == jwtHelper.ErrTokenExpired {
				res := jsonHttpResponse.FailedResponse{
					Status:       jsonHttpResponse.FailedStatus,
					ResponseCode: "00",
					Message:      err.Error(),
				}
				jsonHttpResponse.Unauthorized(c, res)
				c.Abort()
				return
			}

			res := jsonHttpResponse.FailedResponse{
				Status:       jsonHttpResponse.FailedStatus,
				ResponseCode: "00",
				Message:      "invalid token",
			}
			jsonHttpResponse.InternalServerError(c, res)
			c.Abort()
			return
		}
		//put into user context
		c.Set("employee", employee)
	}
}

func (auth *authMiddleware) getUserFromJWTWithRoleValidation(jwtToken string) (employee employeeModel.Employee, err error) {
	if jwtToken == "" {
		err = ErrInvalidToken
		return
	}
	jwtTokenClaims, err := jwtHelper.DecodeJWT(jwtToken)
	if err != nil {
		err = ErrInvalidToken
		return
	}
	employee, err = auth.employeeRepo.FindByEmailAndAccessToken(jwtTokenClaims.Data.Email, jwtTokenClaims.Data.AccessToken)
	if err != nil {
		err = ErrUserNotFound
		return
	}
	return
}

func (auth *authMiddleware) GetEmployee(c *gin.Context) (entity employeeModel.Employee, err error) {
	value, exists := c.Get("employee")
	if !exists {
		err = ErrUserNotFound
		return
	}
	entity, ok := value.(employeeModel.Employee)
	if !ok {
		err = ErrInvalidToken
		return
	}
	entity.Role = entity.Position
	whitelistUser := strings.Split(os.Getenv("WHITELIST_USER"), ",")
	accessRoleDepartement := strings.Split(os.Getenv("ACCESS_ROLE_DEPARTEMENT"), ",")
	entity.Role = "employee"
	if slices.Contains(accessRoleDepartement, entity.Department) {
		entity.Role = "hr"
	} else if slices.Contains(whitelistUser, entity.Email) {
		entity.Role = "head"
	}
	return
}
