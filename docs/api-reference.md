# API Reference

Base URL:

```text
http://localhost:{APP_PORT}/v1/api
```

## Common Rules

### Success Envelope

All successful JSON responses are wrapped in a top-level `response` field.

```json
{
  "response": {}
}
```

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

For uploads use `multipart/form-data` instead of JSON.

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
  "roles": [
    {
      "id": 1,
      "name": "member",
      "permissions": ["VIEW_OWN_FILES"],
      "is_system": true,
      "created_at": "2026-04-11T12:00:00Z",
      "created_by": 1,
      "color": "#00a86b"
    }
  ],
  "invite": {
    "id": 3,
    "code": "ABC12D",
    "creator": {
      "id": 1,
      "username": "admin"
    },
    "is_active": false,
    "user": {
      "id": 2,
      "username": "alice"
    }
  }
}
```

### Invite

```json
{
  "id": 3,
  "code": "ABC12D",
  "creator": {
    "id": 1,
    "username": "admin"
  },
  "is_active": true,
  "user": {
    "id": 0,
    "username": ""
  }
}
```

Notes:

- `created_by` and `used_by` are not serialized.
- Related `creator` and `user` are included by repository queries.

### Image

```json
{
  "id": 10,
  "url": "https://cdn.example.com/images/2/2026-04/img__ABC123DEF456image/png",
  "original_name": "photo.png",
  "mime_type": "image/png",
  "size": 245901,
  "uploader": {
    "id": 2,
    "username": "alice"
  }
}
```

Notes:

- Storage key `r2_key` is intentionally hidden from JSON responses.
- URL is a public URL assembled from `r2_public_url` and the generated object key.

## Public Routes

### `POST /public/auth/register`

Creates a new account from an invite code and immediately issues an access/refresh token pair.

Request body:

```json
{
  "username": "alice",
  "password": "secret12",
  "invite_code": "ABC12D"
}
```

Validation:

- `username`: required, min `3`, max `20`
- `password`: required, min `6`, max `32`
- `invite_code`: required, exact length `6`

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
- Finds invite by code.
- Rejects inactive invites.
- Creates the user.
- Marks invite as used.
- Assigns the new user to role ID `1`.
- Creates a session and issues tokens.

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | validator message | Bad body or failed field validation |
| `409` | `ERR_USER_ALREADY_EXISTS` | Username already taken |
| `409` | `ERR_INVITE_NOT_FOUND` | Invite code does not exist |
| `409` | `ERR_INVITE_ALREADY_USED` | Invite exists but is inactive |
| `409` | `ERR_USER_REGISTER_HASHCREATION` | Password hash creation failed |
| `409` | `ERR_USER_UNKNOWN_CREATION_ERROR` | User creation failed inside transaction |
| `409` | `ERR_INVITE_MARK_ERR` | Invite could not be marked as used |
| `409` | `ERR_USER_ADDINROLE_FAILED` | Default role assignment failed |
| `500` | `ERR_USER_LOOKUP_AFTER_REGISTER` | User re-read failed after creation |
| `500` | `ERR_TOKEN_GENERATION` | Token generation/session creation failed |

Frontend notes:

- The response already contains `user`, so you can hydrate the current-user store immediately after register.
- The invite is consumed on success.

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

Frontend notes:

- Treat login `404` as invalid credentials, not as a missing route.
- Every login creates a fresh session. Existing sessions are not invalidated automatically.

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

Implementation details:

- Refresh tokens are stored hashed in DB.
- Successful refresh updates session expiry, IP, user-agent, and refresh hash.
- Old refresh token becomes invalid immediately.

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | validator message | Bad body or failed field validation |
| `500` | `ERR_TOKEN_GENERATION` | Current implementation also uses this for invalid/missing/expired refresh-session lookups |
| `500` | `ERR_USER_NOTFOUND` | Session user no longer exists |

Frontend notes:

- Always overwrite both tokens after refresh.
- Current implementation does not distinguish invalid refresh token from internal token-generation failure with separate status codes. Frontend should treat refresh `500 ERR_TOKEN_GENERATION` as "session recovery failed, force re-login".

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

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `401` | `ERR_TOKEN_NOTFOUND` | No access token supplied |
| `401` | `ERR_SESSION_NOTFOUND` | Token session missing in DB |
| `403` | `ERR_USER_SESSION_REVOKED` | Session disabled |
| `403` | `ERR_USER_SESSION_EXPIRED` | Session expired |

### `GET /private/user/lookup/:id`

Looks up a user by numeric ID.

Path params:

| Param | Type | Example |
| --- | --- | --- |
| `id` | integer | `42` |

Intended behavior:

- allow self-lookup
- allow privileged lookup when user has `VIEW_OTHER_PROFILES`

Current implementation caveat:

- The service currently returns `403 ERR_NO_ACCESS` for most non-self lookups and can also reject privileged callers because of the condition in `internal/service/user/user.go`.
- If you are building frontend around this route, assume self-lookup is safe via `/user/me` and treat `/user/lookup/:id` as unstable until the service logic is corrected.

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `403` | `ERR_NO_ACCESS` | Current permission logic rejected lookup |
| `401` | `ERR_TOKEN_NOTFOUND` | No access token supplied |
| `401` | `ERR_SESSION_NOTFOUND` | Token session missing in DB |
| `403` | `ERR_USER_SESSION_REVOKED` | Session disabled |
| `403` | `ERR_USER_SESSION_EXPIRED` | Session expired |

## Private Invite Admin Routes

All routes below require:

- a valid access token
- permission `MANAGE_USERS`

If the permission check fails, the middleware returns Fiber `403 Forbidden`.

### `GET /private/invite/admin/list`

Returns all invites.

Response `200`:

```json
{
  "response": [
    {
      "id": 1,
      "code": "ABC12D",
      "creator": {
        "id": 1,
        "username": "admin"
      },
      "is_active": true,
      "user": {
        "id": 0,
        "username": ""
      }
    }
  ]
}
```

Frontend notes:

- The Insomnia collection auto-saves the first active invite ID into `invite_id`.
- Useful for admin panels that need invite lifecycle visibility.

### `POST /private/invite/admin/create`

Creates a new invite code.

Request body: empty

Response `201`:

```json
{
  "response": "ABC12D"
}
```

Behavior:

- Code length is `6`
- Character set is uppercase `A-Z` plus digits `0-9`

Frontend notes:

- The Insomnia collection auto-saves the returned code into `invite_code`.
- A frontend admin screen can expose this route as "Generate invite" and display the resulting code immediately.

### `DELETE /private/invite/admin/revoke`

Revokes an invite by ID.

Request body:

```json
{
  "invite_id": 1
}
```

Validation:

- `invite_id`: required, integer, minimum `1`

Response `200`:

```json
{
  "response": {}
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | validator message | Bad body or failed field validation |
| `404` | `ERR_INVITE_NOTFOUND` | Invite ID does not exist |

## Private Storage Routes

All routes below require a valid access token.

### `POST /private/storage/upload`

Uploads one image file.

Request:

- content type: `multipart/form-data`
- form field name: `image`

Accepted MIME types:

- `image/jpeg`
- `image/png`
- `image/webp`
- `image/gif`

Limit:

- `10 MiB`

Response `201`:

```json
{
  "response": {
    "id": 10,
    "url": "https://cdn.example.com/images/2/2026-04/img__ABC123DEF456image/png",
    "original_name": "photo.png",
    "mime_type": "image/png",
    "size": 245901,
    "uploader": {
      "id": 2,
      "username": "alice"
    }
  }
}
```

Implementation details:

- The route itself does not apply a permission middleware.
- Object key pattern is:

```text
images/{userID}/{YYYY-MM}/{generatedID+ext}
```

Current caveat:

- The code stores the MIME type string itself as the suffix when building the key, not a file extension. That affects the generated URL shape.

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | `ERR_IMAGE_MISSING` | No multipart file named `image` |
| `413` | `ERR_IMAGE_TOO_LARGE` | File exceeds `10 MiB` |
| `415` | `ERR_IMAGE_UNSUPPORTED_TYPE` | Unsupported MIME type |
| `500` | `ERR_OPEN_IMAGE` | Multipart file open failed |
| `500` | `ERR_FILE_READING` | File read failed |
| `500` | `ERR_IMAGE_ID_GEN` | Generated image ID failed |
| `500` | `ERR_IMAGE_UPLOAD` | Object storage upload failed |

Frontend notes:

- Submit as `FormData`, not JSON.
- The Insomnia collection auto-saves `image_id` and `image_url` from successful uploads.

### `POST /private/storage/delete`

Deletes an image by ID.

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
  "response": "OK"
}
```

Current implementation caveat:

- The permission check is stricter than the endpoint name suggests.
- Because of the condition in `internal/service/image/image.go`, the request succeeds only when the requester is both:
  - the uploader of the image
  - a user who has `MANAGE_FILES`
- A plain uploader without `MANAGE_FILES` currently receives `403 ERR_IMAGE_FORBIDDEN`.

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `400` | validator message | Bad body or failed field validation |
| `403` | `ERR_IMAGE_FORBIDDEN` | Current permission logic rejected deletion |
| `404` | `ERR_IMAGE_NOTFOUND` | Image ID does not exist |
| `500` | `ERR_DELETE_IMAGE` | DB delete or object delete failed |

Frontend notes:

- Do not expose delete in UI unless the user is known to have the corresponding admin capability.

### `GET /private/storage/my`

Returns images uploaded by the authenticated user.

Required permission:

- `VIEW_OWN_FILES`

Response `200`:

```json
{
  "response": []
}
```

Possible handled errors:

| HTTP | `error` | Meaning |
| --- | --- | --- |
| `403` | `Forbidden` | User lacks `VIEW_OWN_FILES` |

Frontend notes:

- Suitable for "My uploads" pages.
- The Insomnia collection auto-saves the first returned image into `image_id` and `image_url`.

## Private Image Admin Routes

All routes below require:

- a valid access token
- permission `VIEW_OTHER_FILES`

### `GET /private/image/admin/list`

Returns all stored images ordered by ID descending.

Response `200`:

```json
{
  "response": []
}
```

### `GET /private/image/admin/list/:id`

Returns images uploaded by the user with ID `:id`.

Path params:

| Param | Type | Example |
| --- | --- | --- |
| `id` | integer | `2` |

Response `200`:

```json
{
  "response": []
}
```

Frontend notes:

- Useful for moderator/admin dashboards.
- The Insomnia collection includes both all-images and per-user lookup requests.

## Route Matrix

| Method | Path | Auth | Permission |
| --- | --- | --- | --- |
| `POST` | `/v1/api/public/auth/register` | no | no |
| `POST` | `/v1/api/public/auth/login` | no | no |
| `POST` | `/v1/api/public/auth/refresh` | no | no |
| `GET` | `/v1/api/private/user/me` | yes | no |
| `GET` | `/v1/api/private/user/lookup/:id` | yes | current logic is restrictive |
| `GET` | `/v1/api/private/invite/admin/list` | yes | `MANAGE_USERS` |
| `POST` | `/v1/api/private/invite/admin/create` | yes | `MANAGE_USERS` |
| `DELETE` | `/v1/api/private/invite/admin/revoke` | yes | `MANAGE_USERS` |
| `POST` | `/v1/api/private/storage/upload` | yes | no route-level permission |
| `POST` | `/v1/api/private/storage/delete` | yes | current logic effectively requires uploader + `MANAGE_FILES` |
| `GET` | `/v1/api/private/storage/my` | yes | `VIEW_OWN_FILES` |
| `GET` | `/v1/api/private/image/admin/list` | yes | `VIEW_OTHER_FILES` |
| `GET` | `/v1/api/private/image/admin/list/:id` | yes | `VIEW_OTHER_FILES` |
