package utils

import (
	"fmt"

	"github.com/albrow/forms"
	"github.com/gin-gonic/gin"
	"github.com/inhies/go-bytesize"
	"github.com/spf13/cast"
)

type ValidateFormFileRule struct {
	Key              string
	ValidExtensions  []string
	IsRequired       *bool
	FileMaxSizeBytes int64
}

func ValidateFormFile(c *gin.Context, rules []ValidateFormFileRule) (*forms.Validator, error) {
	reqBody, err := forms.Parse(c.Request)
	if err != nil {
		return nil, err
	}

	val := reqBody.Validator()

	for _, r := range rules {
		if r.IsRequired == nil || *r.IsRequired {
			val.RequireFile(r.Key)
		}
		val.AcceptFileExts(r.Key, r.ValidExtensions...)
		//Check file size
		if r.FileMaxSizeBytes > 0 {
			fileHeader, err := c.FormFile(r.Key)
			if err != nil {
				return nil, err
			}
			if fileHeader.Size > r.FileMaxSizeBytes {
				val.AddError(r.Key, fmt.Sprintf("File size can't more than %s",
					bytesize.New(cast.ToFloat64(r.FileMaxSizeBytes))))
			}
		}
	}

	return val, nil
}
