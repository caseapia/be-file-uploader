package roadmap

import (
	"fmt"
	"time"

	"be-file-uploader/internal/models"
	roadmapEnum "be-file-uploader/pkg/enums/roadmap"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func (s *Service) RoadmapList(ctx fiber.Ctx) (roadmap []models.RoadmapTask, err error) {
	roadmap, err = s.repo.GetRoadmapList(ctx.Context())
	if err != nil {
		return nil, err
	}

	return roadmap, nil
}

func (s *Service) AddTask(ctx fiber.Ctx, sender *models.User, title string) (task *models.RoadmapTask, err error) {
	task = &models.RoadmapTask{
		Title:     title,
		CreatedAt: time.Now(),
		CreatorID: sender.ID,
		Status:    roadmapEnum.Planned,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		err = s.repo.AddTask(ctx.Context(), tx, *task)
		if err != nil {
			return err
		}

		s.notify.CreateNotification(ctx.Context(), sender.ID, fmt.Sprintf("NOTIFY_ROADMAP_ADD_TASK+%s", task.Title))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) EditTask(ctx fiber.Ctx, sender *models.User, id int, title string, status roadmapEnum.Status) (task *models.RoadmapTask, err error) {
	now := time.Now()

	task = &models.RoadmapTask{
		ID:        id,
		Title:     title,
		Status:    status,
		UpdatedAt: &now,
		UpdatorID: &sender.ID,
	}

	err = s.repo.EditTask(ctx.Context(), s.repo.DB, *task)
	if err != nil {
		return nil, err
	}

	return task, nil
}
