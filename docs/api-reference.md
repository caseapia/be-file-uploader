# API Reference

Base URL:

```text
http://localhost:{APP_PORT}/v1/api
```

## Common Rules

### Success Envelope

Most successful JSON responses are wrapped in a top-level `response` field.

```json
{
  "response": {}
}
```

Exception:

- `POST /public/storage/upload/sharex` returns `{ "url": "..." }` directly for ShareX compatibility.

### Error Envelopes

Handled application errors:

```json
{
  "error": "ERR_CODE",
  "code": 400
}
```

Unhandled internal errors:

```json
{
  "code": 500,
  "message": "raw error text"
}
```

### Authentication

- Public routes live under `/public`.
- Private routes live under `/private`.
- Private routes require an access token:
  - `Authorization: Bearer <jwt>`, or
  - `Cookie: auth_token=<jwt>`

### Frontend Headers

Recommended on every request:

```http
Content-Type: application/json
X-User-Agent: your-frontend-name/version
```

For multipart upload chunks and ShareX upload, use `multipart/form-data` instead of JSON.

## Data Shapes

### User

```json
{
  "id": 1,
  "username": "alice",
  "discord_uid": null,
  "discord_name": null,
  "created_at": "2026-04-11T12:00:00Z",
  "updated_at": "2026-04-11T12:00:00Z",
  "roles": [],
  "private": {
    "register_ip": "127.0.0.1",
    "last_ip": "127.0.0.1",
    "useragent": "frontend/1.0.0",
    "cf_ray_id": "",
    "locale": "en",
    "sharex_token": "token"
  },
  "images": [],
  "upload_limit": 1073741824,
  "used_storage": 0,
  "is_verified": false,
  "albums": []
}
```

Notes:

- `private` is present only when the requester has `VIEW_PRIVATE_DATA`.
- Password, IP fields, locale, Cloudflare ray ID, and ShareX token are hidden unless exposed through `private`.

### File

```json
{
  "id": 10,
  "url": "https://cdn.example.com/images/2/2026-04/file.jpg",
  "original_name": "photo.png",
  "mime_type": "image/png",
  "size": 245901,
  "uploader": {
    "id": 2,
    "username": "alice"
  },
  "is_private": false,
  "album": null,
  "comments": [],
  "likes": [],
  "downloads": 0,
  "created_at": "2026-04-11T12:00:00Z"
}
```

Notes:

- Storage key `r2_key` is intentionally hidden from JSON responses.
- `url` can be blank when the file is private or non-image content should not expose a direct URL.

### Album

```json
{
  "id": 1,
  "name": "Screenshots",
  "created_by": {
    "id": 2,
    "username": "alice"
  },
  "created_at": "2026-04-11T12:00:00Z",
  "updated_at": "2026-04-11T12:00:00Z",
  "items": [],
  "options": {
    "is_public": true
  }
}
```

### Role

```json
{
  "id": 1,
  "name": "member",
  "permissions": ["UPLOAD_FILES", "VIEW_OWN_FILES"],
  "is_system": false,
  "created_at": "2026-04-11T12:00:00Z",
  "color": "#00a86b"
}
```

### Notification

```json
{
  "id": 1,
  "content": "NOTIFY_VERIFY_SUCCESS",
  "created_at": "2026-04-11T12:00:00Z",
  "is_readed": false
}
```

### Roadmap Task

```json
{
  "id": 1,
  "title": "Add public profiles",
  "status": 0,
  "created_at": "2026-04-11T12:00:00Z",
  "updated_at": null,
  "created_by": {
    "id": 1,
    "username": "admin"
  },
  "updated_by": null
}
```

## Public Routes

### `GET /public/ping`

Health-style utility route.

Response `200`:

```json
{
  "response": "pong"
}
```

### `POST /public/auth/register`

Creates a new account and immediately issues an access/refresh token pair.

Request body:

```json
{
  "username": "alice",
  "password": "secret12"
}
```

Validation:

- `username`: required, min `3`, max `20`
- `password`: required, min `6`, max `32`

Response `200`:

```json
{
  "response": {
    "user": {},
    "access_token": "jwt",
    "refresh_token": "44-char-token"
  }
}
```

Behavior:

