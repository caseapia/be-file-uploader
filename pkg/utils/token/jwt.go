package token

import (
	"time"

	"be-file-uploader/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/slog"
)

func GenerateAccessToken(userID, tokenVer int, sessionID string) (string, error) {
	claims := models.SessionClaims{
		UserID:    userID,
		TokenVer:  tokenVer,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(models.JWTSecret)
}

func ParseAccessToken(tokenStr string) (*models.SessionClaims, error) {
	if len(tokenStr) == 0 {
		slog.WithData(slog.M{
			"tokenStr": tokenStr,
		}).Error("tokenStr empty")

		return nil, jwt.ErrTokenMalformed
	}

	claims := &models.SessionClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.WithData(slog.M{
				"tokenStr": tokenStr,
				"claims":   claims,
			}).Error("ErrTokenMalformed")

			return nil, jwt.ErrTokenMalformed
		}

		return models.JWTSecret, nil
	})
	if err != nil || !token.Valid {
		slog.WithData(slog.M{
			"tokenStr": tokenStr,
			"err":      err,
			"claims":   claims,
		}).Error("Failed to parse JWT")
		return nil, err
	}

	return claims, nil
}

func GetJWTSecret() []byte {
	return models.JWTSecret
}
