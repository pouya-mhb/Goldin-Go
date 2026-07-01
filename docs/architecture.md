# Goldin-Go Architecture

Goldin-Go is currently a small full-stack project centered on an Identity bounded context. The Go API follows a ports-and-adapters shape: HTTP handlers adapt incoming requests into application commands, application services coordinate use cases through interfaces, and concrete adapters handle persistence, password hashing, and token issuing.

## System Context

```mermaid
flowchart LR
    User["User"]
    Browser["Browser"]
    Web["React/Vite web app"]
    API["Goldin Go API"]
    MySQL["MySQL 8.4"]
    JWT["JWT tokens"]

    User --> Browser
    Browser --> Web
    Web -- "GET /health" --> API
    Web -- "POST /identity/register" --> API
    Browser -- "POST /identity/login" --> API
    API -- "read/write identity_users" --> MySQL
    API -- "issues signed access and refresh tokens" --> JWT
```

The frontend is a Vite React app. During development, Vite proxies `/health` and `/identity/*` to the API target configured by `VITE_API_TARGET`, defaulting to `http://127.0.0.1:8081`.

## Startup And Dependency Wiring

```mermaid
flowchart TD
    Main["cmd/api/main.go"]
    Build["bootstrap.Build()"]
    Config["config.Load() + Validate()"]
    Logger["logger.New(cfg)"]
    DB["database.OpenMySQL()"]
    Identity["identity.NewModule(db, cfg.JWT)"]
    Router["platform/http.NewRouter()"]
    Server["platform/http.New(cfg.Server, router)"]
    App["bootstrap.App"]
    Run["app.Run(ctx)"]

    Main --> Build
    Build --> Config
    Build --> Logger
    Build --> DB
    Build --> Identity
    Build --> Router
    Build --> Server
    Build --> App
    Main --> Run
    Run --> Server
```

Startup is intentionally centralized in `internal/platform/bootstrap`. That keeps construction of infrastructure and modules out of business logic. `App.Run` starts the HTTP server, listens for context cancellation, then performs graceful HTTP shutdown and closes the database connection.

## Backend Layers

```mermaid
flowchart TB
    subgraph Transport["Transport adapters"]
        Health["/health route"]
        IdentityHTTP["Identity HTTP handlers"]
    end

    subgraph Application["Application layer"]
        Inbound["Inbound ports"]
        Registration["RegistrationService"]
        Login["LoginService"]
        Outbound["Outbound ports"]
    end

    subgraph Domain["Domain layer"]
        User["User aggregate"]
        VOs["Email, UserID, PasswordHash"]
        Event["UserRegistered event"]
    end

    subgraph Infrastructure["Infrastructure adapters"]
        Repo["MySQLUserRepository"]
        Bcrypt["BcryptHasher"]
        Token["JWTIssuer"]
        SQL["identity_users table"]
    end

    IdentityHTTP --> Inbound
    Inbound --> Registration
    Inbound --> Login
    Registration --> Domain
    Login --> Domain
    Registration --> Outbound
    Login --> Outbound
    Outbound --> Repo
    Outbound --> Bcrypt
    Outbound --> Token
    Repo --> SQL
    User --> VOs
    User --> Event
```

The dependency direction points inward. The application services know about ports, commands, results, and domain objects. They do not know the details of HTTP, SQL queries, bcrypt, or JWT signing. Those details are plugged in by adapters during module construction.

## Identity Module Composition

```mermaid
flowchart LR
    DB["*sql.DB"]
    JWTConfig["JWT config"]
    Module["identity.Module"]
    Repo["MySQLUserRepository"]
    Hasher["BcryptHasher"]
    Issuer["JWTIssuer"]
    RegisterSvc["RegistrationService"]
    LoginSvc["LoginService"]
    RegisterPort["RegisterUser inbound port"]
    LoginPort["LoginUser inbound port"]

    DB --> Repo
    JWTConfig --> Issuer
    Repo --> RegisterSvc
    Hasher --> RegisterSvc
    Repo --> LoginSvc
    Hasher --> LoginSvc
    Issuer --> LoginSvc
    RegisterSvc --> RegisterPort
    LoginSvc --> LoginPort
    RegisterPort --> Module
    LoginPort --> Module
```

`identity.NewModule` is the composition root for the Identity bounded context. It creates a MySQL repository, a bcrypt password hasher/verifier, a JWT issuer, then exposes only the inbound use case interfaces needed by HTTP.

## Registration Workflow

```mermaid
sequenceDiagram
    participant Client
    participant Handler as RegistrationHandler
    participant Service as RegistrationService
    participant Repo as UserRepository
    participant Hasher as PasswordHasher
    participant Domain as User aggregate
    participant DB as MySQL

    Client->>Handler: POST /identity/register
    Handler->>Handler: Decode strict JSON
    Handler->>Service: RegisterUser(command)
    Service->>Service: Validate email value object
    Service->>Service: Enforce minimum password length
    Service->>Repo: ExistsByEmail(email)
    Repo->>DB: SELECT EXISTS(...)
    DB-->>Repo: exists true/false
    Repo-->>Service: exists
    alt email already exists
        Service-->>Handler: ErrEmailAlreadyRegistered
        Handler-->>Client: 409 email_already_registered
    else email is available
        Service->>Hasher: HashPassword(plaintext)
        Hasher-->>Service: bcrypt hash
        Service->>Domain: RegisterUser(id, email, hash, now)
        Domain-->>Service: User + UserRegistered event
        Service->>Repo: Save(user)
        Repo->>DB: INSERT identity_users
        DB-->>Repo: inserted
        Service-->>Handler: RegisterUser result
        Handler-->>Client: 201 user_id, email, registered_at
    end
```

