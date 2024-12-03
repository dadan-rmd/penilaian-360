package driver

import (
	"gorm.io/gorm"
	"time"
)

// DBPostgreOption options for postgre connection
type DBPostgreOption struct {
	IsEnable          bool
	Host              string
	Port              int
	Username          string
	Password          string
	DBName            string
	MaxPoolSize       int
	ConnMaxLifetime   time.Duration
	MaxIdleConnection int
}

// NewPostgreDatabase return gorp dbmap object with postgre options param
func NewPostgreDatabase(option DBPostgreOption) (*gorm.DB, error) {
	panic("implement me please!")
}
