package requests

type Login struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type Register struct {
	Username   string `json:"username" validate:"required,min=3,max=20"`
	Password   string `json:"password" validate:"required,min=6,max=32"`
	InviteCode string `json:"invite_code" validate:"required,min=6,max=6"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token" validate:"required,len=44"`
}
