# Getting Started

## Stack

- Go `1.25.0` according to [`go.mod`](../go.mod)
- Fiber v3
- MySQL
- Redis
- Bun ORM
- S3-compatible object storage configured for Yandex Object Storage

## Prerequisites

Before starting the API locally you need:

- A MySQL database reachable from the app.
- A Redis instance reachable from the app.
- A schema containing the tables used by the Bun models:
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
  - `roadmap`
- At least one default role row that can be assigned during registration.
  - Current code adds every newly registered user to role ID `1`.
- If you need upload/delete testing, valid object-storage credentials and bucket/public URL settings.

## Environment Variables

Copy `.env.example` to `.env` and fill the values.

```bash
cp .env.example .env
```

### Application

| Variable | Required | Example | Notes |
| --- | --- | --- | --- |
| `APP_MODE` | no | `DEV` | Enables Bun SQL logging when equal to `DEV`; also prefixes upload keys with `images/dev` |
| `APP_PORT` | no | `8080` | HTTP listen port; defaults to `8080` when empty |
| `JWT_SECRET` | yes | `change-me` | Shared HMAC secret used for access-token signing |

### Database

| Variable | Required | Example | Notes |
| --- | --- | --- | --- |
| `DB_USER` | yes | `file_uploader` | Use a dedicated DB user |
| `DB_PASSWORD` | yes | `strong-password` | Avoid root credentials |
| `DB_HOST` | yes | `127.0.0.1` | MySQL host |
| `DB_PORT` | yes | `3306` | MySQL port |
| `DB_NAME` | yes | `fileuploader` | Database name |

### Redis

| Variable | Required | Example | Notes |
| --- | --- | --- | --- |
| `REDIS_HOST` | yes | `127.0.0.1` | Redis host |
| `REDIS_PORT` | yes | `6379` | Redis port |
| `REDIS_PASSWORD` | no | `redis-password` | Leave empty for a local Redis without auth |

### Object Storage

These variables are required by upload/delete flows.

| Variable | Required for uploads | Example | Notes |
| --- | --- | --- | --- |
| `R2_ACCESS_KEY` | yes | `...` | Access key passed into storage client |
| `R2_SECRET_KEY` | yes | `...` | Secret key passed into storage client |
| `R2_BUCKET` | yes | `be-file-uploader-dev` | Bucket name passed into storage client |
| `R2_PUBLIC_URL` | yes | `https://cdn.example.com` | Public base URL used to build file URLs |

### Reserved OAuth Variables

The repository includes Discord-related variables in `.env.example`, but the current app wiring does not use them yet.

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
- registered routes at startup

Use it only in development because it can log request bodies containing credentials.

### CORS

Current allowed origins:

- `http://localhost:3000`
- `http://localhost:8080`
- `https://uploader.dontkillme.lol`

Current allowed request headers:

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
2. Ensure MySQL schema exists and the app can connect.
3. Ensure Redis is reachable.
4. Seed roles, including role ID `1`, and add privileged roles/users if you need admin routes.
5. If you want to test uploads, configure object storage variables.
6. Start the server.
7. Import the Insomnia collection from [`docs/insomnia/be-file-uploader.insomnia.json`](insomnia/be-file-uploader.insomnia.json), or call the routes from [api-reference.md](api-reference.md).
8. Register, log in, then use the returned access token for private routes.

## Database / Seed Expectations

SQL files are present in [`migrations`](../migrations), but there is no migration runner wired into application startup. Contributors need to know the following assumptions:

- Registration attaches users to role ID `1`.
- User administration requires `MANAGE_USERS`.
- Role administration requires `MANAGE_ROLES`.
- Uploads and album mutations require `UPLOAD_FILES`.
- Own file listing requires `VIEW_OWN_FILES`.
- Cross-user file listing and post actions require `VIEW_OTHER_FILES`.
- Download action requires `DOWNLOAD_OTHERS_FILES`.
- Roadmap editing requires `DEVELOPER`.

## Storage Caveat

`CreateApp()` initializes object storage and stores the returned pointer even if initialization returns an error. That means:

- the server can boot even with broken upload configuration
- upload/delete requests may fail later at runtime

For local development, validate storage env vars before testing `/private/storage/*` or `/public/storage/upload/sharex`.
