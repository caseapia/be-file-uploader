package relations

import (
	"be-file-uploader/pkg/enums/role"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID       int    `bun:"id,pk,autoincrement,unique" json:"id"`
	Username string `bun:"username,unique" json:"username"`
	Avatar   string `bun:"avatar" json:"avatar"`
}

type UserRole struct {
	ID          int               `bun:"id,pk,autoincrement" json:"id"`
	Color       string            `bun:"color" json:"color"`
	Name        string            `bun:"name" json:"name"`
	Permissions []role.Permission `bun:"permissions,type:json" json:"permissions"`
}
