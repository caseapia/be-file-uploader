package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Notification struct {
	bun.BaseModel `bun:"table:notifications" alias:"n"`

	ID        int       `bun:"id,pk,notnull" json:"id"`
	UserID    int       `bun:"user_id,notnull" json:"-"`
	Content   string    `bun:"content" json:"content"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	IsReaded  bool      `bun:"is_readed,notnull" json:"is_readed"`
}
