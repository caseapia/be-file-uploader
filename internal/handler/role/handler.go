package role

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/role"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	roleService *role.Service
	repo        *mysql.Repository
}

func NewHandler(role *role.Service, repo *mysql.Repository) *Handler {
	return &Handler{roleService: role, repo: repo}
}

func (h *Handler) LookupAllRoles(ctx fiber.Ctx) error {
	roles, err := h.repo.LookupAllRoles(ctx.Context())
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, roles)
}
