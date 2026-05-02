package auth

import (
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/utils/generate"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

var unhashedPassword string

func (s *Service) validateRegistration(ctx fiber.Ctx, username string) error {
	user, err := s.repo.LookupUserByName(ctx, username)
	if user != nil && err == nil {
		return fiber.NewError(fiber.StatusConflict, "ERR_USER_ALREADY_EXISTS")
	}

	return nil
}

func (s *Service) createUserWithInvite(ctx fiber.Ctx, username, password string) (user *models.User, err error) {
	ip := ctx.IP()
	useragent := ctx.Get("X-User-Agent")
	rayid := ctx.Get("Cf-Ray")
	code, country, city := s.geo.GetGeoString(ip)

	hashPassword, err := generate.HashPassword(password)
	unhashedPassword = password
	if err != nil {
		slog.WithData(slog.M{"error": err, "ip": ip}).Error("Password hashing failed")
		return nil, fiber.NewError(fiber.StatusConflict, "ERR_USER_REGISTER_HASHCREATION")
	}

	user = &models.User{
		Username:    username,
		Password:    hashPassword,
		RegisterIP:  ip,
		LastIP:      ip,
		Useragent:   useragent,
		DiscordName: nil,
		DiscordUID:  nil,
		CreatedAt:   time.Now(),
		UploadLimit: 1073741824,
		UsedStorage: 0,
		CFRayID:     rayid,
		GeoString:   country + ", " + city,
		Geolocation: &models.Geolocation{
			CountryCode: code,
			City:        city,
			Country:     country,
		},
	}

	err = s.repo.WithTx(ctx, func(tx bun.Tx) error {
		return s.registerTx(ctx, tx, user)
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

func (s *Service) registerTx(ctx fiber.Ctx, tx bun.Tx, user *models.User) error {
	created, err := s.repo.CreateUser(ctx.Context(), tx, *user)
	if err != nil {
		slog.WithData(slog.M{"error": err, "user": user}).Error("Failed to create user")
		return fiber.NewError(fiber.StatusConflict, "ERR_USER_UNKNOWN_CREATION_ERROR")
	}

	if err := s.repo.AddUserInRole(ctx, tx, created.ID, 1); err != nil {
		slog.WithData(slog.M{"error": err, "user": created}).Error("Failed to add user in default role")
		return fiber.NewError(fiber.StatusConflict, "ERR_USER_ADDINROLE_FAILED")
	}

	*user = *created
	return nil
}

func (s *Service) createSessionTokens(ctx fiber.Ctx, user *models.User) (access, refresh string, err error) {
	user, access, refresh, err = s.Login(ctx, user.Username, unhashedPassword)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Register(ctx fiber.Ctx, username, password string) (*models.User, string, string, error) {
	err := s.validateRegistration(ctx, username)
	if err != nil {
		return nil, "", "", err
	}

	user, err := s.createUserWithInvite(ctx, username, password)
	if err != nil {
		return nil, "", "", err
	}

	access, refresh, err := s.createSessionTokens(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	return user, access, refresh, nil
}
