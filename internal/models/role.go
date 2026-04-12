package models

import (
	"time"

	"be-file-uploader/pkg/enums/role"

	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID          int               `bun:"id,pk,autoincrement,unique" json:"id"`
	Name        string            `bun:"name" json:"name"`
	Permissions []role.Permission `bun:"permissions,type:json" json:"permissions"`
	IsSystem    bool              `bun:"is_system,default:false" json:"is_system"`
	CreatedAt   time.Time         `bun:"created_at,default:current_timestamp" json:"created_at"`
	CreatedBy   int               `bun:"created_by" json:"created_by"`
	Color       string            `bun:"color" json:"color"`
}

func (r Role) HasPermission(permission role.Permission) bool {
	for _, p := range r.Permissions {
		if p == permission || p == "ADMIN" {
			return true
		}
	}
	return false
}
