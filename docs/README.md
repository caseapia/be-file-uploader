# Documentation Index

This folder contains developer-facing and contributor-facing documentation for `be-file-uploader`.

## Audience

- Backend developers who need to run, extend, or review the service.
- Frontend and API integrators who need exact request/response contracts and auth flow details.
- Open-source contributors who need a quick map of the codebase and current caveats.

## Documents

| File                                                                               | Purpose                                                                                         |
|------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| [getting-started.md](getting-started.md)                                           | Local setup, infrastructure assumptions, environment variables, runtime flags, and first-run flow |
| [api-reference.md](api-reference.md)                                               | Full route catalog, request/response contracts, auth rules, and permission gates              |
| [authentication.md](authentication.md)                                             | JWT + refresh-token flow, session lifecycle, middleware behavior, and ShareX token auth        |
| [errors.md](errors.md)                                                             | Response envelopes, validation behavior, and current error catalog by feature area             |
| [architecture.md](architecture.md)                                                 | Project structure, request flow, persistence model, permission model, and implementation caveats |
| [insomnia/be-file-uploader.insomnia.json](insomnia/be-file-uploader.insomnia.json) | Importable Insomnia collection aligned with current routes                                      |

## What This Service Does

`be-file-uploader` is a Go/Fiber REST API that provides:

- Username/password registration and login.
- DB-backed sessions.
- JWT access tokens plus rotating refresh tokens.
- User profile lookup and user administration.
- Role and permission administration.
- Multipart file upload/list/delete endpoints backed by object storage.
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
| Private storage       | `POST /private/storage/upload/init`, `POST /private/storage/upload/chunk`, `POST /private/storage/upload/complete`, `POST /private/storage/action/delete`, `GET /private/storage/my`, `GET /private/storage/list`, `GET /private/storage/list/:id`, `GET /private/storage/post/:id` |
| Private storage actions | `PUT /private/storage/action/album/put`, `DELETE /private/storage/action/album/delete`, `PATCH /private/storage/action/like/:id`, `DELETE /private/storage/action/likeRemove/:id`, `GET /private/storage/action/download/:id`, `POST /private/storage/action/addComment` |
| Private albums        | `POST /private/album/action/create`, `DELETE /private/album/action/delete/:id`, `GET /private/album/lookup/:id`, `GET /private/album/lookupAll` |
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

## Important Caveats

- Registration no longer consumes an invite code; it creates a user and assigns role ID `1`.
- Access-token middleware accepts `Authorization: Bearer <jwt>` and cookie `access_token`.
- Access-token lifetime is currently `7 days` (same duration as session expiry).
- Redis env vars exist, but Redis is not currently initialized by `CreateDatabase()`.
- `migrations` does not include a `roadmap` table migration; contributors need to create it manually when setting up from scratch.
- `files_grants` is registered in Bun models but has no migration file in `migrations`; contributors testing `/private/storage/my` should verify schema compatibility.
