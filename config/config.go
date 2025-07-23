package config

import (
	"fmt"

	"github.com/MagicRodri/grpc_with_go/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type GrpcConfig struct {
	Address string `mapstructure:"address" validate:"required"`
}

type Config struct {
	GRPC   GrpcConfig    `mapstructure:"grpc" validate:"required"`
	Logger logger.Config `mapstructure:"logger" validate:"required"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("failed to validate config file %s: %w", path, err)
	}

	return &config, nil
}
