package auth

import (
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/utils/generate"
	"be-file-uploader/pkg/utils/token"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (s *Service) Login(ctx fiber.Ctx, username, password string) (user *models.User, access string, refresh string, err error) {
	user, err = s.repo.LookupUserByName(ctx, username)
	if err != nil {
		return nil, "", "", fiber.NewError(fiber.StatusNotFound, "ERR_WRONG_CREDENTIALS")
	}
	if !generate.CheckPassword(user.Password, password) {
		return nil, "", "", fiber.NewError(fiber.StatusNotFound, "ERR_WRONG_CREDENTIALS")
	}

	sessionID := uuid.NewString()
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("Refresh token generation failed")
		return nil, "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_TOKEN_GENERATION")
	}

	ip := ctx.IP()
	useragent := ctx.Get("X-User-Agent")

	refreshHash := generate.HashToken(refreshToken)
	session := &models.Session{
		ID:           sessionID,
		UserID:       user.ID,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     true,
		IPAddress:    ip,
		UserAgent:    useragent,
		LastActiveAt: time.Now(),
		RefreshHash:  refreshHash,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		if err = s.repo.CreateSession(ctx.Context(), tx, session); err != nil {
			slog.WithData(slog.M{"error": err}).Error("Session creation failed inside tx")
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_SESSION_CREATION")
		}
		return nil
	})
	if err != nil {
		return nil, "", "", err
	}

	go func(user models.User) {
		if err := s.repo.CleanupExpiredSessions(ctx, s.repo.DB, &user); err != nil {
			slog.WithData(slog.M{
				"error": err,
			}).Error("Cleanup expired sessions failed")
		}
	}(*user)

	access, err = token.GenerateAccessToken(user.ID, 1, sessionID)
	if err != nil {
		slog.WithData(slog.M{
			"error":     err,
			"user":      user,
			"tokenVer":  1,
			"sessionID": sessionID,
		}).Error("Access token generation failed")
		return nil, "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_TOKEN_GENERATION")
	}

	return user, access, refreshToken, nil
}

func (s *Service) RefreshToken(ctx fiber.Ctx, refreshToken string) (access string, refresh string, err error) {
	refreshHash := generate.HashToken(refreshToken)

	session, err := s.repo.SearchSessionByRefreshHash(ctx, refreshHash)
	if err != nil || session.IsActive == false || session.ExpiresAt.Before(time.Now()) {
		slog.WithData(slog.M{
			"error":        err,
			"refreshToken": refreshToken,
		}).Error("Refresh token lookup failed")
		return "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_TOKEN_GENERATION")
	}

	user, err := s.repo.LookupUserByID(ctx, session.UserID)
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
			"user":  user,
		}).Error("User lookup failed")
		return "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_USER_NOTFOUND")
	}

	ip := ctx.IP()
	useragent := ctx.Get("X-User-Agent")

	err = s.repo.WithTx(ctx, func(tx bun.Tx) (err error) {
		refresh, err = GenerateRefreshToken()
		if err != nil {
			slog.WithData(slog.M{
				"error": err,
				"user":  user,
			}).Error("Refresh token generation failed")
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_TOKEN_GENERATION")
		}

		access, err = token.GenerateAccessToken(user.ID, 1, session.ID)
		if err != nil {
			slog.WithData(slog.M{
				"error":   err,
				"user":    user,
				"refresh": refresh,
			}).Error("Access token generation failed")
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_TOKEN_GENERATION")
		}

		session.RefreshHash = generate.HashToken(refresh)
		session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
		session.IPAddress = ip
		session.UserAgent = useragent

		if err := s.repo.UpdateSession(ctx, tx, session); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Logout(ctx fiber.Ctx, session *models.Session, user *models.User) error {
	user, err := s.repo.LookupUserByID(ctx, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
	}

	if session.UserID != user.ID {
		return fiber.NewError(fiber.StatusNotFound, "ERR_SESSION_NOTFOUND")
	}
	if !session.IsActive {
		return fiber.NewError(fiber.StatusNotFound, "ERR_SESSION_NOTACTIVE")
	}

	session.IsActive = false
	session.IPAddress = ctx.IP()
	session.UserAgent = ctx.Get("X-User-Agent")
	session.ExpiresAt = time.Now()

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		err = s.repo.UpdateSession(ctx, tx, session)
		if err != nil {
			slog.WithData(slog.M{
				"error":   err,
				"user":    session.UserID,
				"session": session,
			}).Error("Session update failed")
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_SESSION_UPDATE")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
