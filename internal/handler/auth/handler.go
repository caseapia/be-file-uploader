package auth

import (
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/service/auth"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	authService *auth.Service
}

func NewHandler(auth *auth.Service) *Handler {
	return &Handler{authService: auth}
}

func (h *Handler) Login(ctx fiber.Ctx) error {
	var req requests.Login
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	user, access, refresh, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, &fiber.Map{
		"user":          user,
		"refresh_token": refresh,
		"access_token":  access,
	})
}

func (h *Handler) Register(ctx fiber.Ctx) error {
	var req requests.Register
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	user, access, refresh, err := h.authService.Register(ctx, req.Username, req.Password, req.InviteCode)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, &fiber.Map{
		"user":          user,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *Handler) Refresh(ctx fiber.Ctx) error {
	var req requests.Refresh
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	access, refresh, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, &fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