Registration is defensive at several layers. The HTTP layer rejects invalid JSON and unknown fields. The service validates email and password rules, checks uniqueness before hashing, and still handles the database unique constraint as a conflict. The domain aggregate records a `UserRegistered` event, although no event publisher is wired yet.

## Login Workflow

```mermaid
sequenceDiagram
    participant Client
    participant Handler as LoginHandler
    participant Service as LoginService
    participant Repo as UserRepository
    participant Verifier as PasswordVerifier
    participant Issuer as TokenIssuer
    participant DB as MySQL

    Client->>Handler: POST /identity/login
    Handler->>Handler: Decode strict JSON
    Handler->>Service: LoginUser(command)
    Service->>Service: Normalize and validate email
    Service->>Repo: FindByEmail(email)
    Repo->>DB: SELECT user by email
    DB-->>Repo: row or no rows
    alt user missing or email invalid
        Service-->>Handler: ErrInvalidCredentials
        Handler-->>Client: 401 invalid_credentials
    else user found
        Repo-->>Service: User aggregate
        Service->>Verifier: VerifyPassword(plaintext, hash)
        alt password mismatch
            Service-->>Handler: ErrInvalidCredentials
            Handler-->>Client: 401 invalid_credentials
        else password matches
            Service->>Issuer: IssueTokens(user id, email)
            Issuer-->>Service: access token + refresh token
            Service-->>Handler: LoginUser result
            Handler-->>Client: 200 tokens and expiry metadata
        end
    end
```

Login deliberately collapses invalid email, missing user, and password mismatch into the same public error: `invalid_credentials`. That prevents the API from revealing which email addresses are registered.

## HTTP Request Pipeline

```mermaid
flowchart LR
    Request["Incoming request"]
    RequestID["RequestID middleware"]
    Recovery["Recovery middleware"]
    Logging["Logging middleware"]
    Route["chi route"]
    Handler["Handler"]
    Response["JSON or text response"]
    Logs["structured slog entry"]

    Request --> RequestID
    RequestID --> Recovery
    Recovery --> Logging
    Logging --> Route
    Route --> Handler
    Handler --> Response
    Logging --> Logs
```

Every request receives or reuses an `X-Request-ID`. Panics are recovered into a 500 response and logged with request context. Completed requests are logged with method, path, status, request ID, and duration.

## Data Model

```mermaid
erDiagram
    IDENTITY_USERS {
        char36 id PK
        varchar254 email UK
        varchar255 password_hash
        timestamp6 registered_at
        timestamp6 created_at
        timestamp6 updated_at
    }
```

The only database table in the current migrations is `identity_users`. The email uniqueness constraint is part of the business rule enforcement, not just an index optimization.

## Frontend Flow

```mermaid
flowchart TD
    React["React App"]
    Health["API status button"]
    Form["Registration form"]
    Vite["Vite dev proxy"]
    API["Go API"]
    Result["Result panel"]

    React --> Health
    React --> Form
    Health -- "GET /health" --> Vite
    Form -- "POST /identity/register" --> Vite
    Vite --> API
    API --> Vite
    Vite --> Result
```

The web app currently focuses on registration and API health. It uses local component state for form fields, submission state, API status, success payloads, and errors. It does not yet include a login screen even though the backend login endpoint exists.

## Configuration And Operations

```mermaid
flowchart TD
    Env["Environment variables or .env"]
    Config["config.Load()"]
    Validate["config.Validate()"]
    MySQLCompose["docker compose MySQL"]
    Migrations["migrate up/down"]
    API["Go API"]
    Tests["go test ./..."]

    Env --> Config
    Config --> Validate
    Validate --> API
    MySQLCompose --> API
    Migrations --> MySQLCompose
    Tests --> API
```

Important runtime configuration includes `DB_HOST`, `DB_USER`, `DB_NAME`, `JWT_SECRET`, server host/port, database pool limits, token durations, and log level. Docker Compose provides MySQL only; schema creation is handled separately through the migration commands in the `Makefile`.

## Extension Points

```mermaid
flowchart LR
    App["Current app"]
    Identity["Identity bounded context"]
    FutureModules["Future bounded contexts"]
    Infra["Shared platform infrastructure"]
    Events["Domain events"]
    Async["Future Redis/Kafka/tracing"]

    App --> Identity
    App --> FutureModules
    Identity --> Events
    FutureModules --> Infra
    Identity --> Infra
    Infra -. "placeholders exist" .-> Async
```

The bootstrap structs already leave room for Redis, Kafka, and tracing, but those are placeholders. The next natural extensions would be adding event publication for `UserRegistered`, adding authenticated routes that verify JWTs, and expanding the web app to include login and token-aware API calls.
