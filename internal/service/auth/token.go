package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/utils/token"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/slog"
)

func (s *Service) ValidateAccessToken(tokenStr string) (int, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenStr, &models.SessionClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_UNEXPECTED_SIGNING_METHOD")
		}
		return models.JWTSecret, nil
	})
	if err != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	claims, ok := token.Claims.(*models.SessionClaims)
	if !ok || !token.Valid {
		return 0, fiber.NewError(fiber.StatusInternalServerError, "ERR_INVALID_TOKEN")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return 0, fiber.NewError(fiber.StatusInternalServerError, "ERR_TOKEN_EXPIRED")
	}

	return claims.UserID, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *Service) ParseJWT(tokenString string) (*models.User, *models.SessionClaims, error) {
	claims, err := token.ParseAccessToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.repo.LookupUserByID(context.Background(), claims.UserID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		slog.WithData(slog.M{
			"error":  err,
			"user":   user,
			"claims": claims,
		}).Error("user seems to be nil on JWT Parsing")
		return nil, nil, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOT_FOUND")
	}

	return user, claims, nil
}
