package apiserver

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

var (
	validate            = validator.New()
	ErrorRequiredString = "must not be empty"
)

func init() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})
}

func Validate(u any) map[string]string {
	err := validate.Struct(u)
	if err == nil {
		return nil
	}
	validationError := err.(validator.ValidationErrors)
	apiError := make(map[string]string)
	for _, fieldError := range validationError {
		switch fieldError.Tag() {
		case "required", "required_if":
			apiError[fieldError.Field()] = ErrorRequiredString
		}
	}
	return apiError
}
