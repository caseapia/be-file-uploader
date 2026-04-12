package developer

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	repository *mysql.Repository
}

func NewHandler(repository *mysql.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Ping(ctx fiber.Ctx) error {
	return validation.Response(ctx, 200, "pong")
}
