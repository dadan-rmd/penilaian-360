package repository

import (
	"central-auth/internal/app/commons"
	"central-auth/internal/app/repository/healtyRepository"
	"central-auth/internal/app/repository/platformRepository"
	"central-auth/internal/app/repository/userRepository"
)

// Option anything any repo object needed
type Option struct {
	commons.Options
}

type Repositories struct {
	HealtyRepository   healtyRepository.IHealtyRepository
	UserRepository     userRepository.IUserRepository
	PlatformRepository platformRepository.IPlatformRepository
}
