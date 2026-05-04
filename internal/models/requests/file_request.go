package requests

type Part struct {
	PartNumber int32  `json:"part_number"`
	ETag       string `json:"etag"`
}

type InitUpload struct {
	OriginalName string `json:"original_name" validate:"required"`
	MimeType     string `json:"mime_type" validate:"required"`
	Size         int64  `json:"size" validate:"required"`
	IsPrivate    bool   `json:"is_private"`
}

type InitUploadResponse struct {
	UploadID string `json:"upload_id"`
	Key      string `json:"key"`
}

type CompleteUpload struct {
	UploadID     string `json:"upload_id" validate:"required"`
	Key          string `json:"key" validate:"required"`
	OriginalName string `json:"original_name" validate:"required"`
	MimeType     string `json:"mime_type" validate:"required"`
	Size         int64  `json:"size" validate:"required"`
	IsPrivate    bool   `json:"is_private"`
	Parts        []Part `json:"parts" validate:"required"`
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

type EditAccess struct {
	FileID int `json:"file_id" validate:"required,min=1"`
	UserID int `json:"user_id" validate:"required,min=1"`
}

type EditFileDetails struct {
	FileID    int    `json:"file_id" validate:"required,min=1"`
	FileName  string `json:"file_name" validate:"required"`
	IsPrivate bool   `json:"is_private" validate:"boolean"`
}
