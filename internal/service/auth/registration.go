package auth

import (
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/utils/generate"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (s *Service) validateRegistration(ctx fiber.Ctx, username, inviteCode string) (*models.Invite, error) {
	user, err := s.repo.LookupUserByName(ctx, username)
	if user != nil && err == nil {
		return nil, fiber.NewError(fiber.StatusConflict, "ERR_USER_ALREADY_EXISTS")
	}

	invite, err := s.repo.SearchInviteByCode(ctx, inviteCode)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "ERR_INVITE_NOT_FOUND")
	}

	if !invite.IsActive {
		return nil, fiber.NewError(fiber.StatusConflict, "ERR_INVITE_ALREADY_USED")
	}

	return invite, nil
}

func (s *Service) createUserWithInvite(ctx fiber.Ctx, username, password string, invite *models.Invite) (*models.User, error) {
	ip := ctx.IP()
	useragent := ctx.Get("X-User-Agent")

	hashPassword, err := generate.HashPassword(password)
	if err != nil {
		slog.WithData(slog.M{"error": err, "ip": ip}).Error("Password hashing failed")
		return nil, fiber.NewError(fiber.StatusConflict, "ERR_USER_REGISTER_HASHCREATION")
	}

	user := &models.User{
		Username:    username,
		Password:    hashPassword,
		RegisterIP:  ip,
		Useragent:   useragent,
		InviteID:    invite.ID,
		CreatedAt:   time.Now(),
		DiscordUID:  nil,
		DiscordName: nil,
	}

	err = s.repo.WithTx(ctx, func(tx bun.Tx) error {
		return s.registerTx(ctx, tx, user, invite)
	})
	if err != nil {
		return nil, err
	}

	user, err = s.repo.LookupUserByName(ctx, username)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_USER_LOOKUP_AFTER_REGISTER")
	}

	return user, nil
}

func (s *Service) registerTx(ctx fiber.Ctx, tx bun.Tx, user *models.User, invite *models.Invite) error {
	created, err := s.repo.CreateUser(ctx.Context(), tx, *user)
	if err != nil {
		slog.WithData(slog.M{"error": err, "user": user}).Error("Failed to create user")
		return fiber.NewError(fiber.StatusConflict, "ERR_USER_UNKNOWN_CREATION_ERROR")
	}

	invite.UsedBy = &created.ID
	if err := s.repo.UseInvite(ctx.Context(), tx, *invite); err != nil {
		slog.WithData(slog.M{"error": err, "invite_id": invite.ID}).Error("Failed to mark invite as used")
		return fiber.NewError(fiber.StatusConflict, "ERR_INVITE_MARK_ERR")
	}

	if err := s.repo.AddUserInRole(ctx, tx, created.ID, 1); err != nil {
		slog.WithData(slog.M{"error": err, "user": created}).Error("Failed to add user in default role")
		return fiber.NewError(fiber.StatusConflict, "ERR_USER_ADDINROLE_FAILED")
	}

	*user = *created
	return nil
}

func (s *Service) createSessionTokens(ctx fiber.Ctx, user *models.User) (access, refresh string, err error) {
	go func(u models.User) {
		if err := s.repo.CleanupExpiredSessions(ctx, s.repo.DB, &u); err != nil {
			slog.WithData(slog.M{"error": err}).Error("Cleanup expired sessions failed")
		}
	}(*user)

	user, access, refresh, err = s.Login(ctx, user.Username, user.Password)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Register(ctx fiber.Ctx, username, password, inviteCode string) (*models.User, string, string, error) {
	invite, err := s.validateRegistration(ctx, username, inviteCode)
	if err != nil {
		return nil, "", "", err
	}

	user, err := s.createUserWithInvite(ctx, username, password, invite)
	if err != nil {
		return nil, "", "", err
	}

	access, refresh, err := s.createSessionTokens(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	return user, access, refresh, nil
}
