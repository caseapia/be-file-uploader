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

## Validation Errors

Validation is performed by `go-playground/validator` through `ParseAndValidate`.

When validation fails, the app returns:

- HTTP `400`
- `error` contains the validator message string

Example:

```json
{
  "error": "Key: 'Register.Username' Error:Field validation for 'Username' failed on the 'min' tag",
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
| `409` | `ERR_INVITE_NOT_FOUND` | Register invite code not found |
| `409` | `ERR_INVITE_ALREADY_USED` | Register invite inactive |
| `409` | `ERR_USER_REGISTER_HASHCREATION` | Password hash generation failed |
| `409` | `ERR_USER_UNKNOWN_CREATION_ERROR` | User creation transaction failed |
| `409` | `ERR_INVITE_MARK_ERR` | Invite update failed after registration |
| `409` | `ERR_USER_ADDINROLE_FAILED` | Default role assignment failed |
| `500` | `ERR_USER_LOOKUP_AFTER_REGISTER` | Re-fetching the created user failed |
| `500` | `ERR_SESSION_CREATION` | Session insert failed during login |
| `500` | `ERR_TOKEN_GENERATION` | Access/refresh generation failed, and currently also some invalid refresh-session cases |
| `500` | `ERR_USER_NOTFOUND` | Session user missing during refresh |

### Access-token middleware

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `401` | `ERR_TOKEN_NOTFOUND` | No cookie/header token found |
| `401` | `ERR_SESSION_NOTFOUND` | JWT session missing in DB |
| `403` | `ERR_USER_SESSION_REVOKED` | Session inactive |
| `403` | `ERR_USER_SESSION_EXPIRED` | Session expired |
| `404` | `ERR_USER_NOT_FOUND` | JWT parsed but referenced user missing |

Additional note:

- JWT parsing can also surface library-driven token errors through `fiber.NewError`, so some malformed/expired JWT responses may expose library text instead of a stable project-specific code.

### User routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `403` | `ERR_NO_ACCESS` | Current lookup permission logic denied access |

### Invite admin routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `404` | `ERR_INVITE_NOTFOUND` | Revoke invite by unknown ID |
| `403` | `Forbidden` | Missing `MANAGE_USERS` permission from Fiber middleware |

### Storage routes

| HTTP | `error` | Where it comes from |
| --- | --- | --- |
| `400` | `ERR_IMAGE_MISSING` | Multipart field `image` missing |
| `403` | `ERR_IMAGE_FORBIDDEN` | Delete rejected by current ownership/permission logic |
| `404` | `ERR_IMAGE_NOTFOUND` | Delete target image missing |
| `413` | `ERR_IMAGE_TOO_LARGE` | Upload exceeds `10 MiB` |
| `415` | `ERR_IMAGE_UNSUPPORTED_TYPE` | Unsupported MIME |
| `500` | `ERR_OPEN_IMAGE` | Failed opening multipart file |
| `500` | `ERR_FILE_READING` | Failed reading multipart file |
| `500` | `ERR_IMAGE_ID_GEN` | Image ID generation failed |
| `500` | `ERR_IMAGE_UPLOAD` | Object upload failed |
| `500` | `ERR_DELETE_IMAGE` | DB/object delete failure |
| `403` | `Forbidden` | Missing `VIEW_OWN_FILES` or `VIEW_OTHER_FILES` permission on protected list routes |

## Frontend Handling Guidance

Recommended high-level handling:

- `400`: validation or malformed request, show actionable form feedback
- `401`: missing or no longer valid access context, attempt refresh once
- `403`: authenticated but not allowed, hide/disable capability and show permission/session-expired message
- `404`: often business-level not-found in this API, not only missing route
- `409`: conflict during registration/invite flows
- `413`: uploaded file too large
- `415`: unsupported file type
- `500`: generic server failure or currently overloaded refresh failure category

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
  case "ERR_INVITE_NOT_FOUND":
    return "Invite code was not found.";
  case "ERR_INVITE_ALREADY_USED":
    return "Invite code has already been used or revoked.";
  case "ERR_IMAGE_TOO_LARGE":
    return "Image is too large. Maximum size is 10 MiB.";
  case "ERR_IMAGE_UNSUPPORTED_TYPE":
    return "Unsupported image format.";
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
