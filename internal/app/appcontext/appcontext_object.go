package appcontext

import "time"

const (
	ClientIndomaret  = "INDOMARET"
	ClientSkyParking = "SKYPARKING"
	ClientAlfamart   = "ALFAMART"
	ClientSepulsa    = "SEPULSA"
	ScopeInternal    = "internal"
	ScopeIndomaret   = "indomaret"
	ScopeAlfamart    = "alfamart"
)

type Config struct {
	App   App   `mapstructure:"app"`
	MySQL MySQL `mapstructure:"mysql"`
}

type App struct {
	Host                     string        `mapstructure:"host"`
	Port                     int           `mapstructure:"port"`
	Name                     string        `mapstructure:"name"`
	LogLevel                 string        `mapstructure:"log_level"`
	Clients                  []Client      `mapstructure:"clients"`
	Errors                   []Error       `mapstructure:"errors"`
	TokenExpired             time.Duration `mapstructure:"token_expired"`
	TransactionExpired       time.Duration `mapstructure:"transaction_expired"`
	LocalIP                  string        //SET on init config
	MigrationPath            string        `mapstructure:"migration_path"`
	JWTSecretKey             string        `mapstructure:"jwt_secret_key"`
	JWTRefreshSecretKey      string        `mapstructure:"jwt_refresh_secret_key"`
	JWTExpirationDurationDay int           `mapstructure:"jwt_expiration_duration_day"`
}

type MySQL struct {
	AdditionalParam string `mapstructure:"additional_param"`
	IsEnabled       bool   `mapstructure:"is_enabled"`
	Host            string `mapstructure:"Host"`
	Port            int    `mapstructure:"Port"`
	DbName          string `mapstructure:"db_name"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int64  `mapstructure:"conn_max_lifetime"`
}

type Client struct {
	ID     string `mapstructure:"id"`
	Name   string `mapstructure:"name"`
	Prefix string `mapstructure:"prefix"`
	Secret string `mapstructure:"secret"`
	Scope  string `mapstructure:"scope"`
}

type Error struct {
	Code string `mapstructure:"code"`
	Desc string `mapstructure:"desc"`
}
