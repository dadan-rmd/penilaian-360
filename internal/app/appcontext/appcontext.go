package appcontext

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"penilaian-360/config"
	"penilaian-360/internal/app/driver"

	"gorm.io/gorm"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/rs/zerolog/log"
)

const (
	// DBDialectMysql rdbms dialect name for MySQL
	DBDialectMysql = "mysql"

	// DBDialectPostgres rdbms dialect name for PostgreSQL
	DBDialectPostgres = "postgres"
)

// AppContext the app context struct
type AppContext struct {
	config config.Provider
	Config Config
}

// NewAppContext initiate appcontext object
func NewAppContext(config config.Provider) *AppContext {
	var c Config
	err := config.Unmarshal(&c)
	if err != nil {
		panic(fmt.Sprint("failed mapping config ", err))
	}

	c.App.LocalIP = getLocalIP()
	return &AppContext{
		Config: c,
	}
}

// GetDBInstance getting gorp instance, param: dbType can be "mysql" or "postgre"
func (a *AppContext) GetDBInstance(dbType string) (*gorm.DB, error) {
	var gormDB *gorm.DB
	var err error
	if dbType == DBDialectMysql {
		dbOption := a.GetMysqlOption(dbType)
		gormDB, err = driver.NewMysqlDatabase(dbOption)
		if err != nil {
			return nil, err
		}
	} else {
		err = errors.New("Error get db instance, unknown db type")
	}

	return gormDB, err
}

// GetMysqlOption returns mysql options
func (a *AppContext) GetMysqlOption(dbType string) driver.DBMysqlOption {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 3306
	}

	dbConnMaxOpen, err := strconv.Atoi(os.Getenv("DB_CONN_MAX_OPEN"))
	if err != nil {
		dbConnMaxOpen = 5
	}

	dbConnMaxIdle, err := strconv.Atoi(os.Getenv("DB_CONN_MAX_IDLE"))
	if err != nil {
		dbConnMaxIdle = 5
	}

	dbConnMaxLifetime, err := strconv.Atoi(os.Getenv("DB_CONN_MAX_LIFETIME"))
	if err != nil {
		dbConnMaxLifetime = 120000000000
	}
	return driver.DBMysqlOption{
		IsEnable:             true,
		Host:                 os.Getenv("DB_HOST"),
		Port:                 dbPort,
		Username:             os.Getenv("DB_USERNAME"),
		Password:             os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_DATABASE_NAME"),
		AdditionalParameters: os.Getenv("DB_ADDITIONAL_PARAM"),
		MaxOpenConns:         dbConnMaxOpen,
		MaxIdleConns:         dbConnMaxIdle,
		ConnMaxLifetime:      time.Duration(dbConnMaxLifetime),
	}
}

func (a *AppContext) GetClient(clientName string) Client {
	var c Client
	var isFound bool
	for _, v := range a.Config.App.Clients {
		if v.Name == clientName {
			isFound = true
			c = v
		}
	}

	if !isFound {
		log.Fatal().Msgf("Client %s not found on config", clientName)
	}

	return c
}

func (a *AppContext) GetClientOss() (client *oss.Client) {
	client, err := oss.New(os.Getenv("OSS_ENDPOINT"), os.Getenv("OSS_ACCESS_KEY_ID"), os.Getenv("OSS_ACCESS_KEY_SECRET"))
	if err != nil {
		log.Fatal().Msgf("Client OSS not connect, %s", err)
	}
	return
}

// getLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal().Msgf("Fail get local IP | %s", err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	log.Fatal().Msg("local IP is empty")
	return ""
}
