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

## Directory Map

| Path | Responsibility |
| --- | --- |
| `cmd/app/main.go` | Bootstraps env loading, app startup, graceful shutdown |
| `internal/app/app.go` | Builds the Fiber app, middleware, CORS, dependency wiring, route registration |
| `internal/database` | MySQL connection creation and Bun setup |
| `internal/handler/auth` | Public auth handlers and private auth middleware |
| `internal/handler/user` | User lookup endpoints |
| `internal/handler/invite` | Invite admin endpoints |
| `internal/handler/image` | Storage/image endpoints |
| `internal/middleware` | Permission middleware |
| `internal/models` | Bun entities and request payload structs |
| `internal/repository/mysql` | Concrete Bun-backed repository operations |
| `internal/service/auth` | Registration, login, refresh, JWT parsing, session rules |
| `internal/service/user` | User lookup rules |
| `internal/service/invite` | Invite creation/revoke rules |
| `internal/service/image` | Upload/delete rules |
| `pkg/storage` | S3-compatible storage adapter |
| `pkg/utils` | Token, validation, hashing, ID generation, account helpers |
| `docs` | Maintained documentation and Insomnia export |

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
  -> repository
  -> response
```

## Routing Layout

Base prefix: `/v1/api`

Public:

- `/public/auth/register`
- `/public/auth/login`
- `/public/auth/refresh`

Private:

- `/private/user/*`
- `/private/invite/admin/*`
- `/private/storage/*`
- `/private/image/admin/*`

## Persistence Model

### Users

`users` contains:

- identity fields
- password hash
- IP/user-agent registration metadata
- invite relation
- many-to-many roles relation

### Sessions

`sessions` contains:

- session ID
- user ID
- IP address
- user agent
- active flag
- hard expiry
- refresh-token hash

### Invites

`invites` contains:

- generated code
- creator
- active flag
- optional consumed user

### Images

`images` contains:

- internal object key
- public URL
- original file name
- MIME type
- size
- uploader

## Permissions

Permission strings live in [`pkg/enums/role/permissions.go`](../pkg/enums/role/permissions.go).

Current permissions used by the app:

- `UPLOAD_FILES`
- `VIEW_OWN_FILES`
- `VIEW_OTHER_FILES`
- `VIEW_OTHER_PROFILES`
- `MANAGE_USERS`
- `MANAGE_FILES`
- `MANAGE_ROLES`
- `ADMIN`

Permission evaluation:

- a user has many roles
- a role has a JSON array of permissions
- `ADMIN` behaves as an override in `Role.HasPermission`

## Storage Design

The storage adapter is S3-compatible and currently configured with:

- endpoint: `https://storage.yandexcloud.net`
- region: `ru-central1`
- path-style object access

Public image URLs are constructed as:

```text
{r2_public_url}/{object-key}
```

Current object-key pattern:

```text
images/{userID}/{YYYY-MM}/{generatedID+ext}
```

## Developer Caveats

These are current implementation caveats worth knowing before contributing.

### `user/lookup/:id` permission logic

In [`internal/service/user/user.go`](../internal/service/user/user.go), the condition currently rejects lookup when:

- requester is not the target, or
- requester has `VIEW_OTHER_PROFILES`

This means the route does not currently behave like a normal "self or privileged lookup" rule.

### `storage/delete` permission logic

In [`internal/service/image/image.go`](../internal/service/image/image.go), deletion is denied when:

- requester is not uploader, or
- requester lacks `MANAGE_FILES`

Effectively, success requires both ownership and `MANAGE_FILES`.

### Storage init is not fail-fast

`CreateApp()` calls storage initialization but does not stop startup on storage initialization error. Upload/delete problems can surface only at request time.

### Route and permission naming split

Storage routes are under `/storage/*`, while admin list routes are under `/image/admin/*`. This is worth preserving in docs because frontend integrators otherwise expect one consistent resource prefix.

## Contribution Notes

When adding a new endpoint:

1. Add request structs to `internal/models/requests` if needed.
2. Add repository method(s) for DB access.
3. Add service method(s) for business logic.
4. Add handler method(s) and route registration.
5. Update [api-reference.md](api-reference.md).
6. Update [errors.md](errors.md) if new error codes are introduced.
7. Update the Insomnia export in [`docs/insomnia/be-file-uploader.insomnia.json`](insomnia/be-file-uploader.insomnia.json).
