package service

import (
	"penilaian-360/internal/app/commons"
	"penilaian-360/internal/app/middleware/authMiddleware"
	"penilaian-360/internal/app/repository"
	"penilaian-360/internal/app/service/authService"
	"penilaian-360/internal/app/service/departmentService"
	"penilaian-360/internal/app/service/employeeService"
	"penilaian-360/internal/app/service/evaluationService"
)

// Option anything any service object needed
type Option struct {
	commons.Options
	*repository.Repositories
}

type Services struct {
	AuthService       authService.IAuthService
	AuthMiddleware    authMiddleware.IAuthMiddleware
	DepartmentService departmentService.IDepartmentService
	EmployeeService   employeeService.IEmployeeService
	EvaluationService evaluationService.IEvaluationService
}
