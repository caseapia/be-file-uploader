package requests

import "be-file-uploader/pkg/enums/role"

type CreateRole struct {
	Name        string            `json:"name" validate:"required"`
	Color       string            `json:"color" validate:"required,hexcolor"`
	IsSystem    bool              `json:"is_system"`
	Permissions []role.Permission `json:"permissions" validate:"required"`
}

type UpdateRole struct {
	RoleID      int               `json:"role_id" validate:"required"`
	Name        string            `json:"name" validate:"required"`
	Color       string            `json:"color" validate:"required,hexcolor"`
	IsSystem    bool              `json:"is_system"`
	Permissions []role.Permission `json:"permissions" validate:"required"`
}

type DeleteRole struct {
	ID int `json:"id" validate:"required"`
}
