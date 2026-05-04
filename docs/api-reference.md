# API Reference

Base URL:

```text
http://localhost:{APP_PORT}/v1/api
```

## Common Rules

### Success Envelope

Most successful JSON responses are wrapped in top-level `response`.

```json
{
  "response": {}
}
```

Exception:

- `POST /public/storage/upload/sharex` returns `{ "url": "..." }` directly.

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
  - `Cookie: access_token=<jwt>`

### Frontend Headers

Common headers used by current handlers/middleware:

```http
Content-Type: application/json
X-User-Agent: your-client-name/version
```

For multipart routes (`/upload/chunk`, `/upload/sharex`) use `multipart/form-data`.

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
  "albums": [],
  "last_seen": "2026-04-11T12:00:00Z",
  "geolocation": {
    "country_code": "US",
    "country": "United States",
    "city": "New York"
  }
}
```

Notes:

- `private` and `geolocation` are included only when requester has `VIEW_PRIVATE_DATA`.
- Password and sensitive fields are never serialized directly.

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

- `r2_key` and `uploaded_by` are hidden from JSON.
- Depending on current visibility logic, `url` may be blank even when metadata is returned.

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

Response `200`:

```json
{
  "response": "pong"
}
```

### `POST /public/auth/register`

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

### `POST /public/auth/login`

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

### `POST /public/auth/refresh`

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

### `GET /public/roadmap/list`

Response `200`:

```json
{
  "response": {
    "roadmap": []
  }
}
```

### `POST /public/storage/upload/sharex`

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

## Private Auth Routes

All routes below require a valid access token.

### `DELETE /private/auth/logout`

Response `200`:

```json
{
  "response": "OK"
}
```

## Private User Routes

All routes below require a valid access token.

### `GET /private/user/me`

Response `200`:

```json
{
  "response": {}
}
```

### `GET /private/user/lookup/:id`

Required permission:

- `VIEW_OTHER_PROFILES`

Response `200`:

```json
{
  "response": {}
}
```

### `GET /private/user/shareX/generate`

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

- valid access token
- `MANAGE_USERS`

### `GET /private/user/admin/users`

Response `200`:

```json
{
  "response": []
}
```

### `PUT /private/user/admin/role/add`

Request body:

```json
{
  "user": 2,
  "role": 3
}
```

### `DELETE /private/user/admin/role/delete`

Request body:

```json
{
  "user": 2,
  "role": 3
}
```

### `PATCH /private/user/admin/storage-limit/update`

Request body:

```json
{
  "user": 2,
  "limit": 1073741824
}
```

### `PATCH /private/user/admin/verify/:id`

Path param `id`: integer.

### `DELETE /private/user/admin/shareX/reset/:id`

Path param `id`: integer.

## Private Storage Upload Routes

All routes below require:

- valid access token
- `UPLOAD_FILES`

### `POST /private/storage/upload/init`

Request body:

```json
{
  "original_name": "photo.png",
  "mime_type": "image/png",
  "size": 245901,
  "is_private": false
}
```

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

Limits:

- `4 GiB` file limit
- user quota limit (`upload_limit - used_storage`)

Response `200`:

```json
{
  "response": {
    "upload_id": "s3-upload-id",
    "key": "images/dev/2/2026-04/file__ABCD1234.png"
  }
}
```

### `POST /private/storage/upload/chunk`

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

### `POST /private/storage/upload/complete`

Request body:

```json
{
  "upload_id": "s3-upload-id",
  "key": "images/dev/2/2026-04/file__ABCD1234.png",
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

## Private Storage File Routes

### `POST /private/storage/action/delete`

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "image_id": 10
}
```

Response `200`:

```json
{
  "response": {
    "used_storage": 12345,
    "status": "OK"
  }
}
```

### `GET /private/storage/my`

Required permission:

- `VIEW_OWN_FILES`

Response `200`:

```json
{
  "response": []
}
```

### `GET /private/storage/list`

Required permission:

- `VIEW_OTHER_FILES`

Response `200`:

```json
{
  "response": []
}
```

### `GET /private/storage/list/:id`

Required permission:

- `VIEW_OTHER_FILES`

Returns public files for the target user ID.

### `GET /private/storage/post/:id`

Required permission:

- `VIEW_OTHER_FILES`

Response `200`:

```json
{
  "response": {}
}
```

## Private Storage Action Routes

### `PUT /private/storage/action/album/put`

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "image_id": 10,
  "album_id": 1
}
```

