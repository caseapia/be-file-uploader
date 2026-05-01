package auth

import (
	"context"
	"strings"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
)

func ValidateSession(r *mysql.Repository, session *models.Session) error {
	if session == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "ERR_SESSION_NOTFOUND")
	}

	if session.IsActive == false {
		slog.Error("Session revoked", "session_id", session.ID)
		return fiber.NewError(fiber.StatusForbidden, "ERR_USER_SESSION_REVOKED")
	}

	now := time.Now()

	if session.ExpiresAt.Before(now) {
		slog.Error("Session expired",
			"session_id", session.ID,
			"expires_at", session.ExpiresAt,
			"current_time", now,
		)

		go func(sid string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, _ = r.TerminateSession(ctx, r.DB, sid)
		}(session.ID)

		return fiber.NewError(fiber.StatusForbidden, "ERR_USER_SESSION_EXPIRED")
	}

	return nil
}

func Middleware(auth *auth.Service, repo *mysql.Repository) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		token := extractToken(ctx)
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "ERR_TOKEN_NOTFOUND")
		}

		ip := ctx.Get("CF-Connecting-IP")
		useragent := ctx.Get("X-User-Agent")
		rayid := ctx.Get("Cf-Ray")
		locale := ctx.Get("X-Locale")

		user, claims, err := auth.ParseJWT(token)
		if err != nil {
			return err
		}

		_, err = repo.UpdateUser(ctx.Context(), repo.DB, &models.User{ID: user.ID, LastIP: ip, Useragent: useragent, CFRayID: rayid, Locale: locale}, "last_ip", "useragent", "cf_ray_id", "locale")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_DATABASE_UPDATE")
		}

		session, err := repo.SearchSessionByID(ctx, claims.SessionID)
		if err != nil || session == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "ERR_SESSION_NOTFOUND")
		}

		if err := ValidateSession(repo, session); err != nil {
			return err
		}

		ctx.Locals("user", user)
		ctx.Locals("session", session)

		return ctx.Next()
	}
}

func extractToken(ctx fiber.Ctx) string {
	header := ctx.Get("Authorization")
	if header != "" {
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1])
		}
		return strings.TrimSpace(strings.TrimPrefix(parts[0], "Bearer"))
	}

	if cookie := ctx.Cookies("access_token"); cookie != "" {
		return strings.TrimSpace(cookie)
	}

	return ""
}
