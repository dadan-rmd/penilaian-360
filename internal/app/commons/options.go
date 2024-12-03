package commons

import (
	"central-auth/internal/app/appcontext"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Options common option for all object that needed
type Options struct {
	AppCtx    *appcontext.AppContext
	Db        *gorm.DB
	UUID      Iuuid
	Validator *validator.Validate
	OssClient *oss.Client
}
