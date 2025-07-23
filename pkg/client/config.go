package client

import "github.com/MagicRodri/grpc_with_go/pkg/validation"

type StatusServiceConfig struct {
	Host    string `mapstructure:"host" validate:"required"`
	Name    string `mapstructure:"name" validate:"required"`
	Timeout int    `mapstructure:"timeout" validate:"required"`
}

func (cfg *StatusServiceConfig) Validate() error {
	return validation.Validate(cfg)
}
