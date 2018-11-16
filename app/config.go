package app

import (
	"fmt"

	"github.com/tomochain/backend-matching-engine/utils"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/viper"
)

// Config stores the application-wide configurations
var Config appConfig
var logger = utils.Logger

type appConfig struct {
	// the path to the error message file. Defaults to "config/errors.yaml"
	ErrorFile string `mapstructure:"error_file"`
	// the server port. Defaults to 8080
	ServerPort int `mapstructure:"server_port"`
	// the data source name (DSN) for connecting to the database. required.
	DSN string `mapstructure:"dsn"`
	// the data source name (DSN) for connecting to the database. required.
	DBName string `mapstructure:"db_name"`
	// the make fee is the percentage to charged from maker
	MakeFee float64 `mapstructure:"make_fee"`
	// the take fee is the percentage to charged from maker
	TakeFee float64 `mapstructure:"take_fee"`
	// the Rabbitmq is the URI of rabbitmq to use
	Rabbitmq string `mapstructure:"rabbitmq"`
	// the redis is the URI of redis to use
	Redis string `mapstructure:"redis"`
	// the signing method for JWT. Defaults to "HS256"
	JWTSigningMethod string `mapstructure:"jwt_signing_method"`
	// JWT signing key. required.
	JWTSigningKey string `mapstructure:"jwt_signing_key"`
	// JWT verification key. required.
	JWTVerificationKey string `mapstructure:"jwt_verification_key"`
	// TickDuration is user by tick streaming cron
	TickDuration map[string][]int64 `mapstructure:"tick_duration"`

	Logs map[string]string `mapstructure:"logs"`

	Ethereum map[string]string `mapstructure:"ethereum"`
}

func (config appConfig) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.DSN, validation.Required),
	)
}

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
// Environment variables with the prefix "RESTFUL_" in their names are also read automatically.
func LoadConfig(configPath string, env string) error {
	v := viper.New()

	if env != "" {
		v.SetConfigName("config." + env)
	}

	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}

	v.SetEnvPrefix("tomo")
	v.AutomaticEnv()

	err = v.Unmarshal(&Config)
	if err != nil {
		return err
	}

	// log information

	logger.Infof("Server port: %v", Config.ServerPort)
	logger.Infof("Ethereum node HTTP url: %v", Config.Ethereum["http_url"])
	logger.Infof("Ethereum node WS url: %v", Config.Ethereum["ws_url"])
	logger.Infof("MongoDB url: %v", Config.DSN)
	logger.Infof("Redis url: %v", Config.Redis)
	logger.Infof("RabbitMQ url: %v", Config.Rabbitmq)
	logger.Infof("Exchange contract address: %v", Config.Ethereum["exchange_address"])
	logger.Infof("Fee Account: %v", Config.Ethereum["fee_account"])

	return Config.Validate()
}
