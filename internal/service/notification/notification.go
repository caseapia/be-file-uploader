package notification

import (
	"context"
	"time"

	"be-file-uploader/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
)

func (s *Service) LookupMyNotifications(ctx fiber.Ctx, sender *models.User) (notifications *[]models.Notification, err error) {
	notifications, err = s.repo.LookupNotificationsByUserID(ctx.Context(), sender)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *Service) CreateNotification(ctx context.Context, user int, content string) {
	notification := models.Notification{
		UserID:    user,
		Content:   content,
		CreatedAt: time.Now(),
		IsReaded:  false,
	}

	_, err := s.repo.CreateNotification(ctx, s.repo.DB, notification)
	if err != nil {
		slog.Errorf("Failed to create notification: %v", err)
		return
	}
}

func (s *Service) ReadNotification(ctx context.Context, sender *models.User, notificationID int) (status bool, err error) {
	notification, err := s.repo.SearchNotificationByID(ctx, notificationID)
	if err != nil {
		return false, fiber.NewError(fiber.StatusNotFound, "ERR_NOTIFICATION_NOTFOUND")
	}
	if notification.UserID != sender.ID {
		return false, fiber.NewError(fiber.StatusForbidden, "ERR_NOTIFICATION_FORBIDDEN")
	}

	status, err = s.repo.ReadNotification(ctx, s.repo.DB, notification.ID)
	if err != nil {
		return false, err
	}

	return status, nil
}