- Checks that username does not already exist.
- Creates the user.
- Assigns the new user to role ID `1`.
- Creates a session and issues tokens.

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | validator message | Bad body or failed field validation |
| `409` | `ERR_USER_ALREADY_EXISTS` | Username already taken |
| `409` | `ERR_USER_REGISTER_HASHCREATION` | Password hash creation failed |
| `409` | `ERR_USER_UNKNOWN_CREATION_ERROR` | User creation failed inside transaction |
| `409` | `ERR_USER_ADDINROLE_FAILED` | Default role assignment failed |
| `500` | `ERR_USER_LOOKUP_AFTER_REGISTER` | User re-read failed after creation |
| `500` | `ERR_SESSION_CREATION` | Session row could not be created |
| `500` | `ERR_TOKEN_GENERATION` | Token generation failed |

### `POST /public/auth/login`

Authenticates an existing user and creates a new DB session.

Request body:

```json
{
  "username": "alice",
  "password": "secret12"
}
```

Validation:

- `username`: required, min `3`, max `20`
- `password`: required, min `6`, max `32`

Response `200`:

```json
{
  "response": {
    "user": {},
    "access_token": "jwt",
    "refresh_token": "44-char-token"
  }
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | validator message | Bad body or failed field validation |
| `404` | `ERR_WRONG_CREDENTIALS` | Username not found or password mismatch |
| `500` | `ERR_SESSION_CREATION` | Session row could not be created |
| `500` | `ERR_TOKEN_GENERATION` | Access/refresh token generation failed |

### `POST /public/auth/refresh`

Rotates the refresh token and returns a new access/refresh pair.

Request body:

```json
{
  "refresh_token": "existing-44-char-token"
}
```

Validation:

- `refresh_token`: required, exact length `44`

Response `200`:

```json
{
  "response": {
    "access_token": "jwt",
    "refresh_token": "new-44-char-token"
  }
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | validator message | Bad body or failed field validation |
| `500` | `ERR_TOKEN_GENERATION` | Invalid/missing/expired refresh session or token generation failure |
| `500` | `ERR_USER_NOTFOUND` | Session user no longer exists |

### `GET /public/roadmap/list`

Returns roadmap tasks.

Response `200`:

```json
{
  "response": {
    "roadmap": []
  }
}
```

### `POST /public/storage/upload/sharex`

Uploads one file through a generated ShareX token.

Request:

- content type: `multipart/form-data`
- form field `token`: ShareX token from `GET /private/user/shareX/generate`
- form field `image`: uploaded file

Response `200`:

```json
{
  "url": "https://cdn.example.com/images/..."
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | `ERR_IMAGE_MISSING` | No multipart file named `image` |
| `401` | `ERR_INVALID_TOKEN` | ShareX token did not match a user |
| `413` | `ERR_IMAGE_TOO_LARGE` / `ERR_QUOTA_EXCEEDED` | File exceeds max size or user's quota |
| `500` | `ERR_MIMETYPE` | MIME type is not supported |
| `500` | `ERR_UPLOAD_CHUNK` | Object part upload failed |
| `500` | `ERR_COMPLETE_MULTIPART` | Multipart completion failed |

## Private Auth Routes

All routes below require a valid access token.

### `DELETE /private/auth/logout`

Disables the current session.

Response `200`:

```json
{
  "response": "OK"
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `404` | `ERR_USER_NOTFOUND` | User no longer exists |
| `404` | `ERR_SESSION_NOTFOUND` | Session does not belong to user |
| `404` | `ERR_SESSION_NOTACTIVE` | Session already inactive |
| `500` | `ERR_SESSION_UPDATE` | Session update failed |

## Private User Routes

All routes below require a valid access token.

### `GET /private/user/me`

Returns the authenticated user.

Response `200`:

```json
{
  "response": {}
}
```

### `GET /private/user/lookup/:id`

Looks up a user by numeric ID.

Required permission:

- `VIEW_OTHER_PROFILES`

Path params:

| Param | Type | Example |
| --- | --- | --- |
| `id` | integer | `42` |

Response `200`:

```json
{
  "response": {}
}
```

Notes:

- Private user data is included only when the requester has `VIEW_PRIVATE_DATA`.
- Private files are filtered when the requester is not the target and lacks `MANAGE_FILES`.

### `GET /private/user/shareX/generate`

Generates and stores a ShareX upload token for the authenticated user.

Required permission:

- `UPLOAD_FILES`

Response `200`:

```json
{
  "response": "generated-token"
}
```

## Private User Admin Routes

All routes below require:

- a valid access token
- permission `MANAGE_USERS`

### `GET /private/user/admin/users`

Returns up to 30 users.

Response `200`:

```json
{
  "response": []
}
```

### `PUT /private/user/admin/role/add`

Adds a role to a user.

Request body:

```json
{
  "user": 2,
  "role": 3
}
```

Validation:

- `user`: required
- `role`: required

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `403` | `ERR_ROLE_ISSUE_FORBIDDEN` | Non-admin tried to assign a system role |
| `404` | `ERR_USER_NOTFOUND` | User ID does not exist |
| `404` | `ERR_ROLE_NOTFOUND...` | Role ID does not exist |

### `DELETE /private/user/admin/role/delete`

Removes a role from a user.

Request body:

```json
{
  "user": 2,
  "role": 3
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `403` | `ERR_ROLE_ISSUE_FORBIDDEN` | Non-admin tried to remove protected role |
| `404` | `ERR_USER_NOTFOUND` | User ID does not exist |
| `404` | `ERR_ROLE_NOTFOUND...` | Role ID does not exist |

### `PATCH /private/user/admin/storage-limit/update`

Updates a user's upload quota.

Request body:

```json
{
  "user": 2,
  "limit": 1073741824
}
```

### `PATCH /private/user/admin/verify/:id`

Toggles a user's verification status.

Path params:

| Param | Type | Example |
| --- | --- | --- |
| `id` | integer | `2` |

### `DELETE /private/user/admin/shareX/reset/:id`

Clears a user's ShareX token.

Path params:

| Param | Type | Example |
| --- | --- | --- |
| `id` | integer | `2` |

## Private Storage Routes

All routes below require a valid access token.

### `POST /private/storage/upload/init`

Starts a multipart upload.

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "original_name": "photo.png",
  "mime_type": "image/png",
  "size": 245901,
  "is_private": false
}
```

Validation:

- `original_name`: required
- `mime_type`: required
- `size`: required
- `is_private`: optional boolean

Accepted MIME types:

- `image/jpeg`
- `image/png`
- `image/webp`
- `image/gif`
- `application/pdf`
- `text/plain`
- `application/zip`
- `application/x-rar-compressed`
- `application/x-7z-compressed`

Limit:

- `4 GiB`
- also constrained by the user's remaining `upload_limit`

Response `200`:

```json
{
  "response": {
    "upload_id": "s3-upload-id",
    "key": "images/dev/2/2026-04/generated.png"
  }
}
```

### `POST /private/storage/upload/chunk`

Uploads one multipart part.

Required permission:

- `UPLOAD_FILES`

Request:

- content type: `multipart/form-data`
- form field `upload_id`
- form field `key`
- form field `part_number`
- file field `chunk`

Response `200`:

```json
{
  "response": {
    "eTag": "\"etag-value\""
  }
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | `ERR_INVALID_PARAMS` | Missing upload ID/key or invalid part number |
| `400` | `ERR_CHUNK_MISSING` | No multipart file named `chunk` |
| `500` | `ERR_OPEN_IMAGE` | Multipart file open failed |
| `500` | `ERR_FILE_READING` | File read failed |
| `500` | `ERR_UPLOAD_CHUNK` | Object storage part upload failed |

### `POST /private/storage/upload/complete`

Completes a multipart upload and creates the DB file row.

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "upload_id": "s3-upload-id",
  "key": "images/dev/2/2026-04/generated.png",
  "original_name": "photo.png",
  "mime_type": "image/png",
  "size": 245901,
  "is_private": false,
  "parts": [
    {
      "part_number": 1,
      "etag": "\"etag-value\""
    }
  ]
}
```

Response `201`:

```json
{
  "response": {}
}
```

### `POST /private/storage/delete`

Deletes a file by ID.

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "image_id": 10
}
```

Validation:

- `image_id`: required, integer, minimum `1`

Response `200`:

```json
{
  "response": {
    "used_storage": 12345,
    "status": "OK"
  }
}
```

Deletion succeeds when the requester owns the file or has `MANAGE_FILES`.

### `GET /private/storage/my`

Returns files uploaded by the authenticated user.

Required permission:

- `VIEW_OWN_FILES`

Response `200`:

```json
{
  "response": []
}
```

### `GET /private/storage/list`

Returns all files visible to the requester.

Required permission:

- `VIEW_OTHER_FILES`

Notes:

- Private files are filtered unless the requester has `MANAGE_FILES`.
- URLs can still be blanked by `ResolveURL`.

### `GET /private/storage/list/:id`

Returns files uploaded by the user with ID `:id`.

Required permission:

- `VIEW_OTHER_FILES`

### `GET /private/storage/post/:id`

Returns one file/post by ID.

Required permission:

- `VIEW_OTHER_FILES`

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `404` | `ERR_IMAGE_NOTFOUND` | File does not exist or is private to another user |

### `PATCH /private/storage/post/action/like/:id`

Adds a like to a file/post.

Required permission:

- `VIEW_OTHER_FILES`

Response `200`:

```json
{
  "response": true
}
```

### `DELETE /private/storage/post/action/likeRemove/:id`

Removes a like from a file/post.

Required permission:

- `VIEW_OTHER_FILES`

Response `200`:

```json
{
  "response": true
}
```

### `GET /private/storage/post/action/download/:id`

Increments download count and returns the file URL.

Required permission:

- `DOWNLOAD_OTHERS_FILES`

Response `201`:

```json
{
  "response": "https://cdn.example.com/images/..."
}
```

### `POST /private/storage/post/action/addComment`

Adds a comment to a file/post.

Required permission:

- `VIEW_OTHER_FILES`

Request body:

```json
{
  "post_id": 10,
  "content": "Looks good"
}
```

Validation:

- `post_id`: required, minimum `1`
- `content`: required, min `1`

Response `201`:

```json
{
  "response": {}
}
```

### `PUT /private/storage/album/put`

Adds a file to an album.

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "image_id": 10,
  "album_id": 1
}
```

### `DELETE /private/storage/album/delete`

Removes a file from its album.

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "image_id": 10
}
```

## Private Album Routes

All routes below require a valid access token.

### `POST /private/album/create`

Creates an album.

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "album_name": "Screenshots",
  "is_private": false
}
```

