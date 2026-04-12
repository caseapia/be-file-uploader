package relations

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID       int    `bun:"id,pk,autoincrement,unique" json:"id"`
	Username string `bun:"username,unique" json:"username"`
}
