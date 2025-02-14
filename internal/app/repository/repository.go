package repository

import (
	"penilaian-360/internal/app/commons"
	"penilaian-360/internal/app/repository/departmentRepository"
	"penilaian-360/internal/app/repository/employeeRepository"
	"penilaian-360/internal/app/repository/evaluatedEmployeeRepository"
	"penilaian-360/internal/app/repository/evaluationAnswerRepository"
	"penilaian-360/internal/app/repository/evaluationRepository"
	"penilaian-360/internal/app/repository/evaluatorEmployeeRepository"
	"penilaian-360/internal/app/repository/questionRepository"
	"penilaian-360/internal/app/repository/userRepository"
)

// Option anything any repo object needed
type Option struct {
	commons.Options
}

type Repositories struct {
	UserRepository              userRepository.IUserRepository
	EvaluationRepository        evaluationRepository.IEvaluationRepository
	QuestionRepository          questionRepository.IQuestionRepository
	EvaluationAnswerRepository  evaluationAnswerRepository.IEvaluationAnswerRepository
	DepartmentRepository        departmentRepository.IDepartmentRepository
	EmployeeRepository          employeeRepository.IEmployeeRepository
	EvaluatorEmployeeRepository evaluatorEmployeeRepository.IEvaluatorEmployeeRepository
	EvaluatedEmployeeRepository evaluatedEmployeeRepository.IEvaluatedEmployeeRepository
}
