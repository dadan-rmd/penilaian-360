package jsonHttpResponse

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func ErrorWithHeader(c *gin.Context, err error) {
	//switch error content
	if strings.Contains(err.Error(), "[Err 4") {
		NewFailedBadRequestResponse(c, err.Error())
		return
	}

	if strings.Contains(err.Error(), "[Err 5") {
		NewFailedInternalServerResponse(c, err.Error())
		return
	}

	NewFailedInternalServerResponse(c, err.Error())
	return
}
