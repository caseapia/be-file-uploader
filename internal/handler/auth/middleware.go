package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/auth"
	userRelation "be-file-uploader/pkg/enums/user"
	"be-file-uploader/pkg/geo"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
)

func ValidateSession(r *mysql.Repository, session *models.Session) error {
	if session == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "ERR_SESSION_NOTFOUND")
	}

	if !session.IsActive {
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

		deletionThreshold := session.ExpiresAt.Add(2 * 24 * time.Hour)

		if now.After(deletionThreshold) {
			go func(sid string) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, _ = r.TerminateSession(ctx, r.DB, sid)
				slog.Info("Session permanently deleted after grace period", "session_id", sid)
			}(session.ID)
		} else {
			slog.Info("Session expired but kept in DB (grace period)", "session_id", session.ID)
		}

		return fiber.NewError(fiber.StatusForbidden, "ERR_USER_SESSION_EXPIRED")
	}

	return nil
}

func Middleware(auth *auth.Service, geo *geo.Service, repo *mysql.Repository) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		token := extractToken(ctx)
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "ERR_TOKEN_NOTFOUND")
		}

		user, claims, err := auth.ParseJWT(token)
		if err != nil {
			return err
		}

		enrichUserMeta(ctx, user, geo)

		if err := processExpiredBan(ctx.Context(), user, repo); err != nil {
			return err
		}

		_, err = repo.UpdateUser(ctx.Context(), repo.DB, user, "last_ip", "useragent", "cf_ray_id", "locale", "last_seen", "geo_string")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_DATABASE_UPDATE")
		}

		session, err := getAndValidateSession(ctx.Context(), repo, claims.SessionID)
		if err != nil {
			return err
		}

		ctx.Locals("user", user)
		ctx.Locals("session", session)

		return ctx.Next()
	}
}

func enrichUserMeta(ctx fiber.Ctx, user *models.User, geo *geo.Service) {
	ip := ctx.Get("CF-Connecting-IP")
	code, country, city := geo.GetGeoString(ip)

	user.LastIP = ip
	user.Useragent = ctx.Get("X-User-Agent")
	user.CFRayID = ctx.Get("Cf-Ray")
	user.Locale = ctx.Get("X-Locale")
	user.LastSeen = time.Now()
	user.GeoString = fmt.Sprintf("%s, %s", country, city)
	user.Geolocation = &models.Geolocation{
		CountryCode: code,
		City:        city,
		Country:     country,
	}
}

func processExpiredBan(ctx context.Context, user *models.User, repo *mysql.Repository) error {
	if user.ActiveRestriction == nil || user.ActiveRestriction.UnbanAt == nil {
		return nil
	}

	if time.Now().Before(*user.ActiveRestriction.UnbanAt) {
		return nil
	}

	expiredBanID := *user.ActiveRestrictionID

	err := repo.RemoveBan(ctx, repo.DB, expiredBanID, userRelation.BanStatusExpired, nil)
	if err != nil {
		return err
	}

	user.ActiveRestrictionID = nil
	user.ActiveRestriction = nil

	return nil
}

func getAndValidateSession(ctx context.Context, repo *mysql.Repository, sessionID string) (*models.Session, error) {
	session, err := repo.SearchSessionByID(ctx, sessionID)
	if err != nil || session == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "ERR_SESSION_NOTFOUND")
	}

	if err := ValidateSession(repo, session); err != nil {
		return nil, err
	}

	return session, nil
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
