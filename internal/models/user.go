package models

import (
	"time"

	"be-file-uploader/pkg/enums/role"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          int                    `bun:"id,pk,autoincrement,unique" json:"id"`
	Username    string                 `bun:"username,unique" json:"username"`
	DiscordUID  *int                   `bun:"discord_uid,unique" json:"discord_uid"`
	DiscordName *string                `bun:"discord_name" json:"discord_name"`
	Password    string                 `bun:"password" json:"-"`
	CreatedAt   time.Time              `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time              `bun:"updated_at,nullzero" json:"updated_at"`
	Roles       []Role                 `bun:"m2m:user_roles,join:User=Role" json:"roles"`
	RegisterIP  string                 `bun:"register_ip" json:"-"`
	LastIP      string                 `bun:"last_ip" json:"-"`
	Useragent   string                 `bun:"useragent" json:"-"`
	Private     map[string]interface{} `bun:"-" json:"private,omitempty"`
	Storage     []File                 `bun:"rel:has-many,join:id=uploaded_by" json:"images"`
	UploadLimit int64                  `bun:"upload_limit,default:1073741824" json:"upload_limit"`
	UsedStorage int64                  `bun:"used_storage,default:0" json:"used_storage"`
	IsVerified  bool                   `bun:"is_verified,default:false" json:"is_verified"`
	CFRayID     string                 `bun:"cf_ray_id" json:"-"`
	Albums      []Album                `bun:"rel:has-many,join:id=created_by" json:"albums"`
	Locale      string                 `bun:"locale" json:"-"`
	ShareXToken *string                `bun:"sharex_token" json:"-"`
	LastSeen    time.Time              `bun:"last_seen,default:current_timestamp" json:"last_seen"`
	GeoString   string                 `bun:"geo_string" json:"-"`
	Geolocation Geolocation            `bun:"-" json:"geolocation"`
}

type UserRole struct {
	bun.BaseModel `bun:"table:user_roles"`

	UserID int   `bun:"user_id,pk"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`

	RoleID int   `bun:"role_id,pk"`
	Role   *Role `bun:"rel:belongs-to,join:role_id=id"`
}

type Geolocation struct {
	Code    string
	Country string
	City    string
}

func (u *User) HasPermission(permission role.Permission) bool {
	for _, r := range u.Roles {
		if r.HasPermission(permission) {
			return true
		}
	}
	return false
}

func (u *User) GetPrivateData() map[string]interface{} {
	return map[string]interface{}{
		"register_ip":  u.RegisterIP,
		"last_ip":      u.LastIP,
		"useragent":    u.Useragent,
		"cf_ray_id":    u.CFRayID,
		"locale":       u.Locale,
		"sharex_token": u.ShareXToken,
	}
}
