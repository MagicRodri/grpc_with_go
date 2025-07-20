package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Kafka KafkaConfig `mapstructure:"kafka"`
		GRPC  GRPC        `mapstructure:"grpc"`
	}

	KafkaConfig struct {
		Brokers     []string `mapstructure:"brokers"`
		Credentials struct {
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"credentials"`
	}

	GRPC struct {
		Address string `mapstructure:"address"`
	}
)

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}
