package role

type Permission string

const (
	FileUpload         Permission = "UPLOAD_FILES"         // FileUpload grants permission to upload files to the system.
	ViewOwnFiles       Permission = "VIEW_OWN_FILES"       // ViewOwnFiles grants permission to view files owned by the authenticated user.
	ViewOtherFiles     Permission = "VIEW_OTHER_FILES"     // ViewOtherFiles grants permission to view files that are not owned by the authenticated user.
	DownloadOtherFiles Permission = "DOWNLOAD_OTHER_FILES" // DownloadOtherFiles grants permission to download files that are not owned by the authenticated user.
	ViewOtherProfiles  Permission = "VIEW_OTHER_PROFILES"  // ViewOtherProfiles grants permission to view profiles of other users.
	ManageUsers        Permission = "MANAGE_USERS"         // ManageUsers grants permission to manage users.
	ManageFiles        Permission = "MANAGE_FILES"         // ManageFiles grants permission to manage files.
	ManageRoles        Permission = "MANAGE_ROLES"         // ManageRoles grants permission to manage roles.
	Admin              Permission = "ADMIN"                // Admin grants permission to all functionality and services.
	ViewPrivateData    Permission = "VIEW_PRIVATE_DATA"    // ViewPrivateData grants permission to view private data.
	AdminCP            Permission = "ADMIN_CP"             // AdminCP grants permission to access the Control Panel.
	ShowBadge          Permission = "SHOW_BADGE"           // ShowBadge when this flag is set, the user will be able to display the badge on their profile.
	Developer          Permission = "DEVELOPER"            // Developer grants permission to access developer tools, runtime logging, and roadmap list editing.
)
