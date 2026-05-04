# Architecture

## High-Level Structure

The service follows a layered structure:

```text
HTTP route -> handler -> service -> repository -> MySQL / object storage
```

Main goals of the current structure:

- keep HTTP parsing in handlers
- keep business rules in services
- keep SQL in repositories
- keep storage access in `pkg/storage`
- keep permission checks at route level and in sensitive service logic

## Directory Map

| Path                            | Responsibility                                                                       |
|---------------------------------|--------------------------------------------------------------------------------------|
| `cmd/app/main.go`               | Loads `.env`, initializes geo DB, starts Fiber, handles shutdown                     |
| `internal/app/app.go`           | Fiber config, CORS, dependency wiring, route registration                            |
| `internal/database`             | MySQL connection setup (Redis code exists but is currently disabled at runtime)      |
| `internal/handler/auth`         | Public auth handlers, private logout handler, auth middleware                        |
| `internal/handler/user`         | User profile, ShareX token generation, user admin endpoints                          |
| `internal/handler/file`         | Multipart upload flow, file listing/deletion, likes/comments/download, ShareX upload |
| `internal/handler/album`        | Album create/delete/lookup endpoints                                                 |
| `internal/handler/role`         | Role admin endpoints                                                                 |
| `internal/handler/notification` | Notification list/read endpoints                                                     |
| `internal/handler/roadmap`      | Public roadmap listing and developer-only roadmap editing                            |
| `internal/handler/developer`    | Public utility endpoint `/ping`                                                      |
| `internal/middleware`           | Route-level permission middleware (`ERR_NO_ACCESS`)                                  |
| `internal/models`               | Bun entities and JSON contracts                                                      |
| `internal/models/requests`      | Request payload structs + validation tags                                            |
| `internal/repository/mysql`     | Bun-backed repository operations                                                     |
| `internal/service/auth`         | Registration, login, refresh, logout, JWT/session rules                              |
| `internal/service/user`         | User lookup, private data exposure, user admin, ShareX token auth                    |
| `internal/service/file`         | Upload/delete/list, likes/comments/download, album assignment rules                  |
| `internal/service/album`        | Album create/delete/lookup rules                                                     |
| `internal/service/role`         | Role create/edit/delete rules                                                        |
| `internal/service/notification` | Notification creation/read rules                                                     |
| `internal/service/roadmap`      | Roadmap list/add/edit rules                                                          |
| `pkg/storage`                   | S3-compatible storage adapter                                                        |
| `pkg/enums`                     | Permission and roadmap status enums                                                  |
| `pkg/utils`                     | Token, validation, hashing, account helpers, ID generation                           |
| `migrations`                    | SQL schema files (not auto-executed on startup)                                      |
| `docs`                          | Maintained documentation and Insomnia export                                         |

## App Wiring

Application composition happens in [`internal/app/app.go`](../internal/app/app.go):

1. Parse `--debug` flag.
2. Connect to MySQL.
3. Build repository.
4. Initialize object storage.
5. Build services.
6. Build handlers.
7. Register route groups:
   - `/v1/api/public`
   - `/v1/api/private`
8. Attach auth middleware to the private group.
9. Attach route-level permission middleware where needed.

## Request Flow

### Public Auth Request

```text
Client
  -> Fiber route
  -> handler parses and validates body
  -> service executes business logic
  -> repository performs DB actions
  -> JSON response { "response": ... }
```

### Private Request

```text
Client
  -> CORS
  -> optional debug logger
  -> auth middleware
  -> optional permission middleware
  -> handler
  -> service
  -> repository/storage
  -> response
```

### Multipart Upload Request

```text
Client
  -> POST /private/storage/upload/init (JSON metadata)
  -> service validates quota and MIME type
  -> storage creates multipart upload
  -> POST /private/storage/upload/chunk (multipart chunk)
  -> storage uploads part
  -> POST /private/storage/upload/complete (parts metadata)
  -> storage completes multipart upload
  -> repository reserves quota and inserts file row
```

## Routing Layout

Base prefix: `/v1/api`

Public:

- `/public/ping`
- `/public/auth/register`
- `/public/auth/login`
- `/public/auth/refresh`
- `/public/storage/upload/sharex`
- `/public/roadmap/list`

Private:

