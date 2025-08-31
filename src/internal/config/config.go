// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	GitHub    GitHubConfig    `mapstructure:"github"`
	Detection DetectionConfig `mapstructure:"detection"`
}

type ServerConfig struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type GitHubConfig struct {
	WebhookSecret string `mapstructure:"webhook_secret"`
	Token         string `mapstructure:"token"`
	AppID         int64  `mapstructure:"app_id"`
	PrivateKey    string `mapstructure:"private_key"`
}

type DetectionConfig struct {
	DefaultThreshold float64 `mapstructure:"default_threshold"`
	MinSamples       int     `mapstructure:"min_samples"`
	MaxSamples       int     `mapstructure:"max_samples"`
}

func Load() (*Config, error) {
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("database.path", "./regression.db")
	viper.SetDefault("detection.default_threshold", 10.0)
	viper.SetDefault("detection.min_samples", 5)
	viper.SetDefault("detection.max_samples", 50)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}