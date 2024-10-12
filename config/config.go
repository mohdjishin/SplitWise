package config

import (
	"time"

	log "github.com/mohdjishin/SplitWise/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Port      string `mapstructure:"port"`
	JwtString string `mapstructure:"jwtString"`
	DSN       string `mapstructure:"dsn"`
	LogLevel  string `mapstructure:"logLevel"` // not used as of now kept in .env to change it dynamically from docker env
	ENV       string `mapstructure:"env"`
}

var config Config

func init() {

	viper.SetConfigFile("config.json")
	viper.SetConfigType("json")

	maxRetries := 5
	for retries := 0; retries < maxRetries; retries++ {
		if err := viper.ReadInConfig(); err != nil {
			log.Error("Error reading config file: %v. Retrying in 5 seconds...", zap.Any("error", err))
			time.Sleep(5 * time.Second)
		} else {
			if err := viper.Unmarshal(&config); err != nil {
				log.Panic("Error unmarshalling config: ", zap.Error(err))
			}
			return
		}
	}

	log.Panic("Failed to load config file after multiple attempts.")
}

func GetConfig() Config {
	return config
}
