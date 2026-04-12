# Documentation Index

This folder contains developer-facing and frontend-facing documentation for `be-file-uploader`.

## Audience

- Backend developers who need to run, extend, or review the service.
- Frontend developers who need exact request/response contracts, auth flow details, and operational caveats.
- Open-source contributors who need a quick map of the codebase and local setup.

## Documents

| File | Purpose |
| --- | --- |
| [getting-started.md](getting-started.md) | Local setup, required infrastructure, environment variables, runtime flags, and first-run workflow |
| [api-reference.md](api-reference.md) | Full route catalog, request/response contracts, auth rules, and frontend notes per endpoint |
| [authentication.md](authentication.md) | JWT + refresh-token flow, session model, browser/API-client integration guidance |
| [errors.md](errors.md) | Response envelope, validation behavior, known error codes, and frontend error-handling guidance |
| [architecture.md](architecture.md) | Project structure, request flow, persistence model, permissions, and known implementation caveats |
| [insomnia/be-file-uploader.insomnia.json](insomnia/be-file-uploader.insomnia.json) | Importable Insomnia collection with environments, auth automation, and storage route helpers |

## What This Service Does

`be-file-uploader` is a Go/Fiber REST API that provides:

- Invite-based user registration.
- Login with DB-backed sessions.
- JWT access tokens plus rotating refresh tokens.
- User profile lookup.
- Invite management for privileged users.
- Image upload/list/delete endpoints backed by object storage.

## API Summary

Base prefix: `/v1/api`

| Scope | Routes |
| --- | --- |
| Public | `POST /public/auth/register`, `POST /public/auth/login`, `POST /public/auth/refresh` |
| Private user | `GET /private/user/me`, `GET /private/user/lookup/:id` |
| Private storage | `POST /private/storage/upload`, `POST /private/storage/delete`, `GET /private/storage/my` |
| Private admin | `GET /private/invite/admin/list`, `POST /private/invite/admin/create`, `DELETE /private/invite/admin/revoke`, `GET /private/image/admin/list`, `GET /private/image/admin/list/:id` |

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
- The server stores `X-User-Agent` in session metadata and updates it on refresh. Set it consistently from frontend clients.
- Refresh token rotation is mandatory: after every successful `/public/auth/refresh`, overwrite both stored tokens immediately.
- Image upload uses `multipart/form-data` with a single file field named `image`.
- The collection export already saves tokens, invite IDs, and uploaded image IDs into the active environment.

## Important Caveats

The docs intentionally describe both the intended contracts and a few implementation details that matter in practice:

- `user/lookup/:id` currently has permission logic that can reject lookups unexpectedly. See [api-reference.md](api-reference.md) and [architecture.md](architecture.md).
- `storage/delete` currently behaves more strictly than the route name suggests. See the route notes before wiring a delete button in the frontend.
- Storage configuration variables are required for upload/delete flows, but the app currently does not fail fast when storage initialization fails.

Those caveats are documented because this repository is expected to be consumed as open source and integrators need the current behavior, not only the intended one.
