package utils

import "github.com/spf13/viper"

// Config these values are read by viper from the app.env configuration file
type Config struct {
	DriverName    string `mapstructure:"DRIVER_NAME"`
	SourceName    string `mapstructure:"SOURCE_NAME"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile("app")
	viper.SetConfigType("env")

	// checks if variables exists and loads them into viper
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
