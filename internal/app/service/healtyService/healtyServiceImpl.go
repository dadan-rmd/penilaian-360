package healtyService

import (
	"errors"

	"central-auth/internal/app/model/healtyModel"
	"central-auth/internal/app/repository/healtyRepository"

	"central-auth/internal/app/commons/loggers"
)

var (
	ErrHealtyNotFound = errors.New("healty not found")
)

type healtyService struct {
	healtyRepo healtyRepository.IHealtyRepository
}

func NewHealtyService(healty healtyRepository.IHealtyRepository) IHealtyService {
	return &healtyService{healtyRepo: healty}
}

func (h healtyService) FindAllHealty(record *loggers.Data) (*[]healtyModel.Healty, error) {
	loggers.Logf(record, "Info, FindAllHealty")
	data, err := h.healtyRepo.FindAll()
	if err != nil {
		return nil, ErrHealtyNotFound
	}
	return data, err
}
