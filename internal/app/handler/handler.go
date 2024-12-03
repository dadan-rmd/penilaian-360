package handler

import (
	"central-auth/internal/app/commons"
	"central-auth/internal/app/service"
)

// HandlerOption option for handler, including all service
type HandlerOption struct {
	commons.Options
	*service.Services
}
