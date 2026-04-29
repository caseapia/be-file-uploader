package user

import (
	"database/sql"
	"errors"
	"fmt"

	"be-file-uploader/internal/models"
	role2 "be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (s *Service) searchAction(ctx fiber.Ctx, userID, roleID int) (user *models.User, role *models.Role, err error) {
	user, err = s.repo.LookupUserByID(ctx.Context(), userID)
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
	}

	role, err = s.repo.LookupRoleByID(ctx.Context(), roleID)
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusNotFound, "ERR_ROLE_NOTFOUND"+err.Error())
	}

	slog.WithData(slog.M{
		"user_id": userID,
		"role_id": roleID,
		"role":    role,
		"user":    user,
	}).Info("searchAction")

	return user, role, nil
}

func (s *Service) SetUploadLimit(ctx fiber.Ctx, userID int, newUploadLimit int64) (user *models.User, err error) {
	user, err = s.repo.LookupUserByID(ctx.Context(), userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		user.UploadLimit = newUploadLimit
		user, err = s.repo.UpdateUser(ctx, tx, user, "upload_limit")
		if err != nil {
			return err
		}

		s.notify.CreateNotification(ctx.Context(), user.ID, fmt.Sprintf("NOTIFY_UPLOAD_LIMIT_CHANGED+%v", newUploadLimit))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AddUserInRole(ctx fiber.Ctx, sender *models.User, userID, roleID int) (user *models.User, err error) {
	user, role, err := s.searchAction(ctx, userID, roleID)
	if err != nil {
		return nil, err
	}

	if role.IsSystem && !sender.HasPermission(role2.Admin) {
		return nil, fiber.NewError(fiber.StatusForbidden, "ERR_ROLE_ISSUE_FORBIDDEN")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		err = s.repo.AddUserInRole(ctx, tx, userID, roleID)
		s.notify.CreateNotification(ctx.Context(), user.ID, fmt.Sprintf("NOTIFY_ADD_IN_ROLE+%s", role.Name))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	user.Roles = append(user.Roles, *role)

	return user, nil
}

func (s *Service) DeleteUserFromRole(ctx fiber.Ctx, sender *models.User, userID, roleID int) (user *models.User, err error) {
	user, role, err := s.searchAction(ctx, userID, roleID)
	if err != nil {
		return nil, err
	}

	if role.IsSystem && !sender.HasPermission(role2.Admin) {
		return nil, fiber.NewError(fiber.StatusForbidden, "ERR_ROLE_ISSUE_FORBIDDEN")
	}
	if role.HasPermission(role2.Admin) && !sender.HasPermission(role2.Admin) {
		return nil, fiber.NewError(fiber.StatusForbidden, "ERR_ROLE_ISSUE_FORBIDDEN")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		err = s.repo.RemoveUserFromRole(ctx, tx, userID, roleID)
		s.notify.CreateNotification(ctx.Context(), user.ID, fmt.Sprintf("NOTIFY_REMOVE_FROM_ROLE+%s", role.Name))
		if err != nil {
			return err
		}

		return nil
	})

	n := 0
	for _, r := range user.Roles {
		if r.ID != roleID {
			user.Roles[n] = r
			n++
		}
	}
	user.Roles = user.Roles[:n]

	return user, nil
}

func (s *Service) VerifyUser(ctx fiber.Ctx, userID int) (user *models.User, err error) {
	user, err = s.repo.LookupUserByID(ctx.Context(), userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
	}
	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		if !user.IsVerified {
			user.IsVerified = true
			s.notify.CreateNotification(ctx.Context(), user.ID, "NOTIFY_VERIFY_SUCCESS")
		} else {
			user.IsVerified = false
			s.notify.CreateNotification(ctx.Context(), user.ID, "NOTIFY_VERIFY_REMOVED")
		}
		user, err = s.repo.UpdateUser(ctx, tx, user, "is_verified")

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) ResetUserAPIToken(ctx fiber.Ctx, userID int) (user *models.User, err error) {
	user, err = s.repo.LookupUserByID(ctx.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
		}
		return nil, err
	}

	user.ShareXToken = nil

	s.notify.CreateNotification(ctx.Context(), user.ID, "NOTIFY_API_TOKEN_RESET")

	user, err = s.repo.UpdateUser(ctx.Context(), s.repo.DB, user, "sharex_token")

	return user, nil
}
