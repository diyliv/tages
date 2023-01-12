package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GrpcServer GrpcServer
}

type GrpcServer struct {
	Host              string
	Port              string
	Timeout           time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
}

func ReadConfig() *Config {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
