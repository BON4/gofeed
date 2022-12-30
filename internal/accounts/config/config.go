package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port    string `mapstructure:"PORT"`
	LogFile string `mapstructure:"LOG_PATH"`

	AcessDuration   time.Duration `mapstructure:"ACESS_DUR"`
	RefreshDuration time.Duration `mapstructure:"REFRESH_DUR"`
	SecretToken     string        `mapstructure:"SECRET_TOKEN"`

	HeaderKey  string `mapstructure:"HEADER_KEY"`
	PaylaodKey string `mapstructure:"PAYLOAD_KEY"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

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
