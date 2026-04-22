package models

import (
	"time"

	"be-file-uploader/internal/models/relations"
)

type Roadmap struct {
	ID          int            `bun:"id,pk,autoincrement" json:"id"`
	Title       string         `bun:"title" json:"title"`
	Tasks       []string       `bun:"tasks,type:json" json:"tasks"`
	CreatedAt   time.Time      `bun:"created_at,default:current_timestamp" json:"created_at"`
	CreatedByID int            `bun:"created_by" json:"-"`
	CreatedBy   relations.User `bun:"rel:belongs-to,join:created_by=id" json:"created_by"`
	FinishedAt  time.Time      `bun:"finished_at" json:"finished_at"`
}
