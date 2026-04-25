package image

var allowedMime = map[string]string{
	"image/jpeg":                   ".jpg",
	"image/png":                    ".png",
	"image/webp":                   ".webp",
	"image/gif":                    ".gif",
	"application/pdf":              ".pdf",
	"text/plain":                   ".txt",
	"application/zip":              ".zip",
	"application/x-rar-compressed": ".rar",
	"application/x-7z-compressed":  ".7z",
}

const maxFileSize = 200 << 20
