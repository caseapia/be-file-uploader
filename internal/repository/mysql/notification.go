package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) LookupNotificationsByUserID(ctx context.Context, user *models.User) (*[]models.Notification, error) {
	notifications := make([]models.Notification, 0)

	err := r.DB.NewSelect().
		Model(&notifications).
		Where("user_id = ?", user.ID).
		Order("is_readed ASC").
		Scan(ctx)

	return &notifications, err
}

func (r *Repository) CreateNotification(ctx context.Context, tx bun.IDB, notification models.Notification) (status bool, err error) {
	_, err = tx.NewInsert().
		Model(&notification).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) ReadNotification(ctx context.Context, tx bun.IDB, notificationID int) (status bool, err error) {
	notification := new(models.Notification)

	_, err = tx.NewUpdate().
		Model(notification).
		Where("id = ?", notificationID).
		Set("is_readed = ?", true).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) SearchNotificationByID(ctx context.Context, notificationID int) (*models.Notification, error) {
	notification := new(models.Notification)

	err := r.DB.NewSelect().
		Model(notification).
		Where("id = ?", notificationID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return notification, nil
}
