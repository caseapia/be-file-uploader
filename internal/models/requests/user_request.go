package requests

import (
	"time"

	userEnum "be-file-uploader/pkg/enums/user"
)

type SetUploadLimitRequest struct {
	User  int   `json:"user" validate:"required"`
	Limit int64 `json:"limit" validate:"required"`
}

type AddUserInRole struct {
	User int `json:"user" validate:"required"`
	Role int `json:"role" validate:"required"`
}

type RemoveUserFromRole struct {
	User int `json:"user" validate:"required"`
	Role int `json:"role" validate:"required"`
}

type RestrictUser struct {
	User    int              `json:"user" validate:"required"`
	UnbanAt time.Time        `json:"unban_at"`
	Type    userEnum.BanType `json:"type" validate:"required"`
	Reason  string           `json:"reason" validate:"required,max=255"`
}
