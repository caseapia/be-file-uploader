package invite

import (
	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/utils/generate"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func (s *Service) CreateInvite(ctx fiber.Ctx, sender *models.User) (code string, err error) {
	code, err = generate.InviteCode()
	if err != nil {
		return "", err
	}

	invite := &models.Invite{
		Code:      code,
		CreatedBy: sender.ID,
		IsActive:  true,
		UsedBy:    nil,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		err = s.repo.CreateInvite(ctx.Context(), tx, invite)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *Service) RevokeInvite(ctx fiber.Ctx, code int) (invite *models.Invite, err error) {
	invite, err = s.repo.SearchInviteByID(ctx, code)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_INVITE_NOTFOUND")
	}

	invite.IsActive = false

	err = s.repo.WithTx(ctx, func(tx bun.Tx) (err error) {
		err = s.repo.RevokeInvite(ctx, tx, invite)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return invite, nil
}
