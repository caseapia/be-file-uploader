package models

import (
	"time"

	"be-file-uploader/internal/models/relations"
	userRelation "be-file-uploader/pkg/enums/user"

	"github.com/uptrace/bun"
)

type Restriction struct {
	bun.BaseModel `bun:"table:user_restrictions,alias:r"`

	ID           *int                   `bun:"id,pk,autoincrement" json:"id"`
	UserID       int                    `bun:"user_id" json:"-"`
	User         *relations.User        `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	ModeratorID  int                    `bun:"moderator_id" json:"-"`
	Moderator    *relations.User        `bun:"rel:belongs-to,join:moderator_id=id" json:"moderator"`
	CreatedAt    *time.Time             `bun:"created_at,default:current_timestamp" json:"created_at"`
	UnbanAt      *time.Time             `bun:"unban_at,default:null" json:"unban_at"`
	Reason       string                 `bun:"reason" json:"reason"`
	Status       userRelation.BanStatus `bun:"status" json:"status"`
	UnbannedByID *int                   `bun:"unbanned_by" json:"-"`
	UnbannedBy   *relations.User        `bun:"rel:belongs-to,join:unbanned_by=id" json:"unbanned_by,omitempty"`
	Type         userRelation.BanType   `bun:"type" json:"type"`
}
