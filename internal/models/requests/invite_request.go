package requests

type RevokeInvite struct {
	InviteID int `json:"invite_id" validate:"required,min=1"`
}
