package role

import (
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/role"
	"be-file-uploader/pkg/utils/account"
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

func (h *Handler) CreateRole(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var req requests.CreateRole
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	createdRole, err := h.roleService.CreateRole(ctx, req.Name, req.Color, req.Permissions, req.IsSystem, *sender)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 201, createdRole)
}

func (h *Handler) UpdateRole(ctx fiber.Ctx) error {
	var req requests.UpdateRole
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	updatedRole, err := h.roleService.EditRole(ctx, req.Color, req.Name, req.Permissions, req.IsSystem, req.RoleID)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 201, updatedRole)
}

func (h *Handler) DeleteRole(ctx fiber.Ctx) error {
	var req requests.DeleteRole
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	if err := h.roleService.DeleteRole(ctx, req.ID); err != nil {
		return err
	}

	return validation.Response(ctx, 200, "OK")
}
