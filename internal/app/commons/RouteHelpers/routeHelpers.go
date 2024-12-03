package RouteHelpers

import (
	"central-auth/internal/app/model/userModel"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrorContextNotExist     = errors.New("user context not exist")
	ErrorParsingUserModel    = errors.New("error parsing user model")
	ErrInvalidJWTForProspect = errors.New("invalid jwt, forbidden access")
)

func GetUserFromJWTContext(c *gin.Context) (*userModel.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, ErrorContextNotExist
	}

	if user == "" {
		return nil, ErrorContextNotExist
	}
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, ErrorParsingUserModel
	}

	var userData userModel.User
	err = json.Unmarshal(jsonData, &userData)
	if err != nil {
		return nil, ErrorParsingUserModel
	}

	return &userData, nil
}
