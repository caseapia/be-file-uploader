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
			return fiber.NewError(fiber.StatusUnauthorized, "ERR_UNAUTHORIZED")
		}
		if !user.HasPermission(permission) {
			return fiber.NewError(fiber.StatusForbidden, "ERR_NO_ACCESS")
		}

		return ctx.Next()
	}
}
