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
- `error` contains the validator message string

Example:

```json
{
  "error": "Key: 'Login.Username' Error:Field validation for 'Username' failed on the 'min' tag",
  "code": 400
}
```

Frontend note:

- Validation messages are not normalized into a structured field-error object.
- If you need field-level UI errors, parse or map these messages on the frontend, or improve the backend later.

## Error Catalog by Area

### Authentication

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_WRONG_CREDENTIALS` | Login username not found or password mismatch |
| `409` | `ERR_USER_ALREADY_EXISTS` | Register username already exists |
| `409` | `ERR_USER_REGISTER_HASHCREATION` | Password hash generation failed |
| `409` | `ERR_USER_UNKNOWN_CREATION_ERROR` | User creation transaction failed |
| `409` | `ERR_USER_ADDINROLE_FAILED` | Default role assignment failed |
| `500` | `ERR_USER_LOOKUP_AFTER_REGISTER` | Re-fetching the created user failed |
| `500` | `ERR_SESSION_CREATION` | Session insert failed during login/register |
| `500` | `ERR_TOKEN_GENERATION` | Access/refresh generation failed, and currently also invalid refresh-session cases |
| `500` | `ERR_USER_NOTFOUND` | Session user missing during refresh |
| `404` | `ERR_SESSION_NOTFOUND` | Logout session does not belong to the current user |
| `404` | `ERR_SESSION_NOTACTIVE` | Logout session is already inactive |
| `500` | `ERR_SESSION_UPDATE` | Logout session update failed |

### Access-token middleware

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `401` | `ERR_TOKEN_NOTFOUND` | No cookie/header token found |
| `401` | `ERR_SESSION_NOTFOUND` | JWT session missing in DB |
| `403` | `ERR_USER_SESSION_REVOKED` | Session inactive |
| `403` | `ERR_USER_SESSION_EXPIRED` | Session expired |
| `404` | `ERR_USER_NOT_FOUND` | JWT parsed but referenced user missing |
| `500` | `ERR_UNEXPECTED_SIGNING_METHOD` | JWT signing algorithm was not HMAC |
| `500` | `ERR_INVALID_TOKEN` | JWT claims were invalid |
| `500` | `ERR_TOKEN_EXPIRED` | Access token expiry check failed |

Additional note:

- JWT parsing can also surface library-driven token errors through `fiber.NewError`, so some malformed/expired JWT responses may expose library text instead of a stable project-specific code.

### User routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_USER_NOTFOUND` | User lookup/admin target does not exist |
| `404` | `ERR_ROLE_NOTFOUND...` | Role lookup during user-role assignment failed |
| `403` | `ERR_ROLE_ISSUE_FORBIDDEN` | Non-admin attempted to assign/remove protected role |
| `401` | `ERR_INVALID_TOKEN` | ShareX token lookup failed |

### Role admin routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_ROLE_NOT_FOUND` | Role lookup/delete target does not exist |
| `409` | `ERR_ROLE_ALREADY_EXISTS` | Role name uniqueness conflict |
| `403` | `Forbidden` | Missing `MANAGE_ROLES` permission from Fiber middleware |

### Storage routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `400` | `ERR_INVALID_PARAMS` | Multipart chunk upload params are missing or invalid |
| `400` | `ERR_CHUNK_MISSING` | Multipart field `chunk` missing |
| `400` | `ERR_IMAGE_MISSING` | ShareX multipart field `image` missing |
| `403` | `ERR_IMAGE_FORBIDDEN` | Delete or album assignment rejected by ownership/permission logic |
| `404` | `ERR_IMAGE_NOTFOUND` | File lookup/delete target missing or hidden as private |
| `413` | `ERR_IMAGE_TOO_LARGE` | Upload exceeds `4 GiB` |
| `413` | `ERR_QUOTA_EXCEEDED` | Upload exceeds the user's remaining storage quota |
| `500` | `ERR_OPEN_IMAGE` | Failed opening multipart file |
| `500` | `ERR_FILE_READING` | Failed reading multipart file |
| `500` | `ERR_MIMETYPE` | MIME type is unsupported |
| `500` | `ERR_UPLOAD_CHUNK` | Object part upload failed |
| `500` | `ERR_COMPLETE_MULTIPART` | Multipart completion failed |
| `500` | `ERR_IMAGE_UPLOAD` | File update failed in album/comment/download flows |
| `500` | `ERR_IMAGE_DELETE` | DB delete failed inside delete transaction |
| `500` | `ERR_DELETE_IMAGE` | DB/object delete failure |
| `403` | `Forbidden` | Missing storage permission from Fiber middleware |

### Album routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ALBUM_NOTFOUND` | Album lookup denied or missing |
| `404` | `ALBUM_NOT_FOUND` | Album delete target missing |
| `403` | `ERR_ALBUM_FORBIDDEN` | Attempted to put an image into another user's album |
| `403` | `Forbidden` | Missing album route permission from Fiber middleware |

### Notification routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_NOTIFICATION_NOTFOUND` | Notification ID does not exist |
| `403` | `ERR_NOTIFICATION_FORBIDDEN` | Notification belongs to another user |

### Roadmap routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_TASK_NOTFOUND` | Roadmap task ID does not exist |
| `403` | `Forbidden` | Missing `DEVELOPER` permission from Fiber middleware |

## Frontend Handling Guidance

Recommended high-level handling:

- `400`: validation or malformed request, show actionable form feedback
- `401`: missing or no longer valid access context, attempt refresh once
- `403`: authenticated but not allowed, hide/disable capability and show permission/session-expired message
- `404`: often business-level not-found in this API, not only missing route
- `409`: conflict during registration or role creation flows
- `413`: uploaded file too large or quota exceeded
- `500`: generic server failure or currently overloaded refresh/upload failure category

## Suggested Frontend Mapping

```ts
export function normalizeApiError(body: any) {
  return {
    status: typeof body?.code === "number" ? body.code : 500,
    code: body?.error ?? body?.message ?? "UNKNOWN_ERROR",
  };
}
```

```ts
switch (error.code) {
  case "ERR_WRONG_CREDENTIALS":
    return "Invalid username or password.";
  case "ERR_USER_ALREADY_EXISTS":
    return "Username is already taken.";
  case "ERR_IMAGE_TOO_LARGE":
    return "File is too large. Maximum size is 4 GiB.";
  case "ERR_QUOTA_EXCEEDED":
    return "You do not have enough storage remaining.";
  case "ERR_MIMETYPE":
    return "Unsupported file format.";
  case "ERR_TOKEN_NOTFOUND":
  case "ERR_SESSION_NOTFOUND":
  case "ERR_USER_SESSION_REVOKED":
  case "ERR_USER_SESSION_EXPIRED":
    return "Your session is no longer valid. Please sign in again.";
  default:
    return "Unexpected server error.";
}
```

## Known Inconsistencies Worth Documenting

- Refresh failure paths are not yet cleanly classified.
- Middleware permission denials often use Fiber's generic `Forbidden` message rather than a project-specific error code.
- Some malformed JWT errors may leak library messages instead of normalized API codes.
- Album not-found codes vary between `ALBUM_NOTFOUND` and `ALBUM_NOT_FOUND`.
- One user-role error concatenates `ERR_ROLE_NOTFOUND` with the underlying error text.
