package validator

import (
	v "github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type CustomValidator interface {
	Validate(s interface{}) error
}

type FlugoValidator struct {
	Validator *v.Validate
}

func (fv *FlugoValidator) Validate(s interface{}) error {
	if err := fv.Validator.Struct(s); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return nil
}
