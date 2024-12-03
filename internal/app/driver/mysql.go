package driver

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql" // defines mysql driver used
)

// DBMysqlOption options for mysql connection
type DBMysqlOption struct {
	IsEnable             bool
	Host                 string
	Port                 int
	Username             string
	Password             string
	DBName               string
	AdditionalParameters string
	MaxOpenConns         int
	MaxIdleConns         int
	ConnMaxLifetime      time.Duration
}

// NewMysqlDatabase return gorp dbmap object with MySQL options param
func NewMysqlDatabase(option DBMysqlOption) (*gorm.DB, error) {
	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", option.Username, option.Password, option.Host, option.Port, option.DBName, option.AdditionalParameters)
	db, err := sql.Open("mysql", dbDsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(option.ConnMaxLifetime)
	db.SetMaxIdleConns(option.MaxIdleConns)
	db.SetMaxOpenConns(option.MaxOpenConns)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation(os.Getenv("LOG_TIME_ZONE"))
			return time.Now().In(ti)
		},
	})
	setMode := os.Getenv("SET_MODE")
	if setMode == "production" {
		gormDB.Config.Logger = logger.Default.LogMode(logger.Error)
	} else {
		gormDB.Config.Logger = logger.Default.LogMode(logger.Info)
	}

	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
