package customvalidators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func IsAValidDate(field validator.FieldLevel) bool {
	var _, err = time.Parse(time.DateOnly, field.Field().String())
	return err == nil
}
