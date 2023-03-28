package apiserver

import (
	"Social-app/internal/model"
	"github.com/go-playground/validator/v10"
	"reflect"
)

var (
	validate        = validator.New()
	ErrorFieldEmpty = "must not be empty"
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
			apiError[fieldError.Field()] = ErrorFieldEmpty
		}
	}
	return apiError
}

func (s *Server) ValidateUserCreateUpdate(username, email string) (map[string]string, error) {
	vucu := make(map[string]string)
	if username != "" {
		exist, err := s.store.User().ExistByUsername(username)
		if err != nil {
			return vucu, err
		} else if exist {
			vucu["username"] = "username is already exist"
		}
	}
	if email != "" {
		exist, err := s.store.User().ExistByEmail(email)
		if err != nil {
			return vucu, err
		} else if exist {
			vucu["email"] = "email is already exist"
		}
	}
	return vucu, nil
}

func (s *Server) ValidateUserLogin(req *UserLogin) (*model.User, error) {
	var (
		u   *model.User
		err error
	)
	if req.Username != "" {
		u = &model.User{Username: req.Username}
		err = s.store.User().GetPasswordByUsername(u)
	} else {
		u = &model.User{Email: req.Email}
		err = s.store.User().GetPasswordByEmail(u)
	}
	return u, err
}
