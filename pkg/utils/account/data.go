package account

import (
	"be-file-uploader/internal/models"

	"github.com/gofiber/fiber/v3"
)

func GetUserFromContext(c fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}

func GetSessionFromContext(c fiber.Ctx) *models.Session {
	return c.Locals("session").(*models.Session)
}
