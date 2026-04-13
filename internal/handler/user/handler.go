package user

import (
	"strconv"

	"be-file-uploader/internal/models"
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/user"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	userService *user.Service
	repo        *mysql.Repository
}

func NewHandler(user *user.Service, repo *mysql.Repository) *Handler {
	return &Handler{userService: user, repo: repo}
}

func (h *Handler) LookupMyAccount(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	acc, err := h.userService.LookupAccount(ctx, sender, sender.ID)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, acc)
}

func (h *Handler) LookupProfile(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)
	sender := account.GetUserFromContext(ctx)

	acc, err := h.userService.LookupAccount(ctx, sender, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, acc)
}

func (h *Handler) SetUploadLimit(ctx fiber.Ctx) error {
	var req requests.SetUploadLimitRequest
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	target, err := h.userService.SetUploadLimit(ctx, req.User, req.Limit)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, target)
}

func (h *Handler) PopulateUserList(ctx fiber.Ctx) error {
	users, err := h.repo.LookupUsers(ctx, 30)
	if err != nil {
		return err
	}

	if users == nil {
		users = make([]models.User, 0)
	}

	return validation.Response(ctx, 200, users)
}

func (h *Handler) AddUserInRole(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var req requests.AddUserInRole
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	u, err := h.userService.AddUserInRole(ctx, sender, req.User, req.Role)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, u)
}

func (h *Handler) RemoveUserFromRole(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var req requests.RemoveUserFromRole
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	u, err := h.userService.DeleteUserFromRole(ctx, sender, req.User, req.Role)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, u)
}
