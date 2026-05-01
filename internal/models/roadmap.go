package models

import (
	"time"

	"be-file-uploader/internal/models/relations"
	roadmapEnum "be-file-uploader/pkg/enums/roadmap"

	"github.com/uptrace/bun"
)

type RoadmapTask struct {
	bun.BaseModel `bun:"table:roadmap"`

	ID        int                `bun:"id,pk,autoincrement" json:"id"`
	Title     string             `bun:"title" json:"title"`
	Status    roadmapEnum.Status `bun:"status,default:0" json:"status"`
	CreatedAt time.Time          `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt *time.Time         `bun:"updated_at" json:"updated_at"`
	CreatorID int                `bun:"created_by" json:"-"`
	CreatedBy relations.User     `bun:"rel:belongs-to,join:created_by=id" json:"created_by"`
	UpdatorID *int               `bun:"updated_by" json:"-"`
	UpdatedBy *relations.User    `bun:"rel:belongs-to,join:updated_by=id" json:"updated_by"`
}
