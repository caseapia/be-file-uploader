# Getting Started

## Stack

- Go `1.25.0` according to [`go.mod`](../go.mod)
- Fiber v3
- MySQL + Bun ORM
- S3-compatible object storage (current adapter targets Yandex Object Storage)
- IP2Location local database for geo metadata (`data/IP2LOCATION-LITE-DB11.IPV6.BIN`)
- Redis dependency is present in codebase, but runtime Redis initialization is currently disabled

## Prerequisites

Before starting the API locally you need:

- A MySQL database reachable from the app.
- A schema containing the tables used by current models and routes:
  - `users`
  - `roles`
  - `user_roles`
  - `sessions`
  - `files`
  - `files_comments`
  - `files_downloads`
  - `files_likes`
  - `albums`
  - `notifications`
  - `roadmap` (no SQL file in `migrations`; create manually)
- At least one default role row that can be assigned during registration.
  - Current code adds every newly registered user to role ID `1`.
- Valid object-storage credentials and bucket/public URL settings for upload/delete flows.
- `data/IP2LOCATION-LITE-DB11.IPV6.BIN` present for geo lookup during auth flows.

## Environment Variables

Copy `.env.example` to `.env` and fill values.

```bash
cp .env.example .env
```

### Application

| Variable | Required | Example | Notes |
| --- | --- | --- | --- |
| `APP_MODE` | no | `DEV` | Enables Bun SQL logging when `DEV`; also prefixes upload keys with `images/dev` |
| `APP_PORT` | no | `8080` | HTTP listen port; defaults to `8080` when empty |
| `JWT_SECRET` | yes | `change-me` | Shared HMAC secret used for access-token signing |

### Database

| Variable | Required | Example | Notes |
| --- | --- | --- | --- |
| `DB_USER` | yes | `file_uploader` | DB user |
| `DB_PASSWORD` | yes | `strong-password` | DB password |
| `DB_HOST` | yes | `127.0.0.1` | MySQL host |
| `DB_PORT` | yes | `3306` | MySQL port |
| `DB_NAME` | yes | `fileuploader` | Database name |

### Redis (currently not wired at startup)

| Variable | Required by runtime now | Example | Notes |
| --- | --- | --- | --- |
| `REDIS_HOST` | no | `127.0.0.1` | Reserved for future runtime Redis wiring |
| `REDIS_PORT` | no | `6379` | Reserved |
| `REDIS_PASSWORD` | no | `redis-password` | Reserved |

### Object Storage

These variables are required by upload/delete flows.

| Variable | Required for uploads | Example | Notes |
| --- | --- | --- | --- |
| `R2_ACCESS_KEY` | yes | `...` | Access key passed into storage client |
| `R2_SECRET_KEY` | yes | `...` | Secret key passed into storage client |
| `R2_BUCKET` | yes | `be-file-uploader-dev` | Bucket name |
| `R2_PUBLIC_URL` | yes | `https://cdn.example.com` | Public base URL used to build file URLs |

### Reserved OAuth Variables

The repository includes Discord-related variables in `.env.example`, but current app wiring does not use them.

- `DISCORD_CLIENT_ID`
- `DISCORD_CLIENT_SECRET`
- `DISCORD_LOCAL_REDIRECT_URI`
- `DISCORD_PRODUCTION_REDIRECT_URI`

## Local Run

Development:

```bash
go run ./cmd/app
```

Development with request logging:

```bash
go run ./cmd/app --debug
```

Production build:

```bash
go build -o be-file-uploader ./cmd/app
./be-file-uploader
```

Default local base URL:

```text
http://localhost:8080
```

## Runtime Behavior

### `--debug`

When `--debug` is enabled, the app logs:

- HTTP method
- route path
- response status
- latency
- client IP
- `X-User-Agent`
- raw request body
- query parameters
- response object
- registered routes at startup

### Routing Flags

Fiber app config currently enables:

- `StrictRouting: true`
- `CaseSensitive: true`

Contributors should use exact path casing and trailing slash behavior from route files.

### CORS

Current allowed origins:

- `http://localhost:3000`
- `http://localhost:8080`
- `https://uploader.dontkillme.lol`

Current allowed header values:

- `Origin`
- `Content-Type`
- `Accept`
- `Authorization`
- `Cache-Control`
- `X-Request-Fingerprint`
- `X-User-Agent`
- `Access-Control-Allow-Origin`
- `X-Locale`

Current allowed methods:

- `GET`
- `POST`
- `PUT`
- `DELETE`
- `PATCH`
- `OPTIONS`

Credentials are allowed.

To change CORS rules, edit [`internal/app/app.go`](../internal/app/app.go).

## First-Run Workflow

1. Create `.env` from `.env.example`.
2. Apply SQL from [`migrations`](../migrations).
3. Create missing schema objects not covered by `migrations` (`roadmap`, and optional `files_grants` if needed by your testing flow).
4. Ensure role ID `1` exists for registration.
5. Ensure object-storage vars are valid if testing `/private/storage/*` or `/public/storage/upload/sharex`.
6. Start the server.
7. Import [`docs/insomnia/be-file-uploader.insomnia.json`](insomnia/be-file-uploader.insomnia.json) or use [api-reference.md](api-reference.md).
8. Register, log in, then use the returned access token for private routes.

## Database / Seed Expectations

`migrations` files are provided, but there is no migration runner wired into app startup. Contributors should know these assumptions:

- Registration attaches users to role ID `1`.
- User administration routes require `MANAGE_USERS`.
- Role administration routes require `MANAGE_ROLES`.
- Upload and album mutation routes require `UPLOAD_FILES`.
- Own file listing requires `VIEW_OWN_FILES`.
- Cross-user file listing and post actions require `VIEW_OTHER_FILES`.
- Download action requires `DOWNLOAD_OTHERS_FILES`.
- Roadmap editing requires `DEVELOPER`.

## Storage Initialization Caveat

`CreateApp()` stores the storage pointer even when storage initialization returns an error value. This means the server can boot with broken upload configuration and fail later during upload/delete requests.
