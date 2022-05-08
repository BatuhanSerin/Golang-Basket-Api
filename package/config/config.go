package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig   ServerConfig
	JWTConfig      JWTConfig
	DatabaseConfig DatabaseConfig
	Logger         Logger
}

type ServerConfig struct {
	AppVersion       string
	Mode             string
	RoutePrefix      string
	Debug            bool
	Port             string
	TimeoutSecs      int64
	ReadTimeoutSecs  int64
	WriteTimeoutSecs int64
}

type JWTConfig struct {
	SessionTime int
	SecretKey   string
}

type DatabaseConfig struct {
	DataSourceName  string
	Name            string
	MigrationFolder string
	MaxOpen         int
	MaxIdle         int
	MaxLifetime     int
}

type Logger struct {
	Development bool
	Encodings   string
	Level       string
}

func LoadConfig(configFileName string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var config Config
	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil

}
func GetSecretKey() string {
	var config Config
	return config.JWTConfig.SecretKey
}
