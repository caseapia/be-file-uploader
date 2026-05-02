package user

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
	"github.com/uptrace/bun"
)

func (s *Service) LookupAccount(ctx context.Context, sender *models.User, target int) (account *models.User, err error) {
	account, err = s.repo.LookupUserByID(ctx, target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
		}
		return nil, err
	}

	if sender.ID != account.ID && !sender.HasPermission(role.ManageFiles) {
		account.Storage = slices.DeleteFunc(account.Storage, func(image models.File) bool {
			return image.IsPrivate
		})
	}

	isSensetiveDataAccess := sender.HasPermission(role.ViewPrivateData)
	if !isSensetiveDataAccess {
		account.Private = nil
		account.Geolocation = nil
	} else {
		account.Private = account.GetPrivateData()
		account.Geolocation = account.GetGeolocationData()
	}

	return account, nil
}

func (s *Service) GenerateAPIToken(ctx context.Context, userID int) (token string, err error) {
	token = utils.GenerateSecureToken(32)
	err = s.repo.WithTx(ctx, func(tx bun.Tx) error {
		_, err := s.repo.UpdateUser(ctx, tx, &models.User{
			ID:          userID,
			ShareXToken: &token,
		}, "sharex_token")

		return err
	})

	return token, err
}

func (s *Service) AuthByToken(ctx context.Context, token string) (user *models.User, err error) {
	user, err = s.repo.LookupUserByToken(ctx, token)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "ERR_INVALID_TOKEN")
	}
	return user, nil
}
