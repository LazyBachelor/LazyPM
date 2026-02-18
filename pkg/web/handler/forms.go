package handler

import (
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

var (
	decoder  = form.NewDecoder()
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func ParseForm[T any](r *http.Request) (*T, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	var data T
	if err := decoder.Decode(&data, r.PostForm); err != nil {
		return nil, err
	}

	return &data, nil
}

func ValidateForm[T any](data *T) error {
	return validate.Struct(data)
}
