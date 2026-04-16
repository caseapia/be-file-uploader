package requests

type CreateAlbum struct {
	AlbumName string `json:"album_name" validate:"required,min=3,max=64"`
	IsPrivate bool   `json:"is_private"`
}
