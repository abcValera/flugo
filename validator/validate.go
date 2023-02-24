package validator

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type CustomValidator interface {
	Validate(s interface{}) error
}

type FlugoValidator struct {
	validate *validator.Validate
}

func (fv *FlugoValidator) Validate(s interface{}) error {
	if err := fv.validate.Struct(s); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return nil
}
