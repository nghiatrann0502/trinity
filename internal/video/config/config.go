package config

import (
	"fmt"
	"strings"

	"github.com/nghiatrann0502/trinity/pkg/config"
	"github.com/spf13/viper"
)

type (
	Config struct {
		config.App   `mapstructure:"app"`
		config.DB    `mapstructure:"db"`
		config.Redis `mapstructure:"redis"`
		GRPC         `mapstructure:"grpc"`
	}

	GRPC struct {
		Port int `mapstructure:"port"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set the file name and path
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Enable environment variable overrides
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the nested configuration
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
