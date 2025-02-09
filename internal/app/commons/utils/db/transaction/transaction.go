package transaction

import (
	"context"
	"gorm.io/gorm"
)

type transactionImpl struct {
	db *gorm.DB
}
type ITransaction interface {
	WithTransaction(ctx context.Context, fc func(tx *gorm.DB) error) error
}

func NewTransaction(db *gorm.DB) ITransaction {
	return &transactionImpl{
		db: db,
	}
}
func (s *transactionImpl) WithTransaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	return s.db.WithContext(ctx).Transaction(fc)
}
