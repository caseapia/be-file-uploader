package models

import "be-file-uploader/internal/models/relations"

type Invite struct {
	ID        int            `bun:"id,pk,autoincrement" json:"id"`
	Code      string         `bun:"code,unique" json:"code"`
	CreatedBy int            `bun:"created_by" json:"-"`
	Creator   relations.User `bun:"rel:belongs-to,join:created_by=id" json:"creator"`
	IsActive  bool           `bun:"is_active" json:"is_active"`
	UsedBy    *int           `bun:"used_by,nullzero" json:"-"`
	User      relations.User `bun:"rel:belongs-to,join:used_by=id" json:"user"`
}
