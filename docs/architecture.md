# Architecture

## High-Level Structure

The service follows a layered structure:

```text
HTTP route -> handler -> service -> repository -> MySQL / Redis / object storage
```

Main goals of the current structure:

- keep HTTP parsing in handlers
- keep business rules in services
- keep SQL in repositories
- keep storage access in `pkg/storage`
- keep permission checks close to route registration and sensitive service logic

## Directory Map

| Path | Responsibility |
| --- | --- |
| `cmd/app/main.go` | Loads `.env`, starts Fiber, handles graceful shutdown |
| `internal/app/app.go` | Builds the Fiber app, middleware, CORS, dependency wiring, route registration |
| `internal/database` | MySQL and Redis connection creation |
| `internal/handler/auth` | Public auth handlers, private logout handler, and private auth middleware |
| `internal/handler/user` | User profile, ShareX token, and user admin endpoints |
| `internal/handler/file` | Multipart upload, file listing, file actions, album assignment, and ShareX upload endpoints |
| `internal/handler/album` | Album create/delete/lookup endpoints |
| `internal/handler/role` | Role admin endpoints |
| `internal/handler/notification` | Notification list/read endpoints |
| `internal/handler/roadmap` | Public roadmap list and developer-only roadmap editing endpoints |
| `internal/handler/developer` | Public utility endpoints such as `/ping` |
| `internal/middleware` | Permission middleware |
| `internal/models` | Bun entities |
| `internal/models/requests` | Request payload structs |
| `internal/repository/mysql` | Concrete Bun-backed repository operations |
| `internal/service/auth` | Registration, login, refresh, logout, JWT/session rules |
| `internal/service/user` | User lookup, private data exposure, user admin, ShareX token rules |
| `internal/service/file` | Upload, delete, visibility, like/comment/download, album assignment rules |
| `internal/service/album` | Album create/delete/lookup rules |
| `internal/service/role` | Role create/edit/delete rules |
| `internal/service/notification` | Notification creation/read rules |
| `internal/service/roadmap` | Roadmap list/add/edit rules |
| `pkg/storage` | S3-compatible storage adapter |
| `pkg/enums` | Permission and roadmap status enums |
| `pkg/utils` | Token, validation, hashing, ID generation, and account helpers |
| `migrations` | SQL schema files; not automatically executed by app startup |
| `docs` | Maintained documentation and Insomnia export |

## App Wiring

Application composition happens in [`internal/app/app.go`](../internal/app/app.go):

1. Parse `--debug` flag.
2. Connect to MySQL.
3. Connect to Redis.
4. Build repository.
5. Initialize object storage.
6. Build services.
7. Build handlers.
8. Register route groups:
   - `/v1/api/public`
   - `/v1/api/private`
9. Attach auth middleware to the private group.
10. Attach route-level permission middleware where needed.

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
  -> upload/init JSON metadata
  -> service validates quota and MIME type
  -> storage creates multipart upload
  -> upload/chunk multipart part(s)
  -> storage uploads part(s)
  -> upload/complete JSON part list
  -> storage completes multipart upload
  -> repository reserves quota and stores file row
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
- `/private/storage/*`
- `/private/album/*`
- `/private/roles/admin/*`
- `/private/notifications/*`
- `/private/roadmap/admin/*`

## Persistence Model

### Users

`users` contains:

- identity fields
- password hash
- IP/user-agent/Cloudflare registration metadata
- locale
- upload quota and used storage
- verification status
- ShareX token
- many-to-many roles relation
- has-many storage and albums relations

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

- internal object key
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
- public/private option
- related file items

### Notifications

`notifications` contains:

- recipient user ID
- content string
- created timestamp
- read state

### Roadmap

`roadmap` contains:

- task title
- status enum
- creator/updater relations
- created/updated timestamps

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
- `ADMIN` behaves as an override in `Role.HasPermission`
- route-level checks use `middleware.RequirePermission(...)`
- some services perform additional ownership/visibility checks

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

These are current implementation caveats worth knowing before contributing.

### Storage init is not fail-fast

`CreateApp()` calls storage initialization but does not stop startup on storage initialization error. Upload/delete problems can surface only at request time.

### Invite documentation is obsolete

Older docs described invite-based registration. Current registration accepts only `username` and `password`, and assigns role ID `1`.

### Album privacy naming mismatch

`CreateAlbum` accepts request field `is_private`, but passes that value into `AlbumOptions.IsPublic`. Frontend integrations should verify desired semantics before depending on this flag.

### Some ID parsing ignores errors

Several handlers convert `:id` with `strconv.Atoi` and ignore the returned error. A non-numeric path parameter can become `0` and then fail through normal lookup or permission behavior.

### File URL visibility has service-level rules

`File.ResolveURL` clears `url` for non-image files and for private files when the requester is not allowed to view the object URL. Some endpoints can return file metadata with an empty URL.

### ShareX upload has a custom response shape

`POST /public/storage/upload/sharex` returns `{ "url": "..." }` directly instead of using the shared `validation.Response` wrapper.

### Route naming is historical

File, post, storage, and album actions are split across `/storage/*` and `/album/*`. Preserve the exact paths in client integrations rather than inferring REST-style resource names.

## Contribution Notes

When adding a new endpoint:

1. Add request structs to `internal/models/requests` if needed.
2. Add repository method(s) for DB access.
3. Add service method(s) for business logic.
4. Add handler method(s) and route registration.
5. Add or update route-level permission middleware.
6. Update [api-reference.md](api-reference.md).
7. Update [errors.md](errors.md) if new error codes are introduced.
8. Update the Insomnia export in [`docs/insomnia/be-file-uploader.insomnia.json`](insomnia/be-file-uploader.insomnia.json) when practical.
