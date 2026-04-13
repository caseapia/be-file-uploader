package role

import (
	"database/sql"
	"errors"
	"time"

	"be-file-uploader/internal/models"
	roleEnum "be-file-uploader/pkg/enums/role"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func (s *Service) CreateRole(ctx fiber.Ctx, roleName, color string, permissions []roleEnum.Permission, isSystem bool, sender models.User) (role *models.Role, err error) {
	role = &models.Role{
		Name:        roleName,
		Permissions: permissions,
		IsSystem:    isSystem,
		CreatedBy:   sender.ID,
		CreatedAt:   time.Now(),
		Color:       color,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		role, err = s.repo.CreateRole(ctx.Context(), tx, *role)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return fiber.NewError(fiber.StatusConflict, "ERR_ROLE_ALREADY_EXISTS")
			}

			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *Service) DeleteRole(ctx fiber.Ctx, roleID int) (err error) {
	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		err = s.repo.DeleteRole(ctx.Context(), tx, roleID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusNotFound, "ERR_ROLE_NOT_FOUND")
			}
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) EditRole(ctx fiber.Ctx, color, name string, permissions []roleEnum.Permission, isSystem bool, id int) (role *models.Role, err error) {
	role, err = s.repo.LookupRoleByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "ERR_ROLE_NOT_FOUND")
		}
	}

	role = &models.Role{
		ID:          id,
		Name:        name,
		Permissions: permissions,
		IsSystem:    isSystem,
		Color:       color,
		CreatedAt:   role.CreatedAt,
		CreatedBy:   role.CreatedBy,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		err = s.repo.EditRole(ctx.Context(), tx, *role)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusNotFound, "ERR_ROLE_NOT_FOUND")
			}
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return fiber.NewError(fiber.StatusConflict, "ERR_ROLE_ALREADY_EXISTS")
			}

			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}
