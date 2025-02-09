package handler

import (
	"penilaian-360/internal/app/commons"
	"penilaian-360/internal/app/service"
)

// HandlerOption option for handler, including all service
type HandlerOption struct {
	commons.Options
	*service.Services
}
