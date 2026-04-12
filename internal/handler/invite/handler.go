package invite

import (
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/invite"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	inviteService *invite.Service
	repository    *mysql.Repository
}

func NewHandler(invite *invite.Service, repository *mysql.Repository) *Handler {
	return &Handler{inviteService: invite, repository: repository}
}

func (h *Handler) CreateInvite(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	newInvite, err := h.inviteService.CreateInvite(ctx, sender)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusCreated, newInvite)
}

func (h *Handler) RevokeInvite(ctx fiber.Ctx) error {
	var req requests.RevokeInvite
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	revokedInvite, err := h.inviteService.RevokeInvite(ctx, req.InviteID)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, revokedInvite)
}

func (h *Handler) SearchAllInvites(ctx fiber.Ctx) error {
	invites, err := h.repository.SearchAllInvites(ctx)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, invites)
}
