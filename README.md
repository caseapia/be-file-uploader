# be-file-uploader

`be-file-uploader` is a Go/Fiber backend for authenticated file uploads, profile and role management, albums, notifications, roadmap items, and ShareX-compatible upload flows backed by MySQL, Redis, and S3-compatible object storage.

## Features

- Username/password registration and login
- DB-backed sessions with JWT access tokens and rotating refresh tokens
- Private route authentication via bearer token or `auth_token` cookie
- Role/permission checks on protected endpoints
- Chunked multipart uploads to S3-compatible storage
- Per-user storage quota tracking
- File listing, deletion, likes, comments, downloads, and album assignment
- Album creation, lookup, and moderation listing
- User administration, verification, role assignment, and ShareX token management
- Notifications and public roadmap endpoints
- Insomnia collection for local API exploration

## Quick Start

```bash
cp .env.example .env
go run ./cmd/app
```

Default local URL:

```text
http://localhost:8080
```

## Documentation

- [docs/README.md](docs/README.md)
- [docs/getting-started.md](docs/getting-started.md)
- [docs/api-reference.md](docs/api-reference.md)
- [docs/authentication.md](docs/authentication.md)
- [docs/errors.md](docs/errors.md)
- [docs/architecture.md](docs/architecture.md)

## Insomnia

Import:

- [docs/insomnia/be-file-uploader.insomnia.json](docs/insomnia/be-file-uploader.insomnia.json)

The collection may lag behind newer API routes. Use the Markdown API reference as the source of truth for the current route catalog.

## Open-Source Notes

The documentation intentionally calls out current implementation caveats so contributors and [frontend](https://github.com/caseapia/fe-file-uploader) integrators can work against the real behavior of the service today, not only the intended behavior.
