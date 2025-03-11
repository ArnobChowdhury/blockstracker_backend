package config

import (
	"fmt"
	"os"
	"time"
)

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

func LoadRedisConfig() (*RedisConfig, error) {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		return nil, fmt.Errorf("REDIS_PASSWORD environment variable is not set")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST environment variable is not set")
	}

	redisAddr := fmt.Sprintf(redisHost + ":6379")

	config := &RedisConfig{
		Addr:         redisAddr,
		Password:     redisPassword,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return config, nil
}