### `DELETE /private/storage/action/album/delete`

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "image_id": 10
}
```

### `PATCH /private/storage/action/like/:id`

Required permission:

- `VIEW_OTHER_FILES`

Response `200`:

```json
{
  "response": true
}
```

### `DELETE /private/storage/action/likeRemove/:id`

Required permission:

- `VIEW_OTHER_FILES`

Response `200`:

```json
{
  "response": true
}
```

### `GET /private/storage/action/download/:id`

Required permission:

- `DOWNLOAD_OTHERS_FILES`

Response `201`:

```json
{
  "response": "https://cdn.example.com/images/..."
}
```

### `POST /private/storage/action/addComment`

Required permission:

- `VIEW_OTHER_FILES`

Request body:

```json
{
  "post_id": 10,
  "content": "Looks good"
}
```

Response `201`:

```json
{
  "response": {}
}
```

### `PUT /private/storage/action/access/grant`

Request body:

```json
{
  "file_id": 10,
  "user_id": 1
}
```

### `DELETE /private/storage/action/access/remove`

Request body:

```json
{
  "file_id": 10,
  "user_id": 1
}
```

## Private Album Routes

All routes below require a valid access token.

### `POST /private/album/action/create`

Required permission:

- `UPLOAD_FILES`

Request body:

```json
{
  "album_name": "Screenshots",
  "is_private": false
}
```

Response `201`:

```json
{
  "response": {}
}
```

### `DELETE /private/album/action/delete/:id`

Required permission:

- `UPLOAD_FILES`

Response `200`:

```json
{
  "response": true
}
```

### `GET /private/album/lookup/:id`

Required permission:

- `VIEW_OWN_FILES`

### `GET /private/album/lookupAll`

Required permission:

- `MANAGE_FILES`

## Private Roles Admin Routes

All routes below require:

- valid access token
- `MANAGE_ROLES`

### `GET /private/roles/admin/all`

Response `200` with role array.

### `POST /private/roles/admin/create`

Request body:

```json
{
  "name": "moderator",
  "color": "#3366ff",
  "is_system": false,
  "permissions": ["VIEW_OTHER_FILES", "MANAGE_FILES"]
}
```

Response `201`.

### `PATCH /private/roles/admin/edit`

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

Response `201`.

### `DELETE /private/roles/admin/delete`

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

Response `200`:

```json
{
  "response": []
}
```

### `PATCH /private/notifications/action/read/:id`

Response `200`:

```json
{
  "response": true
}
```

## Private Roadmap Admin Routes

All routes below require:

- valid access token
- `DEVELOPER`

### `POST /private/roadmap/admin/task/add`

Request body:

```json
{
  "title": "Add public profiles"
}
```

Response `201`.

### `PATCH /private/roadmap/admin/task/edit`

Request body:

```json
{
  "id": 1,
  "title": "Add public profiles",
  "status": 1
}
```

Response `200`.

## Route Matrix

| Method   | Path                                              | Auth         | Permission                |
|----------|---------------------------------------------------|--------------|---------------------------|
| `GET`    | `/v1/api/public/ping`                             | no           | no                        |
| `POST`   | `/v1/api/public/auth/register`                    | no           | no                        |
| `POST`   | `/v1/api/public/auth/login`                       | no           | no                        |
| `POST`   | `/v1/api/public/auth/refresh`                     | no           | no                        |
| `GET`    | `/v1/api/public/roadmap/list`                     | no           | no                        |
| `POST`   | `/v1/api/public/storage/upload/sharex`            | ShareX token | no role check             |
| `DELETE` | `/v1/api/private/auth/logout`                     | yes          | no                        |
| `GET`    | `/v1/api/private/user/me`                         | yes          | no                        |
| `GET`    | `/v1/api/private/user/lookup/:id`                 | yes          | `VIEW_OTHER_PROFILES`     |
| `GET`    | `/v1/api/private/user/shareX/generate`            | yes          | `UPLOAD_FILES`            |
| `GET`    | `/v1/api/private/user/admin/users`                | yes          | `MANAGE_USERS`            |
| `PUT`    | `/v1/api/private/user/admin/role/add`             | yes          | `MANAGE_USERS`            |
| `DELETE` | `/v1/api/private/user/admin/role/delete`          | yes          | `MANAGE_USERS`            |
| `PATCH`  | `/v1/api/private/user/admin/storage-limit/update` | yes          | `MANAGE_USERS`            |
| `PATCH`  | `/v1/api/private/user/admin/verify/:id`           | yes          | `MANAGE_USERS`            |
| `DELETE` | `/v1/api/private/user/admin/shareX/reset/:id`     | yes          | `MANAGE_USERS`            |
| `POST`   | `/v1/api/private/storage/upload/init`             | yes          | `UPLOAD_FILES`            |
| `POST`   | `/v1/api/private/storage/upload/chunk`            | yes          | `UPLOAD_FILES`            |
| `POST`   | `/v1/api/private/storage/upload/complete`         | yes          | `UPLOAD_FILES`            |
| `DELETE` | `/v1/api/private/storage/action/delete`           | yes          | `UPLOAD_FILES`            |
| `GET`    | `/v1/api/private/storage/my`                      | yes          | `VIEW_OWN_FILES`          |
| `GET`    | `/v1/api/private/storage/list`                    | yes          | `VIEW_OTHER_FILES`        |
| `GET`    | `/v1/api/private/storage/list/:id`                | yes          | `VIEW_OTHER_FILES`        |
| `PUT`    | `/v1/api/private/storage/action/album/put`        | yes          | `UPLOAD_FILES`            |
| `DELETE` | `/v1/api/private/storage/action/album/delete`     | yes          | `UPLOAD_FILES`            |
| `PATCH`  | `/v1/api/private/storage/action/like/:id`         | yes          | `VIEW_OTHER_FILES`        |
| `DELETE` | `/v1/api/private/storage/action/likeRemove/:id`   | yes          | `VIEW_OTHER_FILES`        |
| `GET`    | `/v1/api/private/storage/action/download/:id`     | yes          | `DOWNLOAD_OTHERS_FILES`   |
| `POST`   | `/v1/api/private/storage/action/addComment`       | yes          | `VIEW_OTHER_FILES`        |
| `PUT`    | `/v1/api/private/storage/action/access/grant`     | yes          | no                        |
| `DELETE` | `v1/api/private/storage/action/access/remove`     | yes          | no                        |
| `GET`    | `/v1/api/private/storage/post/:id`                | yes          | `VIEW_OTHER_FILES`        |
| `POST`   | `/v1/api/private/album/action/create`             | yes          | `UPLOAD_FILES`            |
| `DELETE` | `/v1/api/private/album/action/delete/:id`         | yes          | `UPLOAD_FILES`            |
| `GET`    | `/v1/api/private/album/lookup/:id`                | yes          | `VIEW_OWN_FILES`          |
| `GET`    | `/v1/api/private/album/lookupAll`                 | yes          | `MANAGE_FILES`            |
| `GET`    | `/v1/api/private/roles/admin/all`                 | yes          | `MANAGE_ROLES`            |
| `POST`   | `/v1/api/private/roles/admin/create`              | yes          | `MANAGE_ROLES`            |
| `PATCH`  | `/v1/api/private/roles/admin/edit`                | yes          | `MANAGE_ROLES`            |
| `DELETE` | `/v1/api/private/roles/admin/delete`              | yes          | `MANAGE_ROLES`            |
| `GET`    | `/v1/api/private/notifications/my`                | yes          | no route-level permission |
| `PATCH`  | `/v1/api/private/notifications/action/read/:id`   | yes          | no route-level permission |
| `POST`   | `/v1/api/private/roadmap/admin/task/add`          | yes          | `DEVELOPER`               |
| `PATCH`  | `/v1/api/private/roadmap/admin/task/edit`         | yes          | `DEVELOPER`               |

## Implementation Caveats

- Several handlers parse `:id` via `strconv.Atoi` and ignore parse errors.
- Current file URL visibility logic can blank `url` on endpoints that still return file metadata.
- `/private/storage/my` depends on current query/schema assumptions around grant alias `fg`.
