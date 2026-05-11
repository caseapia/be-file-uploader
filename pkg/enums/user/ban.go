package userEnum

type BanStatus string

const (
	BanStatusBanned          BanStatus = "banned"
	BanStatusUnbannedByAdmin BanStatus = "unbanned_by_admin"
	BanStatusExpired         BanStatus = "ban_expired"
)

type BanType string

const (
	BanTypeAccount BanType = "account"
	BanTypeUpload  BanType = "upload"
	BanTypeLike    BanType = "like"
	BanTypeComment BanType = "comment"
)
