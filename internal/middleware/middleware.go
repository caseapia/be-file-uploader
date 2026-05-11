package middleware

import (
	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"
	userEnum "be-file-uploader/pkg/enums/user"

	"github.com/gofiber/fiber/v3"
)

func RequirePermission(permission role.Permission) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*models.User)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "ERR_UNAUTHORIZED")
		}
		if !user.HasPermission(permission) {
			return fiber.NewError(fiber.StatusForbidden, "ERR_NO_ACCESS")
		}

		return ctx.Next()
	}
}

type restrictionChecker func(user *models.User) error

var flagCheckers = map[string]restrictionChecker{
	"not_restricted_upload": func(user *models.User) error {
		if user.ActiveRestriction != nil && user.ActiveRestriction.Type == userEnum.BanTypeUpload {
			return fiber.NewError(fiber.StatusForbidden, "ERR_USER_RESTRICTED")
		}
		return nil
	},
	"not_restricted_like": func(user *models.User) error {
		if user.ActiveRestriction != nil && user.ActiveRestriction.Type == userEnum.BanTypeLike {
			return fiber.NewError(fiber.StatusForbidden, "ERR_USER_RESTRICTED")
		}
		return nil
	},
	"not_restricted_comment": func(user *models.User) error {
		if user.ActiveRestriction != nil && user.ActiveRestriction.Type == userEnum.BanTypeComment {
			return fiber.NewError(fiber.StatusForbidden, "ERR_USER_RESTRICTED")
		}
		return nil
	},
	"not_restricted_ban": func(user *models.User) error {
		if user.ActiveRestriction != nil {
			return fiber.NewError(fiber.StatusForbidden, "ERR_USER_RESTRICTED")
		}
		return nil
	},
	"sharex_token": func(user *models.User) error {
		if user.ShareXToken == nil {
			return fiber.NewError(fiber.StatusForbidden, "ERR_TOKEN_NOTFOUND")
		}
		return nil
	},
}

func Require(flags []string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*models.User)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "ERR_UNAUTHORIZED")
		}

		for _, flag := range flags {
			checker, exists := flagCheckers[flag]
			if !exists {
				continue
			}
			if err := checker(user); err != nil {
				return err
			}
		}

		return ctx.Next()
	}
}
