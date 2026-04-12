package user

import (
	"strconv"

	"be-file-uploader/internal/service/user"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	userService *user.Service
}

func NewHandler(user *user.Service) *Handler {
	return &Handler{userService: user}
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
