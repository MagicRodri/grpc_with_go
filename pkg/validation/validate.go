package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validatable interface {
	Validate() error
}

func Validate(cfg Validatable) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cfg)
	if err != nil {
		return fmt.Errorf("error validate config: %w", err)
	}
	return nil
}