- `/private/auth/logout`
- `/private/user/*`
- `/private/user/admin/*`
- `/private/storage/upload/*`
- `/private/storage/action/*`
- `/private/storage/post/:id`
- `/private/album/action/*`
- `/private/album/lookup*`
- `/private/roles/admin/*`
- `/private/notifications/*`
- `/private/roadmap/admin/*`

## Persistence Model

### Users

`users` contains:

- identity fields
- password hash
- IP/user-agent/Cloudflare metadata
- locale
- upload quota and used storage
- verification status
- ShareX token
- many-to-many roles relation
- has-many files and albums

### Sessions

`sessions` contains:

- session ID
- user ID
- IP address
- user agent
- active flag
- hard expiry
- refresh-token hash

### Roles

`roles` contains:

- name
- JSON permission array
- system-role flag
- creator ID
- display color

### Files

`files` contains:

- object key (`r2_key`)
- public URL
- original file name
- MIME type
- size
- uploader
- privacy flag
- optional album ID
- download count
- created timestamp

Related tables store likes and comments.

### Albums

`albums` contains:

- album name
- creator
- timestamps
- `is_public` option
- related file items

### Notifications

`notifications` contains:

- recipient user ID
- content string
- created timestamp
- read state

### Roadmap

`roadmap` table is used by services and models, but no migration file is currently shipped for it.

## Permissions

Permission strings live in [`pkg/enums/role/permissions.go`](../pkg/enums/role/permissions.go).

Current permissions:

- `UPLOAD_FILES`
- `VIEW_OWN_FILES`
- `VIEW_OTHER_FILES`
- `DOWNLOAD_OTHERS_FILES`
- `VIEW_OTHER_PROFILES`
- `MANAGE_USERS`
- `MANAGE_FILES`
- `MANAGE_ROLES`
- `ADMIN`
- `VIEW_PRIVATE_DATA`
- `ADMIN_CP`
- `SHOW_BADGE`
- `DEVELOPER`

Permission evaluation:

- a user has many roles
- a role has a JSON array of permissions
- `ADMIN` works as an override in `Role.HasPermission`
- route-level checks use `middleware.RequirePermission(...)`

## Storage Design

The storage adapter is S3-compatible and currently configured with:

- endpoint: `https://storage.yandexcloud.net`
- region: `ru-central1`
- path-style object access

Public file URLs are constructed as:

```text
{R2_PUBLIC_URL}/{object-key}
```

Current object-key pattern:

```text
images/{userID}/{YYYY-MM}/{generatedID+ext}
```

When `APP_MODE=DEV`, the key starts with:

```text
images/dev/{userID}/{YYYY-MM}/{generatedID+ext}
```

## Developer Caveats

### Storage init is not fail-fast

`CreateApp()` does not stop startup when `NewStorage(...)` returns an error, so upload/delete failures can appear only at request time.

### Invite flow is removed

Registration now accepts only `username` and `password`. It still auto-assigns role ID `1`.

### Access token extraction details

Auth middleware checks `Authorization` header first, then cookie `access_token`.

### Token lifetime is long-lived

`GenerateAccessToken` currently sets JWT expiry to `7 days`.

### Album privacy naming mismatch

`CreateAlbum` request uses field `is_private`, but handler passes it directly into model field `options.is_public`.

### Some ID parsing ignores conversion errors

Several handlers parse `:id` using `strconv.Atoi` and ignore conversion errors, so malformed path params may become `0`.

### `/private/storage/my` depends on schema assumptions

`SearchOwnFiles` query references grant alias `fg`. Contributors setting up a fresh DB should verify `files_grants` schema compatibility with current query behavior.

### ShareX upload response is custom

`POST /public/storage/upload/sharex` returns `{ "url": "..." }` directly instead of using the standard `response` wrapper.

## Contribution Notes

When adding or changing an endpoint:

1. Add request structs in `internal/models/requests` if needed.
2. Add repository method(s) for DB access.
3. Add service method(s) for business logic.
4. Add handler method(s) and route registration.
5. Add or update route-level permission middleware.
6. Update [api-reference.md](api-reference.md).
7. Update [errors.md](errors.md) if new error codes are introduced.
8. Update [`docs/insomnia/be-file-uploader.insomnia.json`](insomnia/be-file-uploader.insomnia.json).
