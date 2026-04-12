package validation

import "github.com/gofiber/fiber/v3"

type ResponseStruct[T any] struct {
	Data T `json:"response"`
}

func Response[T any](c fiber.Ctx, status int, data T) error {
	return c.Status(status).JSON(ResponseStruct[T]{Data: data})
}
