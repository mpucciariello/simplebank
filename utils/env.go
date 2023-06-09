package utils

import (
	"github.com/spf13/viper"
)

// Config these values are read by viper from the config.env configuration file
type Config struct {
	DriverName    string `mapstructure:"DB_DRIVER"`
	SourceName    string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile("config.env")
	viper.SetConfigType("env")

	// checks if variables exists and loads them into viper
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}