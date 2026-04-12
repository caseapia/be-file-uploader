package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}
	return nil
}

func ParseAndValidate(c fiber.Ctx, s interface{}) error {
	if err := c.Bind().Body(s); err != nil {
		return err
	}
	return ValidateStruct(s)
}
