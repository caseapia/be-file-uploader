package models

import (
	"time"

	"be-file-uploader/internal/models/relations"

	"github.com/uptrace/bun"
)

type Album struct {
	bun.BaseModel `bun:"table:albums,alias:al"`

	ID          int            `bun:"id,pk,autoincrement" json:"id"`
	Name        string         `bun:"name" json:"name"`
	CreatedByID int            `bun:"created_by" json:"-"`
	CreatedBy   relations.User `bun:"rel:belongs-to,join:created_by=id" json:"created_by"`
	CreatedAt   time.Time      `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time      `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	Items       []Image        `bun:"rel:has-many,join:id=album_id" json:"items,omitempty"`
	Options     AlbumOptions   `bun:"embed:" json:"options"`
}

type AlbumOptions struct {
	IsPublic bool `bun:"is_public,default:false" json:"is_public"`
}
