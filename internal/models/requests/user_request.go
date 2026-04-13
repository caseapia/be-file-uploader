package requests

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
