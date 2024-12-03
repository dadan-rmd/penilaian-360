package service

import (
	"central-auth/internal/app/commons"
	"central-auth/internal/app/middleware/authMiddleware"
	"central-auth/internal/app/repository"
	"central-auth/internal/app/service/authService"
	"central-auth/internal/app/service/healtyService"
	"central-auth/internal/app/service/userService"
)

// Option anything any service object needed
type Option struct {
	commons.Options
	*repository.Repositories
}

type Services struct {
	HealtyService  healtyService.IHealtyService
	AuthService    authService.IAuthService
	AuthMiddleware authMiddleware.IAuthMiddleware
	UserService    userService.IUserService
}
