package auth

import (
	"fmt"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/utils/generate"
	"be-file-uploader/pkg/utils/token"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/gookit/slog"
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

	ip := ctx.Get("CF-Connecting-IP")
	useragent := ctx.Get("X-User-Agent")
	code, country, city := s.geo.GetGeoString(ip)

	refreshHash := generate.HashToken(refreshToken)
	session := &models.Session{
		ID:          sessionID,
		UserID:      user.ID,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		IsActive:    true,
		IPAddress:   ip,
		UserAgent:   useragent,
		RefreshHash: refreshHash,
		GeoString:   fmt.Sprintf("%s, %s", country, city),
		Geolocation: models.Geolocation{
			Code:    code,
			City:    city,
			Country: country,
		},
	}

	if err = s.repo.CreateSession(ctx.Context(), s.repo.DB, session); err != nil {
		slog.WithData(slog.M{"error": err}).Error("Session creation failed")
		return nil, "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_SESSION_CREATION")
	}

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
	ip := ctx.IP()
	useragent := ctx.Get("X-User-Agent")

	session, err := s.repo.SearchSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "ERR_TOKEN_INVALID")
	}

	if session.UserAgent != useragent {
		slog.Warn("Session hijacking attempt? UA changed", "sid", session.ID, "old", session.UserAgent, "new", useragent)
	}

	newRefresh, _ := GenerateRefreshToken()
	newAccess, _ := token.GenerateAccessToken(session.UserID, 1, session.ID)
	code, country, city := s.geo.GetGeoString(ip)

	session.RefreshHash = generate.HashToken(newRefresh)
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	session.IPAddress = ip
	session.UserAgent = useragent
	session.GeoString = fmt.Sprintf("%s, %s", country, city)
	session.Geolocation = models.Geolocation{
		Code:    code,
		City:    city,
		Country: country,
	}

	if _, err := s.repo.UpdateSession(ctx, s.repo.DB, *session); err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
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

	_, err = s.repo.UpdateSession(ctx, s.repo.DB, *session)
	if err != nil {
		slog.WithData(slog.M{
			"error":   err,
			"user":    session.UserID,
			"session": session,
		}).Error("Session update failed")
		return fiber.NewError(fiber.StatusInternalServerError, "ERR_SESSION_UPDATE")
	}

	return nil
}
