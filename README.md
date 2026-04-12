# be-file-uploader

`be-file-uploader` is a Go/Fiber backend that provides invite-based registration, JWT authentication with rotating refresh tokens, user/profile endpoints, invite administration, and image upload flows backed by object storage.

## Features

- Invite-only user registration
- Login with DB-backed sessions
- Access token + refresh token flow
- Role/permission checks on protected endpoints
- Image upload, listing, and deletion APIs
- Insomnia collection with automated token and ID propagation

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

The collection includes:

- Local, staging, and production environments
- auto-save for `access_token` and `refresh_token`
- helper variables for invite/image/user IDs
- storage upload request with multipart file field

## Open-Source Notes

The documentation intentionally calls out a few current implementation caveats so contributors and frontend integrators can work against the real behavior of the service today, not only the intended behavior.
