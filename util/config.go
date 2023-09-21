package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbSource     string `mapstructure:"DB_SOURCE"`
	ServerAdress string `mapstructure:"SERVER_ADRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // can be json etc.

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
