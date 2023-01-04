package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port    string `mapstructure:"PORT"`
	LogFile string `mapstructure:"LOG_PATH"`

	HeaderKey  string `mapstructure:"HEADER_KEY"`
	PaylaodKey string `mapstructure:"PAYLOAD_KEY"`

	DBconn string `mapstructure:"DB_SOURCE"`
}

func LoadServerConfig(path string) (config ServerConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("cfg")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
