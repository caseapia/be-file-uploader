package requests

type UploadImage struct {
	IsPrivate string `json:"is_private" form:"is_private" validate:"required"`
}

type DeleteImage struct {
	ImageID int `json:"image_id" validate:"required,min=1"`
}

type AddImageInAlbum struct {
	ImageID int `json:"image_id" validate:"required,min=1"`
	AlbumID int `json:"album_id" validate:"required,min=1"`
}

type RemoveImageFromAlbum struct {
	ImageID int `json:"image_id" validate:"required,min=1"`
}

type AddCommentToPost struct {
	PostID  int    `json:"post_id" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1"`
}
