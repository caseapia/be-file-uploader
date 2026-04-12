package user

import (
	"context"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) LookupAccount(ctx context.Context, sender *models.User, target int) (account *models.User, err error) {
	if sender.ID != target || sender.HasPermission(role.ViewOtherProfiles) {
		return nil, fiber.NewError(fiber.StatusForbidden, "ERR_NO_ACCESS")
	}

	account, err = s.repo.LookupUserByID(ctx, target)
	if err != nil {
		return nil, err
	}

	return account, nil
}
