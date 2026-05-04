# Authentication and Sessions

## Model Overview

The API uses a two-token model backed by a persistent session row:

- Access token: JWT, required on private routes.
- Refresh token: opaque random token, exchanged for a new token pair.
- Session row: DB record that ties user, refresh-hash, IP/user-agent metadata, active state, and expiry.

## Access Token

Current access-token properties:

| Property      | Value |
|---------------|-------|
| Algorithm     | `HS256` |
| Secret source | `JWT_SECRET` (`models.JWTSecret`) |
| Lifetime      | `7 days` |
| Claims        | `sub`, `sid`, `tv`, `iat`, `exp` |

Meaning of claims:

- `sub`: user ID
- `sid`: session ID
- `tv`: token version (currently `1`)

## Refresh Token

Current refresh-token properties:

| Property | Value |
| --- | --- |
| Random bytes | `32` |
| Encoding | Base64 URL with padding |
| Typical length | `44` chars |
| DB storage | SHA-256 hash only |
| Session TTL refresh | yes, on successful refresh |

## Session Lifecycle

### Login / Register

During login and register:

1. A session UUID is generated.
2. A refresh token is generated.
3. The refresh token is hashed and stored in `sessions.refresh_hash`.
4. Session metadata is stored:
   - `user_id`
   - `ip_address`
   - `user_agent`
   - `expires_at`
   - `created_at`
5. An access token is generated with the session ID in `sid`.

Registration accepts only `username` and `password`; invite-code flow is not used.

### Refresh

During refresh:

1. The incoming refresh token is hashed.
2. Session is searched by `refresh_hash`.
3. New refresh token and JWT are generated.
4. Session hash, expiry, IP, user-agent, and geo string are updated.

Result:

- old refresh token becomes invalid
- new refresh token must be stored immediately
- access token must be replaced immediately

### Logout

`DELETE /private/auth/logout` disables the current session:

- `is_active` is set to `false`
- `expires_at` is set to current time
- IP and user-agent metadata are updated

## How Private Requests Are Validated

`auth.Middleware` runs on `/v1/api/private/*`.

Validation sequence:

1. Extract token from:
   - `Authorization` header (`Bearer <jwt>`), then
   - `access_token` cookie
2. Parse JWT and read session claims.
3. Load user by JWT `sub`.
4. Update user metadata (`last_ip`, `useragent`, `cf_ray_id`, `locale`, `last_seen`, `geo_string`).
5. Load session by JWT `sid`.
6. Reject when session is missing, inactive, or expired.
7. Put `user` and `session` into `ctx.Locals`.

Session expiry handling:

- expired sessions return `ERR_USER_SESSION_EXPIRED`
- sessions older than 2 days past expiry are deleted asynchronously

## How to Send the Access Token

### Option 1: Authorization Header

```http
Authorization: Bearer <jwt>
```

### Option 2: Cookie

```http
Cookie: access_token=<jwt>
```

Current middleware checks header first, then cookie.

## ShareX Token Auth

ShareX upload does not use JWT auth. It uses a per-user API token:

1. Authenticated users with `UPLOAD_FILES` call `GET /private/user/shareX/generate`.
2. Returned token is stored in `users.sharex_token`.
3. ShareX sends token as multipart form field `token` to `POST /public/storage/upload/sharex`.
4. Admins with `MANAGE_USERS` can reset another user's token via `DELETE /private/user/admin/shareX/reset/:id`.

ShareX upload success response is not wrapped:

```json
{
  "url": "https://cdn.example.com/images/..."
}
```

## `X-User-Agent`

`X-User-Agent` is used in:

- login session creation
- refresh session update
- logout session update
- auth middleware user-metadata update

## Permission Gates

JWT validation confirms identity. Many private routes then apply `RequirePermission(...)`.

Permissions used by current routes:

- `UPLOAD_FILES`
- `VIEW_OWN_FILES`
- `VIEW_OTHER_FILES`
- `DOWNLOAD_OTHERS_FILES`
- `VIEW_OTHER_PROFILES`
- `MANAGE_USERS`
- `MANAGE_FILES`
- `MANAGE_ROLES`
- `VIEW_PRIVATE_DATA`
- `DEVELOPER`

## Known Caveats

- `models.JWTSecret` is initialized at package load time. Contributors relying on `.env` loading in `main()` should verify startup environment behavior.
- `RefreshToken` currently ignores errors from new token generation (`newRefresh, _` and `newAccess, _`).
- JWT parsing errors can surface as raw library errors through the global error handler.
- Refresh lookup failures return `401 ERR_TOKEN_INVALID`.
- `ParseJWT` returns `404 ERR_USER_NOT_FOUND` when JWT claims reference a deleted user.
