package healtyRepository

import "central-auth/internal/app/model/healtyModel"

type IHealtyRepository interface {
	FindAll() (*[]healtyModel.Healty, error)
}
