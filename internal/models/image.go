package models

import (
	"be-file-uploader/internal/models/relations"

	"github.com/uptrace/bun"
)

type Image struct {
	bun.BaseModel `bun:"table:images,alias:i"`

	ID           int            `bun:"id,pk,autoincrement" json:"id"`
	R2Key        string         `bun:"r2_key,unique" json:"-"`
	URL          string         `bun:"url" json:"url"`
	OriginalName string         `bun:"original_name" json:"original_name"`
	MimeType     string         `bun:"mime_type" json:"mime_type"`
	Size         int64          `bun:"size" json:"size"`
	UploadedBy   int            `bun:"uploaded_by" json:"-"`
	Uploader     relations.User `bun:"rel:belongs-to,join:uploaded_by=id" json:"uploader"`
	IsPrivate    bool           `bun:"is_private,default:false" json:"is_private"`
	AlbumID      *int           `bun:"album_id" json:"-"`
	Album        *Album         `bun:"rel:belongs-to,join:album_id=id" json:"album,omitempty"`
	Views        int            `bun:"views,default:0" json:"views"`
}
