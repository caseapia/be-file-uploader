package requests

type DeleteImage struct {
	ImageID int `json:"image_id" validate:"required,min=1"`
}
