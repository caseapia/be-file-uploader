package relations

type Role struct {
	ID          int      `bun:"id,pk,autoincrement,unique" json:"id"`
	Name        string   `bun:"name" json:"name"`
	Permissions []string `bun:"permissions,type:json" json:"permissions"`
	Color       string   `bun:"color" json:"color"`
}
