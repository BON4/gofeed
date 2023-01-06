package config

import (
	"fmt"
	"os"
	"strings"
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

	DBconn         string `mapstructure:"DB_SOURCE"`
	TestDBconn     string `mapstructure:"TEST_DB_SOURCE"`
	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`
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

	for _, k := range viper.AllKeys() {
		k = strings.ToUpper(k)
		fmt.Printf("Setting %s\n", k)
		if err = os.Setenv(k, fmt.Sprintf("%s", viper.Get(k))); err != nil {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