Validation:

- `album_name`: required, min `3`, max `64`
- `is_private`: optional boolean

Response `201`:

```json
{
  "response": {}
}
```

Implementation note:

- The request field is named `is_private`, but the model stores `options.is_public` from that boolean value.

### `DELETE /private/album/delete/:id`

Deletes an album.

Required permission:

- `UPLOAD_FILES`

Response `200`:

```json
{
  "response": true
}
```

### `GET /private/album/lookup/:id`

Looks up an album.

Required permission:

- `VIEW_OWN_FILES`

Private albums are visible only to the owner or users with `MANAGE_FILES`.

### `GET /private/album/lookupAll`

Returns all albums.

Required permission:

- `MANAGE_FILES`

## Private Roles Admin Routes

All routes below require:

- a valid access token
- permission `MANAGE_ROLES`

### `GET /private/roles/admin/all`

Returns all roles.

### `POST /private/roles/admin/create`

Creates a role.

Request body:

```json
{
  "name": "moderator",
  "color": "#3366ff",
  "is_system": false,
  "permissions": ["VIEW_OTHER_FILES", "MANAGE_FILES"]
}
```

Validation:

- `name`: required
- `color`: required, valid hex color
- `permissions`: required

Response `201`:

```json
{
  "response": {}
}
```

