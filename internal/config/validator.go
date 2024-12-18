package config

import "github.com/go-playground/validator"

func NewValidator() *validator.Validate {
	return validator.New()
}
