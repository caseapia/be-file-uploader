package user

import (
	"context"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"
)

func (s *Service) LookupAccount(ctx context.Context, sender *models.User, target int) (account *models.User, err error) {
	account, err = s.repo.LookupUserByID(ctx, target)
	if err != nil {
		return nil, err
	}

	isUserManager := sender.HasPermission(role.ManageUsers)
	if isUserManager {
		account.Private = account.GetPrivateData()
	}

	return account, nil
}
