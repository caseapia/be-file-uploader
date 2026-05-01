package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) GetRoadmapList(ctx context.Context) (roadmap []models.RoadmapTask, err error) {
	roadmap = make([]models.RoadmapTask, 0)

	err = r.DB.NewSelect().
		Model(&roadmap).
		Scan(ctx)
	return roadmap, err
}

func (r *Repository) AddTask(ctx context.Context, tx bun.IDB, task models.RoadmapTask) error {
	_, err := tx.NewInsert().
		Model(&task).
		Exec(ctx)
	return err
}

func (r *Repository) EditTask(ctx context.Context, tx bun.IDB, task models.RoadmapTask) error {
	_, err := tx.NewUpdate().
		Model(&task).
		Where("id = ?", task.ID).
		Exec(ctx)
	return err
}

func (r *Repository) SearchTaskByID(ctx context.Context, id int) (*models.RoadmapTask, error) {
	task := new(models.RoadmapTask)

	err := r.DB.NewSelect().
		Model(task).
		Where("id = ?", id).
		Scan(ctx)

	return task, err
}
