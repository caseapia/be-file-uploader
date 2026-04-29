package notification

import (
	"strconv"

	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/notification"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	notificationService *notification.Service
	repository          *mysql.Repository
}

func NewHandler(notification *notification.Service, repository *mysql.Repository) *Handler {
	return &Handler{notificationService: notification, repository: repository}
}

func (h *Handler) SearchMyNotifications(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	notifications, err := h.notificationService.LookupMyNotifications(ctx, sender)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, notifications)
}

func (h *Handler) ReadNotification(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)
	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	status, err := h.notificationService.ReadNotification(ctx.Context(), sender, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, status)
}
