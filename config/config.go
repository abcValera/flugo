package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	configName = "api"
	configType = "env"
)

// Contains all configuration variables
// The values are read from api.env file
type Config struct {
	PORT                string        `mapstructure:"PORT"`
	DatabaseDriver      string        `mapstructure:"DATABASE_DRIVER"`
	DatabaseUrl         string        `mapstructure:"DATABASE_URL"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// Override variables from file with the environmet variables
	viper.AutomaticEnv()

	config := Config{}
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
