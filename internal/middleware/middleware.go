package middleware

import (
	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func RequirePermission(permission role.Permission) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*models.User)
		if !ok {
			return fiber.ErrUnauthorized
		}

		if !user.HasPermission(permission) {
			return fiber.ErrForbidden
		}

		return ctx.Next()
	}
}
