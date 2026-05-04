# Errors

## Response Shapes

### Handled application errors

```json
{
  "error": "ERR_CODE",
  "code": 400
}
```

### Unhandled internal errors

```json
{
  "code": 500,
  "message": "raw error text"
}
```

### ShareX upload success exception

`POST /public/storage/upload/sharex` returns a direct URL object instead of the standard success wrapper:

```json
{
  "url": "https://cdn.example.com/images/..."
}
```

## Validation Errors

Validation is performed by `go-playground/validator` through `ParseAndValidate`.

When validation fails, the app returns:

- HTTP `400`
- `error` containing the validator message string

Example:

```json
{
  "error": "Key: 'Login.Username' Error:Field validation for 'Username' failed on the 'min' tag",
  "code": 400
}
```

## Error Catalog by Area

### Authentication routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_WRONG_CREDENTIALS` | Login username not found or password mismatch |
| `409` | `ERR_USER_ALREADY_EXISTS` | Register username already exists |
| `409` | `ERR_USER_REGISTER_HASHCREATION` | Password hash generation failed |
| `409` | `ERR_USER_UNKNOWN_CREATION_ERROR` | User creation transaction failed |
| `409` | `ERR_USER_ADDINROLE_FAILED` | Default role assignment failed |
| `500` | `ERR_USER_LOOKUP_AFTER_REGISTER` | Re-fetching created user failed |
| `500` | `ERR_SESSION_CREATION` | Session insert failed during login/register |
| `500` | `ERR_TOKEN_GENERATION` | Token generation failed during login/register |
| `401` | `ERR_TOKEN_INVALID` | Refresh token hash not found in sessions |
| `404` | `ERR_USER_NOTFOUND` | Logout user re-read failed |
| `404` | `ERR_SESSION_NOTFOUND` | Logout session does not belong to user |
| `404` | `ERR_SESSION_NOTACTIVE` | Logout called on already inactive session |
| `500` | `ERR_SESSION_UPDATE` | Logout session update failed |

### Access-token middleware

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `401` | `ERR_TOKEN_NOTFOUND` | No access token in header/cookie |
| `500` | `ERR_DATABASE_UPDATE` | User metadata update failed in middleware |
| `401` | `ERR_SESSION_NOTFOUND` | JWT `sid` missing in DB |
| `403` | `ERR_USER_SESSION_REVOKED` | Session is inactive |
| `403` | `ERR_USER_SESSION_EXPIRED` | Session expired |
| `404` | `ERR_USER_NOT_FOUND` | JWT user lookup failed after token parsing |

Additional note:

- JWT parsing can also return library-driven errors that are not wrapped in `fiber.Error`. Those return the unhandled format with `code: 500` and raw `message`.

### Permission middleware

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `401` | `ERR_UNAUTHORIZED` | No user in `ctx.Locals` for permission middleware |
| `403` | `ERR_NO_ACCESS` | User lacks required permission |

### User and user-admin routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_USER_NOTFOUND` | Target user does not exist |
| `404` | `ERR_ROLE_NOTFOUND...` | Role lookup failed during add/remove role |
| `403` | `ERR_ROLE_ISSUE_FORBIDDEN` | Protected role operation attempted by non-admin |
| `401` | `ERR_INVALID_TOKEN` | ShareX token lookup failed |

### Role admin routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_ROLE_NOT_FOUND` | Role lookup/edit/delete target missing |
| `409` | `ERR_ROLE_ALREADY_EXISTS` | Role name uniqueness conflict |

### Storage routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `400` | `ERR_INVALID_PARAMS` | Missing/invalid `upload_id`, `key`, or `part_number` |
| `400` | `ERR_CHUNK_MISSING` | Multipart field `chunk` missing |
| `400` | `ERR_IMAGE_MISSING` | ShareX multipart field `image` missing |
| `403` | `ERR_IMAGE_FORBIDDEN` | File mutation rejected by ownership/permission checks |
| `403` | `ERR_ALBUM_FORBIDDEN` | Album assignment blocked by ownership checks |
| `404` | `ERR_IMAGE_NOTFOUND` | File not found or hidden by visibility rules |
| `404` | `ERR_ALBUM_NOTFOUND` | Album not found in storage action flow |
| `413` | `ERR_IMAGE_TOO_LARGE` | Upload exceeds `4 GiB` limit |
| `413` / `404` | `ERR_QUOTA_EXCEEDED` | Quota exceeded from size precheck (`413`) or reserve step (`404`) |
| `500` | `ERR_OPEN_IMAGE` | Failed opening multipart file |
| `500` | `ERR_FILE_READING` | Failed reading multipart file |
| `500` | `ERR_MIMETYPE` | Unsupported MIME type |
| `500` | `ERR_UPLOAD_CHUNK` | Object part upload failed |
| `500` | `ERR_COMPLETE_MULTIPART` | Multipart completion failed |
| `500` | `ERR_IMAGE_UPLOAD` | File update failed in album/comment/download operations |
| `500` | `ERR_IMAGE_DELETE` | DB file delete failed in transaction |
| `500` | `ERR_DELETE_IMAGE` | Transaction failed during delete flow |

### Album routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ALBUM_NOTFOUND` | Album lookup denied or missing |
| `404` | `ALBUM_NOT_FOUND` | Album delete target missing |

### Notification routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_NOTIFICATION_NOTFOUND` | Notification ID does not exist |
| `403` | `ERR_NOTIFICATION_FORBIDDEN` | Notification belongs to another user |

### Roadmap routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_TASK_NOTFOUND` | Roadmap task ID does not exist |

## Known Inconsistencies Worth Documenting

- Refresh/session and token failures are not fully normalized across auth code paths.
- JWT parsing failures can surface as raw library text (`message`) instead of project error codes.
- Album not-found codes vary between `ALBUM_NOTFOUND` and `ALBUM_NOT_FOUND`.
- User-role lookup error currently concatenates `ERR_ROLE_NOTFOUND` with raw DB error text.
- `ERR_QUOTA_EXCEEDED` can be emitted with multiple status codes depending on the code path.
