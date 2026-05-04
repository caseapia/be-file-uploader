# be-file-uploader

`be-file-uploader` is a Go/Fiber backend for authenticated file uploads, user and role administration, albums, notifications, roadmap tasks, and ShareX-compatible uploads backed by MySQL and S3-compatible object storage.

## Features

- Username/password registration and login
- DB-backed sessions with JWT access tokens and rotating refresh tokens
- Private route authentication via bearer token or `access_token` cookie
- Route-level permission checks for protected operations
- Multipart upload flow (`init` -> `chunk` -> `complete`) for object storage
- Per-user storage quota tracking
- File listing, deletion, likes, comments, downloads, and album assignment
- Album creation, lookup, deletion, and moderation listing
- User administration, verification, role assignment, and ShareX token reset
- Notification list/read endpoints and public roadmap listing
- Importable Insomnia collection for local API exploration

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

The collection is aligned with the current route layout under `/v1/api/public/*` and `/v1/api/private/*`.

## Open-Source Notes

The docs intentionally describe real implementation behavior and known caveats so contributors can work against the current service behavior, not only intended behavior.
