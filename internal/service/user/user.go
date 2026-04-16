package user

import (
	"context"
	"slices"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"
)

func (s *Service) LookupAccount(ctx context.Context, sender *models.User, target int) (account *models.User, err error) {
	account, err = s.repo.LookupUserByID(ctx, target)
	if err != nil {
		return nil, err
	}

	if sender.ID != account.ID && !sender.HasPermission(role.ManageFiles) {
		account.Storage = slices.DeleteFunc(account.Storage, func(image models.Image) bool {
			return image.IsPrivate
		})
	}

	isUserManager := sender.HasPermission(role.ManageUsers)
	if isUserManager {
		account.Private = account.GetPrivateData()
	}

	return account, nil
}