### `PATCH /private/roles/admin/edit`

Updates a role.

Request body:

```json
{
  "role_id": 3,
  "name": "moderator",
  "color": "#3366ff",
  "is_system": false,
  "permissions": ["VIEW_OTHER_FILES", "MANAGE_FILES"]
}
```

Response `201`:

```json
{
  "response": {}
}
```

### `DELETE /private/roles/admin/delete`

Deletes a role.

Request body:

```json
{
  "id": 3
}
```

Response `200`:

```json
{
  "response": "OK"
}
```

## Private Notifications Routes

All routes below require a valid access token.

### `GET /private/notifications/my`

Returns notifications for the authenticated user.

Response `200`:

```json
{
  "response": []
}
```

### `PATCH /private/notifications/action/read/:id`

Marks a notification as read.

Response `200`:

```json
{
  "response": true
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `403` | `ERR_NOTIFICATION_FORBIDDEN` | Notification belongs to another user |
| `404` | `ERR_NOTIFICATION_NOTFOUND` | Notification ID does not exist |

## Private Roadmap Admin Routes

All routes below require:

- a valid access token
- permission `DEVELOPER`

### `POST /private/roadmap/admin/task/add`

Creates a roadmap task.

Request body:

```json
{
  "title": "Add public profiles"
}
```

Validation:

- `title`: required, min `3`, max `200`

Response `201`:

```json
{
  "response": {}
}
```

### `PATCH /private/roadmap/admin/task/edit`

Updates a roadmap task.

Request body:

```json
{
  "id": 1,
  "title": "Add public profiles",
  "status": 1
}
```

Validation:

- `id`: required
- `title`: required, min `3`, max `200`
- `status`: roadmap enum value

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `404` | `ERR_TASK_NOTFOUND` | Task ID does not exist |

## Route Matrix

| Method | Path | Auth | Permission |
| --- | --- | --- | --- |
| `GET` | `/v1/api/public/ping` | no | no |
| `POST` | `/v1/api/public/auth/register` | no | no |
| `POST` | `/v1/api/public/auth/login` | no | no |
| `POST` | `/v1/api/public/auth/refresh` | no | no |
| `GET` | `/v1/api/public/roadmap/list` | no | no |
| `POST` | `/v1/api/public/storage/upload/sharex` | ShareX token | no role check |
| `DELETE` | `/v1/api/private/auth/logout` | yes | no |
| `GET` | `/v1/api/private/user/me` | yes | no |
| `GET` | `/v1/api/private/user/lookup/:id` | yes | `VIEW_OTHER_PROFILES` |
| `GET` | `/v1/api/private/user/shareX/generate` | yes | `UPLOAD_FILES` |
| `GET` | `/v1/api/private/user/admin/users` | yes | `MANAGE_USERS` |
| `PUT` | `/v1/api/private/user/admin/role/add` | yes | `MANAGE_USERS` |
| `DELETE` | `/v1/api/private/user/admin/role/delete` | yes | `MANAGE_USERS` |
| `PATCH` | `/v1/api/private/user/admin/storage-limit/update` | yes | `MANAGE_USERS` |
| `PATCH` | `/v1/api/private/user/admin/verify/:id` | yes | `MANAGE_USERS` |
| `DELETE` | `/v1/api/private/user/admin/shareX/reset/:id` | yes | `MANAGE_USERS` |
| `POST` | `/v1/api/private/storage/upload/init` | yes | `UPLOAD_FILES` |
| `POST` | `/v1/api/private/storage/upload/chunk` | yes | `UPLOAD_FILES` |
| `POST` | `/v1/api/private/storage/upload/complete` | yes | `UPLOAD_FILES` |
| `POST` | `/v1/api/private/storage/delete` | yes | `UPLOAD_FILES` |
| `GET` | `/v1/api/private/storage/my` | yes | `VIEW_OWN_FILES` |
| `GET` | `/v1/api/private/storage/list` | yes | `VIEW_OTHER_FILES` |
| `GET` | `/v1/api/private/storage/list/:id` | yes | `VIEW_OTHER_FILES` |
| `PUT` | `/v1/api/private/storage/album/put` | yes | `UPLOAD_FILES` |
| `DELETE` | `/v1/api/private/storage/album/delete` | yes | `UPLOAD_FILES` |
| `PATCH` | `/v1/api/private/storage/post/action/like/:id` | yes | `VIEW_OTHER_FILES` |
| `DELETE` | `/v1/api/private/storage/post/action/likeRemove/:id` | yes | `VIEW_OTHER_FILES` |
| `GET` | `/v1/api/private/storage/post/action/download/:id` | yes | `DOWNLOAD_OTHERS_FILES` |
| `POST` | `/v1/api/private/storage/post/action/addComment` | yes | `VIEW_OTHER_FILES` |
| `GET` | `/v1/api/private/storage/post/:id` | yes | `VIEW_OTHER_FILES` |
| `POST` | `/v1/api/private/album/create` | yes | `UPLOAD_FILES` |
| `DELETE` | `/v1/api/private/album/delete/:id` | yes | `UPLOAD_FILES` |
| `GET` | `/v1/api/private/album/lookup/:id` | yes | `VIEW_OWN_FILES` |
| `GET` | `/v1/api/private/album/lookupAll` | yes | `MANAGE_FILES` |
| `GET` | `/v1/api/private/roles/admin/all` | yes | `MANAGE_ROLES` |
| `POST` | `/v1/api/private/roles/admin/create` | yes | `MANAGE_ROLES` |
| `PATCH` | `/v1/api/private/roles/admin/edit` | yes | `MANAGE_ROLES` |
| `DELETE` | `/v1/api/private/roles/admin/delete` | yes | `MANAGE_ROLES` |
| `GET` | `/v1/api/private/notifications/my` | yes | no route-level permission |
| `PATCH` | `/v1/api/private/notifications/action/read/:id` | yes | no route-level permission |
| `POST` | `/v1/api/private/roadmap/admin/task/add` | yes | `DEVELOPER` |
| `PATCH` | `/v1/api/private/roadmap/admin/task/edit` | yes | `DEVELOPER` |
