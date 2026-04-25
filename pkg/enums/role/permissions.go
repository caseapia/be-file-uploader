package role

type Permission string

const (
	FileUpload         Permission = "UPLOAD_FILES"
	ViewOwnFiles       Permission = "VIEW_OWN_FILES"
	ViewOtherFiles     Permission = "VIEW_OTHER_FILES"
	DownloadOtherFiles Permission = "DOWNLOAD_OTHER_FILES"
	ViewOtherProfiles  Permission = "VIEW_OTHER_PROFILES"
	ManageUsers        Permission = "MANAGE_USERS"
	ManageFiles        Permission = "MANAGE_FILES"
	ManageRoles        Permission = "MANAGE_ROLES"
	Admin              Permission = "ADMIN"
	ViewPrivateData    Permission = "VIEW_PRIVATE_DATA"
	InviteUsers        Permission = "INVITE_USERS"
	AdminCP            Permission = "ADMIN_CP"
)
