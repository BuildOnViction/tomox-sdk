package app

import (
    "fmt"

    "github.com/go-ozzo/ozzo-validation"
    "github.com/spf13/viper"
    "github.com/tomochain/tomox-sdk/utils"
)

// Config stores the application-wide configurations
var Config appConfig

var logger = utils.Logger

// var logger = utils.NoopLogger

type appConfig struct {
    // the path to the error message file. Defaults to "config/errors.yaml"
    ErrorFile string `mapstructure:"error_file"`
    // the server port. Defaults to 8080
    ServerPort int `mapstructure:"server_port"`
    // the data source name (MongoURL) for connecting to the database. required.
    MongoURL        string `mapstructure:"mongo_url"`
    MongoDBPassword string `mapstructure:"mongo_password"`
    MongoDBUsername string `mapstructure:"mongo_username"`

    // simulate the environment
    Simulated bool `mapstructure:"simulated"`

    // the data source name (DSN) for connecting to the database. required.
    DBName string `mapstructure:"db_name"`

    // the RabbitMQURL is the URI of rabbitmq to use
    RabbitMQURL string `mapstructure:"rabbitmq_url"`

    // the signing method for JWT. Defaults to "HS256"
    JWTSigningMethod string `mapstructure:"jwt_signing_method"`
    // JWT signing key. required.
    JWTSigningKey string `mapstructure:"jwt_signing_key"`
    // JWT verification key. required.
    JWTVerificationKey string `mapstructure:"jwt_verification_key"`
    // TickDuration is user by tick streaming cron
    TickDuration map[string][]int64 `mapstructure:"tick_duration"`

    Logs map[string]string `mapstructure:"logs"`

    Tomochain map[string]string `mapstructure:"tomochain"`

    CoingeckoAPIUrl string `mapstructure:"coingecko_api_url"`

    SupportedCurrencies string `mapstructure:"supported_currencies"`

    Env string `mapstructure:"env"`
}

func (config appConfig) Validate() error {
    return validation.ValidateStruct(&config,
        validation.Field(&config.MongoURL, validation.Required),
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
    // // update config, if yaml does not presented this config, we can still apply from env
    // Config.Simulated = v.GetBool("simulated")

    // log information
    logger.Infof("Server port: %v", Config.ServerPort)
    logger.Infof("Tomochain node HTTP url: %v", Config.Tomochain["http_url"])
    logger.Infof("Tomochain node WS url: %v", Config.Tomochain["ws_url"])
    logger.Infof("MongoDB url: %v", Config.MongoURL)
    logger.Infof("RabbitMQ url: %v", Config.RabbitMQURL)
    logger.Infof("Exchange contract address: %v", Config.Tomochain["exchange_address"])
    logger.Infof("Env: %v", Config.Env)
    return Config.Validate()
}
