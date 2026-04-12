# Getting Started

## Stack

- Go `1.25.0` according to [`go.mod`](../go.mod)
- Fiber v3
- MySQL
- Bun ORM
- S3-compatible object storage via Yandex Object Storage client configuration

## Prerequisites

Before starting the API locally you need:

- A MySQL database reachable from the app.
- A schema containing at least the tables used by the Bun models:
  - `users`
  - `roles`
  - `user_roles`
  - `invites`
  - `sessions`
  - `images`
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
| `APP_MODE` | no | `DEV` | Enables Bun SQL logging when equal to `DEV` |
| `APP_PORT` | no | `8080` | HTTP listen port |
| `JWT_SECRET` | yes | `change-me` | Shared HMAC secret used for access-token signing |

### Database

| Variable | Required | Example | Notes |
| --- | --- | --- | --- |
| `DB_USER` | yes | `file_uploader` | Use a dedicated DB user |
| `DB_PASSWORD` | yes | `strong-password` | Avoid root credentials |
| `DB_HOST` | yes | `127.0.0.1` | MySQL host |
| `DB_PORT` | yes | `3306` | MySQL port |
| `DB_NAME` | yes | `fileuploader` | Database name |

### Object Storage

These variables are required by upload/delete flows.

| Variable | Required for uploads | Example | Notes |
| --- | --- | --- | --- |
| `R2_ACCESS_KEY` | yes | `...` | Access key passed into storage client |
| `R2_SECRET_KEY` | yes | `...` | Secret key passed into storage client |
| `r2_bucket` | yes | `be-file-uploader-dev` | Exact lowercase key expected by current code |
| `r2_public_url` | yes | `https://cdn.example.com` | Exact lowercase key expected by current code |

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

Current allowed request headers:

- `Origin`
- `Content-Type`
- `Accept`
- `Authorization`
- `Cache-Control`
- `X-Request-Fingerprint`
- `X-User-Agent`

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
3. Seed roles and at least one privileged user if you need invite administration.
4. If you want to test uploads, configure object storage variables.
5. Start the server.
6. Import the Insomnia collection from [`docs/insomnia/be-file-uploader.insomnia.json`](insomnia/be-file-uploader.insomnia.json).
7. Use `Create invite` as an admin, then `Register`, then `Login` or `Refresh` as needed.

## Database / Seed Expectations

The codebase does not currently ship migrations or seeders in this repository. Contributors need to know the following assumptions:

- Registration attaches users to role ID `1`.
- Invite creation requires a role permission `MANAGE_USERS`.
- Admin image listing requires `VIEW_OTHER_FILES`.
- Profile lookup for other users is intended to depend on `VIEW_OTHER_PROFILES`, but current behavior is stricter due to service logic.

If you open-source this repository, consider adding migrations later. Until then, document the schema externally or provide SQL bootstrap scripts.

## Storage Caveat

`CreateApp()` initializes object storage but does not stop startup if storage creation fails. That means:

- the server can boot even with broken upload configuration
- upload/delete requests may fail later at runtime

For local development, validate storage env vars before testing `/private/storage/*`.
