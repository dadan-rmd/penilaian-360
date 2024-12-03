package healtyService

import (
	"central-auth/internal/app/model/healtyModel"

	"central-auth/internal/app/commons/loggers"
)

type IHealtyService interface {
	FindAllHealty(record *loggers.Data) (*[]healtyModel.Healty, error)
}
