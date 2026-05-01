# Documentation Index

This folder contains developer-facing and frontend-facing documentation for `be-file-uploader`.

## Audience

- Backend developers who need to run, extend, or review the service.
- Frontend developers who need exact request/response contracts, auth flow details, and operational caveats.
- Open-source contributors who need a quick map of the codebase and local setup.

## Documents

| File                                                                               | Purpose                                                                                         |
|------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| [getting-started.md](getting-started.md)                                           | Local setup, infrastructure, environment variables, runtime flags, and first-run workflow        |
| [api-reference.md](api-reference.md)                                               | Full route catalog, request/response contracts, auth rules, permissions, and frontend notes      |
| [authentication.md](authentication.md)                                             | JWT + refresh-token flow, session model, logout behavior, ShareX token usage, and client advice |
| [errors.md](errors.md)                                                             | Response envelope, validation behavior, known error codes, and frontend error-handling guidance |
| [architecture.md](architecture.md)                                                 | Project structure, request flow, persistence model, permissions, and implementation caveats      |
| [insomnia/be-file-uploader.insomnia.json](insomnia/be-file-uploader.insomnia.json) | Importable Insomnia collection for API exploration                                              |

## What This Service Does

`be-file-uploader` is a Go/Fiber REST API that provides:

- Username/password registration and login.
- DB-backed sessions.
- JWT access tokens plus rotating refresh tokens.
- User profile lookup and user administration.
- Role and permission administration.
- Chunked file upload/list/delete endpoints backed by object storage.
- Albums, likes, comments, downloads, and notifications.
- Public roadmap listing with developer-only roadmap editing.
- ShareX token generation and public ShareX upload.

## API Summary

Base prefix: `/v1/api`

| Scope                 | Routes |
|-----------------------|--------|
| Public auth           | `POST /public/auth/register`, `POST /public/auth/login`, `POST /public/auth/refresh` |
| Public utility        | `GET /public/ping`, `GET /public/roadmap/list`, `POST /public/storage/upload/sharex` |
| Private auth          | `DELETE /private/auth/logout` |
| Private user          | `GET /private/user/me`, `GET /private/user/lookup/:id`, `GET /private/user/shareX/generate` |
| Private user admin    | `GET /private/user/admin/users`, `PUT /private/user/admin/role/add`, `DELETE /private/user/admin/role/delete`, `PATCH /private/user/admin/storage-limit/update`, `PATCH /private/user/admin/verify/:id`, `DELETE /private/user/admin/shareX/reset/:id` |
| Private storage       | `POST /private/storage/upload/init`, `POST /private/storage/upload/chunk`, `POST /private/storage/upload/complete`, `POST /private/storage/delete`, `GET /private/storage/my`, `GET /private/storage/list`, `GET /private/storage/list/:id`, `GET /private/storage/post/:id` |
| Private post actions  | `PATCH /private/storage/post/action/like/:id`, `DELETE /private/storage/post/action/likeRemove/:id`, `GET /private/storage/post/action/download/:id`, `POST /private/storage/post/action/addComment` |
| Private albums        | `POST /private/album/create`, `DELETE /private/album/delete/:id`, `GET /private/album/lookup/:id`, `GET /private/album/lookupAll` |
| Private roles admin   | `GET /private/roles/admin/all`, `POST /private/roles/admin/create`, `PATCH /private/roles/admin/edit`, `DELETE /private/roles/admin/delete` |
| Private notifications | `GET /private/notifications/my`, `PATCH /private/notifications/action/read/:id` |
| Private roadmap admin | `POST /private/roadmap/admin/task/add`, `PATCH /private/roadmap/admin/task/edit` |

## Response Conventions

Successful JSON responses are wrapped:

```json
{
  "response": {}
}
```

Most handled errors are returned as:

```json
{
  "error": "ERR_CODE",
  "code": 400
}
```

Unhandled internal failures use:

```json
{
  "code": 500,
  "message": "raw error text"
}
```

## Frontend Quick Notes

- Private endpoints accept the access token either in `Authorization: Bearer <jwt>` or in the `auth_token` cookie.
- The server stores `X-User-Agent` in session metadata during login, refresh, and logout. Set it consistently from frontend clients.
- Refresh token rotation is mandatory: after every successful `/public/auth/refresh`, overwrite both stored tokens immediately.
- Uploads are multipart and currently use an init/chunk/complete flow instead of a single private upload request.
- Public ShareX upload uses a generated ShareX token and returns `{ "url": "..." }` without the standard `response` wrapper.
- Private file URLs may be blanked in responses when the requester is not allowed to view the underlying object URL.

## Important Caveats

The docs intentionally describe both the intended contracts and implementation details that matter in practice:

- Registration no longer consumes an invite code; it creates a user and assigns role ID `1`.
- Storage configuration variables are required for upload/delete flows, but the app currently does not fail fast when storage initialization fails.
- Some path parameters are converted with `strconv.Atoi` and ignore conversion errors, so invalid IDs can behave like `0`.
- The Insomnia export may not include every newer route. Use [api-reference.md](api-reference.md) as the current route source of truth.
