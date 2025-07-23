package logger

import "github.com/MagicRodri/grpc_with_go/pkg/validation"

type Config struct {
	Level  string `mapstructure:"level"`
	Path   string `mapstructure:"path" validate:"omitempty,filepath"`
	Format string `mapstructure:"format" validate:"oneof=json text"`
	Output string `mapstructure:"output" validate:"oneof=stdout file"`
}

func (cfg *Config) Validate() error {
	return validation.Validate(cfg)
}
