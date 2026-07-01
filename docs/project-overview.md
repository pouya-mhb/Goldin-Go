# Goldin-Go Project Overview

## Purpose

Goldin-Go is a Go backend service with a small React frontend for identity-focused user operations. The implemented backend currently centers on user registration and login, backed by MySQL persistence, bcrypt password hashing, and JWT token issuance.

The repository is structured as a modular application, with platform infrastructure separated from the identity bounded context.

## Repository Layout

```text
cmd/api/                         Application entrypoint
internal/platform/               Shared platform infrastructure
internal/platform/bootstrap/     Application composition and lifecycle
internal/platform/config/        Environment-based configuration
internal/platform/database/      MySQL connection setup
internal/platform/http/          HTTP server, router, middleware, health route
internal/platform/logger/        slog logger setup
internal/identity/               Identity bounded context
internal/identity/domain/        User aggregate and value objects
internal/identity/application/   Commands, ports, services, and results
internal/identity/adapters/      HTTP, persistence, password, and token adapters
migrations/mysql/                MySQL schema migrations
web/                             React + Vite frontend
```

## Backend Runtime Flow

The backend starts from `cmd/api/main.go`.

Startup flow:

1. Create a cancellable context that responds to interrupt and termination signals.
2. Build the application through `internal/platform/bootstrap.Build`.
3. Load and validate environment configuration.
4. Create the structured logger.
5. Open and verify the MySQL connection.
6. Build the Identity module.
7. Register HTTP routes and start the HTTP server.
8. Shut down gracefully when the process context is canceled.

The app lifecycle is owned by `internal/platform/bootstrap.App`. Its `Run` method starts the HTTP server and closes the server and database connection during graceful shutdown.

## Configuration

Configuration is loaded from environment variables, with `.env` support for local development.

Important groups:

- App: `APP_NAME`, `APP_ENV`, `APP_VERSION`
- Server: `SERVER_HOST`, `SERVER_PORT`
- Database: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Database pool: `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`
- Redis placeholders: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`
- JWT: `JWT_SECRET`, `JWT_ACCESS_DURATION`, `JWT_REFRESH_DURATION`
- Kafka placeholders: `KAFKA_BROKERS`, `KAFKA_CLIENT_ID`
- Logging: `LOG_LEVEL`

Validation currently requires database host, user, and name; a JWT secret of at least 32 characters; valid server/database/Redis ports; a valid log level; and sane database pool settings.

## Platform Layer

The platform layer provides shared infrastructure:

- Bootstrap composition
- HTTP server and router
- Middleware registration
- Health endpoint at `GET /health`
- MySQL database connection and DSN construction
- Structured logging through Go `slog`
- Environment configuration loading and validation

The HTTP server uses `chi` for routing.

## Identity Module

The identity module is the main implemented bounded context.

It exposes two inbound use cases:

- Register user
- Login user

Key components:

- Domain: user entity, email value object, password hash value object, user ID value object, and user registration event.
- Application services: registration and login orchestration.
- Ports: inbound use-case interfaces and outbound dependencies for persistence, hashing, verification, and token issuing.
- Adapters:
  - HTTP handlers for `/identity/register` and `/identity/login`
  - MySQL user repository
  - bcrypt password hasher/verifier
  - JWT issuer

## HTTP API

Implemented endpoints:

```text
GET  /health
POST /identity/register
POST /identity/login
```

Registration accepts an email and password, validates them through the application/domain layer, hashes the password, persists the user, and returns the created user ID, email, and registration timestamp.

Login accepts an email and password, verifies credentials, and returns access and refresh tokens with expiry metadata.

## Database

The project uses MySQL.

The first migration creates the `identity_users` table with:

- `id`
- `email`
- `password_hash`
- `registered_at`
- `created_at`
- `updated_at`

The email field is unique, and registration time is indexed.

## Frontend

The `web/` folder contains a React 19 + Vite frontend.

Current frontend behavior:

- Checks API health through `/health`
- Displays API online/offline state
- Provides a user registration form
- Calls `/identity/register`
- Displays successful registration details or API errors

The frontend is intentionally small and currently focused on the identity registration workflow.

## Testing Status

The backend test suite currently passes with:

```powershell
go test ./...
```

Covered areas include:

- Identity module construction
- Domain and value objects
- Registration and login services
- HTTP identity handlers
- Password hashing
- JWT issuing
- MySQL repository behavior
- Platform database helpers
- Platform HTTP server and routes

## Current State

The backend has moved beyond a scaffold and now has a working identity slice with tests. The project still appears early-stage because several areas are placeholders or minimal:

- `README.md` is empty.
- Top-level `Makefile` is empty.
- Top-level `docker-compose.yml` is empty.
- Redis and Kafka are represented in configuration but not wired into runtime infrastructure.
- The frontend currently supports registration but not login.
- There are no OpenAPI or proto contract files yet.

## Suggested Next Steps

1. Fill in `README.md` with setup, configuration, migration, test, and run instructions.
2. Add Docker Compose services for MySQL and any planned local dependencies.
3. Add a Makefile or task runner commands for common workflows.
4. Document the HTTP API with OpenAPI.
5. Add frontend login support.
6. Add migration tooling instructions.
7. Decide whether Redis and Kafka are near-term requirements or should remain deferred placeholders.
