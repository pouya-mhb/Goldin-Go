# Goldin-Go

Goldin-Go is a Go API service with a small React frontend. The current product slice focuses on identity: user registration, login, MySQL persistence, bcrypt password hashing, and JWT token issuance.

For a fuller architecture summary, see [docs/project-overview.md](docs/project-overview.md).

## Project Structure

```text
cmd/api/                         API entrypoint
internal/platform/               Shared infrastructure
internal/platform/bootstrap/     App composition and lifecycle
internal/platform/config/        Environment configuration
internal/platform/database/      MySQL connection setup
internal/platform/http/          HTTP server, router, middleware, health route
internal/platform/logger/        slog logger setup
internal/identity/               Identity bounded context
migrations/mysql/                MySQL migrations
web/                             React + Vite frontend
```

## Requirements

- Go 1.26.4
- Node.js and npm for the frontend
- Docker for local MySQL
- `make`, optional but recommended
- `migrate` CLI for database migrations

## Configuration

Create a local environment file from the example:

```powershell
Copy-Item .env.example .env
```

Important variables:

- `SERVER_HOST`, `SERVER_PORT`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET`
- `JWT_ACCESS_DURATION`, `JWT_REFRESH_DURATION`
- `LOG_LEVEL`

`JWT_SECRET` must be at least 32 characters.

## Local Development

Start MySQL:

```powershell
docker compose up -d
```

Run database migrations:

```powershell
migrate -path migrations/mysql -database "mysql://root:change-me-local-only@tcp(localhost:3307)/goldin" up
```

Run the API:

```powershell
go run ./cmd/api
```

The API defaults to:

```text
http://localhost:18080
```

Run the frontend:

```powershell
Set-Location web
npm install
npm run dev
```

## API Endpoints

```text
GET  /health
POST /identity/register
POST /identity/login
```

Example registration request:

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:18080/identity/register `
  -ContentType "application/json" `
  -Body '{"email":"user@example.com","password":"very-secure-password"}'
```

Example login request:

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:18080/identity/login `
  -ContentType "application/json" `
  -Body '{"email":"user@example.com","password":"very-secure-password"}'
```

## Make Targets

The repository includes these helper targets:

```powershell
make test
make compose-up
make compose-down
make migrate-up
make migrate-down
make migrate-force VERSION=1
```

The migration targets use `DB_DSN`. Override it when needed:

```powershell
make migrate-up DB_DSN="mysql://root:change-me-local-only@tcp(localhost:3307)/goldin"
```

## Tests

Run backend tests:

```powershell
go test ./...
```

Run frontend build checks:

```powershell
Set-Location web
npm run build
```

## Current Status

Implemented:

- Identity registration and login use cases
- MySQL user repository
- bcrypt password hashing
- JWT access and refresh tokens
- HTTP handlers and health route
- React registration UI
- MySQL migration for `identity_users`

Planned or still minimal:

- README-driven production deployment notes
- OpenAPI or proto contracts
- Frontend login screen
- Redis and Kafka runtime integration
