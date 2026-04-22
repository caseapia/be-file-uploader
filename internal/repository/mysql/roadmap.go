package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) GetRoadmapList(ctx context.Context) (roadmap []*models.Roadmap, err error) {
	roadmap = make([]*models.Roadmap, 0)

	err = r.DB.NewSelect().
		Model(roadmap).
		Scan(ctx)

	return roadmap, err
}

func (r *Repository) AddRoadmapTask(ctx context.Context, tx bun.IDB, task models.Roadmap) error {
	_, err := tx.NewInsert().
		Model(task).
		Exec(ctx)

	return err
}
