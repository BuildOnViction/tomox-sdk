package app

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/viper"
)

// Config stores the application-wide configurations
var Config appConfig

type appConfig struct {
	// the path to the error message file. Defaults to "config/errors.yaml"
	ErrorFile string `mapstructure:"error_file"`
	LogLevel  string `mapstructure:"log_level"`
	// the server port. Defaults to 8080
	ServerPort int `mapstructure:"server_port"`
	// the data source name (MongoURL) for connecting to the database. required.
	MongoURL        string `mapstructure:"mongo_url"`
	MongoDBPassword string `mapstructure:"mongo_password"`
	MongoDBUsername string `mapstructure:"mongo_username"`

	// the data source name (DSN) for connecting to the database. required.
	DBName string `mapstructure:"db_name"`

	ApiAuthKey string `mapstructure:"api_auth_key"`

	// the RabbitMQURL is the URI of rabbitmq to use
	RabbitMQURL string `mapstructure:"rabbitmq_url"`

	// TickDuration is user by tick streaming cron
	TickDuration map[string][]int64 `mapstructure:"tick_duration"`

	Tomochain map[string]string `mapstructure:"tomochain"`

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

	return Config.Validate()
}
