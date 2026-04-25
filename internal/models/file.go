package models

import (
	"time"

	"be-file-uploader/internal/models/relations"

	"github.com/uptrace/bun"
)

type File struct {
	bun.BaseModel `bun:"table:files,alias:f"`

	ID           int            `bun:"id,pk,autoincrement" json:"id"`
	R2Key        string         `bun:"r2_key,unique" json:"-"`
	URL          string         `bun:"url" json:"-"`
	OriginalName string         `bun:"original_name" json:"original_name"`
	MimeType     string         `bun:"mime_type" json:"mime_type"`
	Size         int64          `bun:"size" json:"size"`
	UploadedBy   int            `bun:"uploaded_by" json:"-"`
	Uploader     relations.User `bun:"rel:belongs-to,join:uploaded_by=id" json:"uploader"`
	IsPrivate    bool           `bun:"is_private,default:false" json:"is_private"`
	AlbumID      *int           `bun:"album_id" json:"-"`
	Album        *Album         `bun:"rel:belongs-to,join:album_id=id" json:"album,omitempty"`
	Comments     []FileComment  `bun:"rel:has-many,join:id=image_id" json:"comments,omitempty"`
	Likes        []FileLike     `bun:"rel:has-many,join:id=image_id" json:"likes,omitempty"`
	Downloads    int            `bun:"downloads" json:"downloads,omitempty"`
}

type FileLike struct {
	bun.BaseModel `bun:"table:files_likes,alias:fl"`

	ImageID  int            `bun:"image_id,pk" json:"image_id"`
	AuthorID int            `bun:"author,pk" json:"author_id"`
	Author   relations.User `bun:"rel:belongs-to,join:author=id" json:"author"`
}

type FileComment struct {
	bun.BaseModel `bun:"table:files_comments,alias:fc"`

	ID        int            `bun:"id,pk,autoincrement" json:"id"`
	AuthorID  int            `bun:"author" json:"-"`
	Author    relations.User `bun:"rel:belongs-to,join:author=id" json:"author"`
	ImageID   int            `bun:"image_id" json:"image_id"`
	Content   string         `bun:"content" json:"content"`
	CreatedAt time.Time      `bun:"created_at,default:current_timestamp" json:"created_at"`
}

// type FileDownloads struct {
// 	ImageID  int            `bun:"image_id,pk" json:"image_id"`
// 	AuthorID int            `bun:"author" json:"-"`
// 	Author   relations.User `bun:"rel:belongs-to,join:author=id" json:"author"`
// }
