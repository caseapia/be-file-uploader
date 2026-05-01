# Authentication and Sessions

## Model Overview

The API uses a two-token model backed by a persistent session row:

- Access token: short-lived JWT, sent on every private request.
- Refresh token: opaque random token, exchanged for a new token pair.
- Session row: DB record that ties a user, refresh-token hash, IP metadata, user-agent metadata, active state, and expiry together.

## Access Token

Current access-token properties:

| Property      | Value                            |
|---------------|----------------------------------|
| Algorithm     | `HS256`                          |
| Secret source | `JWT_SECRET`                     |
| Lifetime      | `15 minutes`                     |
| Claims        | `sub`, `sid`, `tv`, `iat`, `exp` |

Meaning of claims:

- `sub`: user ID
- `sid`: session ID
- `tv`: token version, currently hard-coded as `1`

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
   - `last_active_at`
5. An access token is generated with the session ID in `sid`.

Registration currently requires only `username` and `password`; it does not consume an invite code.

### Refresh

During refresh:

1. The incoming refresh token is hashed.
2. A session is searched by `refresh_hash`.
3. The service loads the related user.
4. A brand-new refresh token is generated.
5. A new JWT is generated.
6. The session hash, expiry, IP, and user-agent are updated.

Result:

- old refresh token becomes invalid
- new refresh token must be stored immediately
- existing access token should also be replaced immediately

### Logout

`DELETE /private/auth/logout` disables the current session:

- `is_active` is set to `false`
- `expires_at` is set to the current time
- IP and user-agent metadata are updated

After logout, the same access token should fail private middleware because the backing session is inactive.

## How Private Requests Are Validated

`auth.Middleware` runs on `/v1/api/private/*`.

Validation sequence:

1. Extract token from:
   - `auth_token` cookie, or
   - `Authorization` header
2. Parse JWT.
3. Load the user from DB.
4. Load the session from DB by JWT `sid`.
5. Reject if session is missing, inactive, or expired.
6. Store `user` and `session` in `ctx.Locals`.

## How to Send the Access Token

### Option 1: Authorization Header

```http
Authorization: Bearer <jwt>
```

Recommended for:

- mobile clients
- backend-to-backend calls
- Insomnia/Postman
- SPAs that keep token in memory

### Option 2: Cookie

```http
Cookie: auth_token=<jwt>
```

Recommended for:

- browser flows where your own backend manages cookies

Current code checks cookie first, then header.

## ShareX Token

ShareX upload does not use JWT auth. It uses a per-user API token:

1. Authenticated users with `UPLOAD_FILES` call `GET /private/user/shareX/generate`.
2. The returned token is stored in `users.sharex_token`.
3. ShareX sends the token as a multipart form field named `token` to `POST /public/storage/upload/sharex`.
4. Admins with `MANAGE_USERS` can reset another user's token with `DELETE /private/user/admin/shareX/reset/:id`.

ShareX upload returns:

```json
{
  "url": "https://cdn.example.com/images/..."
}
```

This response is intentionally not wrapped in `response`.

## `X-User-Agent`

The app stores `X-User-Agent` in session metadata during login, refresh, and logout.

Recommended pattern:

```http
X-User-Agent: webapp/0.4.0
```

Why it matters:

- useful for audit logs
- useful when inspecting session anomalies
- part of the data updated on refresh/logout

## Frontend Integration

### Recommended Browser Strategy

- Keep `access_token` in memory.
- Store `refresh_token` in the most secure storage your architecture allows.
- On `401` or expired-access-token flow, call `/public/auth/refresh`.
- Replace both tokens from the refresh response.
- Retry the original request once.
- Call `/private/auth/logout` when the user explicitly signs out.

### Example Fetch Wrapper

```ts
let accessToken = "";
let refreshToken = "";

async function apiFetch(path: string, init: RequestInit = {}) {
  const headers = new Headers(init.headers);
  headers.set("X-User-Agent", "frontend/1.0.0");

  if (accessToken) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  }

  const response = await fetch(`/v1/api${path}`, {
    ...init,
    headers,
  });

  if (response.status !== 401) {
    return response;
  }

  const refreshResponse = await fetch("/v1/api/public/auth/refresh", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User-Agent": "frontend/1.0.0",
    },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  if (!refreshResponse.ok) {
    accessToken = "";
    refreshToken = "";
    throw new Error("refresh failed");
  }

  const refreshBody = await refreshResponse.json();
  accessToken = refreshBody.response.access_token;
  refreshToken = refreshBody.response.refresh_token;

  headers.set("Authorization", `Bearer ${accessToken}`);
  return fetch(`/v1/api${path}`, { ...init, headers });
}
```

### Upload Integration

Private upload is a three-step flow:

1. `POST /private/storage/upload/init` with JSON metadata.
2. `POST /private/storage/upload/chunk` with `multipart/form-data`.
3. `POST /private/storage/upload/complete` with uploaded part metadata.

Example chunk body fields:

- `upload_id`
- `key`
- `part_number`
- `chunk`

Do not set JSON content type manually for the chunk request.

## Permission Gates

JWT validation only proves identity. Many routes then apply `RequirePermission(...)`.

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

- Invalid refresh-token cases currently collapse into `500 ERR_TOKEN_GENERATION` instead of a more specific `401/403` branch.
- Middleware permission denials often use Fiber's generic `Forbidden` response.
- Some malformed JWT errors may expose JWT library text instead of normalized project error codes.
- ShareX upload uses token auth and a custom unwrapped response shape.
