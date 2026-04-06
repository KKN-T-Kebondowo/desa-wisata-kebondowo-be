# Desa Wisata Kebondowo - Backend API

REST API backend for the Desa Wisata Kebondowo (Kebondowo Tourism Village) management system. Built with Go, Gin, GORM, and PostgreSQL.

This service powers both the public-facing tourism website and the admin dashboard, providing endpoints for managing tourism spots, UMKM (small businesses), articles, galleries, and visitor analytics.

## Table of Contents

- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [Database Setup](#database-setup)
  - [Running the Server](#running-the-server)
  - [Running with Docker](#running-with-docker)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [Authentication](#authentication)
  - [Users](#users)
  - [Roles](#roles)
  - [Tourism](#tourism)
  - [Tourism Pictures](#tourism-pictures)
  - [UMKM](#umkm)
  - [UMKM Pictures](#umkm-pictures)
  - [Articles](#articles)
  - [Galleries](#galleries)
  - [Dashboard](#dashboard)
- [Authentication Flow](#authentication-flow)
- [Code Best Practices](#code-best-practices)
  - [Security](#security)
  - [Database Access](#database-access)
  - [Error Handling](#error-handling)
  - [Input Validation](#input-validation)
  - [Concurrency Safety](#concurrency-safety)
- [Testing](#testing)
  - [Test Architecture](#test-architecture)
  - [Running Tests](#running-tests)
  - [Test Categories](#test-categories)
- [Deployment](#deployment)

---

## Architecture

```
Client (Browser / Dashboard)
        |
        v
   Gin HTTP Server
        |
        +-- CORS Middleware
        +-- Auth Middleware (JWT RS256)
        |
        v
   Controllers (request/response handling)
        |
        v
   GORM ORM (query building, validation)
        |
        v
   PostgreSQL Database
```

**Key technology choices:**

| Component       | Technology                  | Why                                           |
|-----------------|-----------------------------|-----------------------------------------------|
| HTTP framework  | Gin v1.9                    | High performance, middleware ecosystem        |
| ORM             | GORM v1.25                  | Auto-migration, relationship mapping          |
| Database        | PostgreSQL                  | Relational data with spatial fields (lat/lng) |
| Authentication  | JWT with RSA (RS256)        | Asymmetric signing, stateless auth            |
| Password hashing| bcrypt (cost 10)            | Industry-standard, adaptive hashing           |
| Validation      | go-playground/validator v10 | Struct tag-based input validation              |

---

## Project Structure

```
.
├── main.go                      # Application entry point
├── go.mod / go.sum              # Go module dependencies
├── Dockerfile                   # Container build file
├── .env                         # Local environment config (not committed)
├── .env.copy                    # Environment variable template
│
├── controllers/                 # HTTP request handlers
│   ├── auth.controller.go       # Login, register, refresh, logout
│   ├── user.controllers.go      # Get current user profile
│   ├── role.controller.go       # Role CRUD
│   ├── tourism.controller.go    # Tourism spot CRUD + visitor tracking
│   ├── tourismpictures.controller.go
│   ├── gallery.controller.go    # Gallery CRUD with pagination
│   ├── article.controller.go    # Article CRUD with pagination + visitors
│   ├── umkm.controller.go       # UMKM CRUD + visitor tracking
│   ├── umkmpictures.controller.go
│   └── dashboard.controller.go  # Aggregated statistics
│
├── models/                      # Database models and DTOs
│   ├── user.model.go            # User, SignUpInput, SignInInput
│   ├── role.model.go            # Role, RoleInput
│   ├── tourism.model.go         # Tourism, TourismInput, TourismUpdate
│   ├── tourismpicture.model.go  # TourismPicture
│   ├── gallery.model.go         # Gallery, GalleryInput
│   ├── article.model.go         # Article, ArticleInput, ArticleUpdate
│   ├── umkm.model.go            # UMKM, UMKMInput, UMKMUpdate
│   ├── umkmpicture.model.go     # UMKMPicture
│   └── dashboard.model.go       # DashboardResponse
│
├── routes/                      # Route registration (maps paths to controllers)
│   ├── auth.routes.go
│   ├── user.routes.go
│   ├── role.routes.go
│   ├── tourism.routes.go
│   ├── tourismpicture.routes.go
│   ├── gallery.routes.go
│   ├── article.routes.go
│   ├── umkm.routes.go
│   ├── umkmpicture.routes.go
│   └── dashboard.routes.go
│
├── middleware/
│   └── deserialize-user.go      # JWT authentication middleware
│
├── initializers/
│   ├── loadEnv.go               # Environment config loader
│   └── connectDB.go             # PostgreSQL connection setup
│
├── utils/
│   ├── password.go              # bcrypt hash/verify helpers
│   └── token.go                 # RSA JWT token create/validate
│
├── migrate/
│   └── migrate.go               # Database migration script
│
└── e2e/                         # End-to-end tests
    ├── setup_test.go            # Test infrastructure and helpers
    ├── auth_test.go             # Authentication flow tests
    ├── crud_test.go             # CRUD operation tests
    └── security_test.go         # Security and edge case tests
```

---

## Getting Started

### Prerequisites

- **Go** 1.20 or later
- **PostgreSQL** 12 or later
- **Git**

### Environment Variables

Copy the template and fill in your values:

```bash
cp .env.copy .env
```

| Variable                         | Description                          | Example                |
|----------------------------------|--------------------------------------|------------------------|
| `POSTGRES_HOST`                  | PostgreSQL host                      | `localhost`            |
| `POSTGRES_USER`                  | PostgreSQL username                  | `postgres`             |
| `POSTGRES_PASSWORD`              | PostgreSQL password                  | `yourpassword`         |
| `POSTGRES_DB`                    | Database name                        | `kebondowo`            |
| `POSTGRES_PORT`                  | PostgreSQL port                      | `5432`                 |
| `PORT`                           | HTTP server port                     | `8080`                 |
| `CLIENT_ORIGIN`                  | Frontend URL (for CORS)             | `http://localhost:3000`|
| `ACCESS_TOKEN_PRIVATE_KEY`       | Base64-encoded RSA private key (PEM) | *(see below)*          |
| `ACCESS_TOKEN_PUBLIC_KEY`        | Base64-encoded RSA public key (PEM)  | *(see below)*          |
| `REFRESH_TOKEN_PRIVATE_KEY`      | Base64-encoded RSA private key (PEM) | *(see below)*          |
| `REFRESH_TOKEN_PUBLIC_KEY`       | Base64-encoded RSA public key (PEM)  | *(see below)*          |
| `ACCESS_TOKEN_EXPIRED_IN_HOURS`  | Access token TTL in hours            | `1`                    |
| `REFRESH_TOKEN_EXPIRED_IN_DAYS`  | Refresh token TTL in days            | `7`                    |

**Generating RSA key pairs:**

```bash
# Generate a 2048-bit RSA private key
openssl genrsa -out private.pem 2048

# Extract the public key
openssl rsa -in private.pem -pubout -out public.pem

# Base64-encode for .env (single line, no newlines)
cat private.pem | base64 | tr -d '\n'
cat public.pem | base64 | tr -d '\n'
```

You need two separate key pairs: one for access tokens and one for refresh tokens.

### Database Setup

```bash
# Create the database
createdb kebondowo

# Run migrations (creates all tables)
go run migrate/migrate.go
```

### Running the Server

```bash
# Install dependencies
go mod tidy

# Run the server
go run main.go
```

The server starts at `http://localhost:8080` (or the port specified in `PORT`).

Verify it's running:

```bash
curl http://localhost:8080/api/healthchecker
# {"message":"Welcome to Golang","status":"success"}
```

### Running with Docker

```bash
# Build the image
docker build -t kebondowo-api .

# Run the container (make sure PostgreSQL is accessible)
docker run -p 8080:8080 --env-file .env kebondowo-api
```

---

## API Endpoints

Base URL: `/api`

All write operations (POST, PUT, DELETE) require authentication via Bearer token unless noted otherwise. All GET endpoints are public.

### Health Check

| Method | Path                 | Auth | Description          |
|--------|----------------------|------|----------------------|
| GET    | `/api/healthchecker` | No   | Server health status |

**Response:**
```json
{
  "status": "success",
  "message": "Welcome to Golang"
}
```

### Authentication

| Method | Path                  | Auth | Description              |
|--------|-----------------------|------|--------------------------|
| POST   | `/api/auth/login`     | No   | Login with credentials   |
| POST   | `/api/auth/register`  | Yes  | Register a new user      |
| POST   | `/api/auth/refresh`   | No   | Refresh access token     |
| GET    | `/api/auth/logout`    | Yes  | Logout (clear cookies)   |

**POST /api/auth/login**
```json
// Request
{
  "username": "admin",
  "password": "admin1234"
}

// Response 200
{
  "status": "success",
  "access_token": "eyJhbGciOi...",
  "refresh_token": "eyJhbGciOi..."
}
```

**POST /api/auth/register** (requires Bearer token)
```json
// Request
{
  "username": "newuser",
  "password": "securepass123",
  "roleid": 1
}

// Response 201
{
  "status": "success",
  "data": {
    "user": {
      "id": 2,
      "username": "newuser",
      "role": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

**POST /api/auth/refresh**
```json
// Request
{
  "refresh_token": "eyJhbGciOi..."
}

// Response 200
{
  "access_token": "eyJhbGciOi..."
}
```

### Users

| Method | Path             | Auth | Description           |
|--------|------------------|------|-----------------------|
| GET    | `/api/users/me`  | Yes  | Get current user info |

**Response 200:**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "role": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### Roles

| Method | Path              | Auth | Description      |
|--------|-------------------|------|------------------|
| GET    | `/api/roles/`     | No   | List all roles   |
| GET    | `/api/roles/:id`  | No   | Get role by ID   |
| POST   | `/api/roles/`     | Yes  | Create a role    |
| PUT    | `/api/roles/:id`  | Yes  | Update a role    |
| DELETE | `/api/roles/:id`  | Yes  | Delete a role    |

**POST /api/roles/** (requires Bearer token)
```json
// Request
{ "name": "editor" }

// Response 201
{
  "status": "success",
  "data": { "role": { "id": 3, "name": "editor" } }
}
```

### Tourism

| Method | Path                  | Auth | Description             |
|--------|-----------------------|------|-------------------------|
| GET    | `/api/tourisms/`      | No   | List all tourism spots  |
| GET    | `/api/tourisms/:slug` | No   | Get by slug (increments visitor count) |
| POST   | `/api/tourisms/`      | Yes  | Create tourism spot     |
| PUT    | `/api/tourisms/:id`   | Yes  | Update tourism spot     |
| DELETE | `/api/tourisms/:id`   | Yes  | Delete tourism spot (cascades to pictures) |

**POST /api/tourisms/** (requires Bearer token)
```json
// Request
{
  "title": "Rawa Pening Lake",
  "slug": "rawa-pening-lake",
  "cover_picture_url": "https://example.com/lake.jpg",
  "description": "A beautiful lake in Central Java",
  "latitude": -7.2906,
  "longitude": 110.3987,
  "pictures": [
    { "picture_url": "https://example.com/pic1.jpg", "caption": "Sunrise view" },
    { "picture_url": "https://example.com/pic2.jpg", "caption": "Boat ride" }
  ]
}

// Response 200
{
  "status": "success",
  "data": { "tourism": { ... } }
}
```

**GET /api/tourisms/:slug** - Automatically increments the visitor counter on each request.

### Tourism Pictures

| Method | Path                          | Auth | Description              |
|--------|-------------------------------|------|--------------------------|
| GET    | `/api/tourism-pictures/`      | No   | List all tourism pictures|
| GET    | `/api/tourism-pictures/:id`   | No   | Get picture by ID        |
| POST   | `/api/tourism-pictures/`      | Yes  | Add a tourism picture    |
| DELETE | `/api/tourism-pictures/:id`   | Yes  | Delete a tourism picture |

**POST /api/tourism-pictures/** (requires Bearer token)
```json
// Request
{
  "picture_url": "https://example.com/pic.jpg",
  "caption": "Beautiful scenery",
  "tourism_id": 1
}
```

### UMKM

| Method | Path               | Auth | Description            |
|--------|--------------------|------|------------------------|
| GET    | `/api/umkms/`      | No   | List all UMKM          |
| GET    | `/api/umkms/:slug` | No   | Get by slug (increments visitor count) |
| POST   | `/api/umkms/`      | Yes  | Create UMKM entry      |
| PUT    | `/api/umkms/:id`   | Yes  | Update UMKM entry      |
| DELETE | `/api/umkms/:id`   | Yes  | Delete UMKM (cascades to pictures) |

**POST /api/umkms/** (requires Bearer token)
```json
// Request
{
  "title": "Batik Kebondowo",
  "slug": "batik-kebondowo",
  "cover_picture_url": "https://example.com/batik.jpg",
  "description": "Traditional batik from Kebondowo village",
  "latitude": -7.2910,
  "longitude": 110.3990,
  "contact": "081234567890",
  "contact_name": "Pak Budi",
  "pictures": [
    { "picture_url": "https://example.com/batik1.jpg", "caption": "Process" }
  ]
}
```

### UMKM Pictures

| Method | Path                       | Auth | Description           |
|--------|----------------------------|------|-----------------------|
| GET    | `/api/umkm-pictures/`      | No   | List all UMKM pictures|
| GET    | `/api/umkm-pictures/:id`   | No   | Get picture by ID     |
| POST   | `/api/umkm-pictures/`      | Yes  | Add a UMKM picture    |
| DELETE | `/api/umkm-pictures/:id`   | Yes  | Delete a UMKM picture |

### Articles

| Method | Path                  | Auth | Description            |
|--------|-----------------------|------|------------------------|
| GET    | `/api/articles/`      | No   | List articles (paginated) |
| GET    | `/api/articles/:slug` | No   | Get by slug (increments visitor count) |
| POST   | `/api/articles/`      | Yes  | Create article         |
| PUT    | `/api/articles/:id`   | Yes  | Update article         |
| DELETE | `/api/articles/:id`   | Yes  | Delete article         |

**GET /api/articles/?limit=10&offset=0&sortby=created_at&orderedby=desc**

Query parameters:

| Param       | Default      | Allowed Values                          |
|-------------|-------------|------------------------------------------|
| `limit`     | `20`        | Any positive integer                     |
| `offset`    | `0`         | Any non-negative integer                 |
| `sortby`    | `created_at`| `created_at`, `updated_at`, `id`, `title`|
| `orderedby` | `desc`      | `asc`, `desc`                            |

**Response 200:**
```json
{
  "status": "success",
  "data": {
    "articles": [ ... ],
    "meta": {
      "limit": 10,
      "offset": 0,
      "total": 42
    }
  }
}
```

**POST /api/articles/** (requires Bearer token)
```json
// Request
{
  "title": "Festival Desa Kebondowo 2024",
  "slug": "festival-desa-kebondowo-2024",
  "author": "Admin",
  "content": "The annual village festival...",
  "picture_url": "https://example.com/festival.jpg"
}
```

### Galleries

| Method | Path                 | Auth | Description              |
|--------|----------------------|------|--------------------------|
| GET    | `/api/galleries/`    | No   | List galleries (paginated) |
| GET    | `/api/galleries/:id` | No   | Get gallery by ID        |
| POST   | `/api/galleries/`    | Yes  | Create gallery entry     |
| DELETE | `/api/galleries/:id` | Yes  | Delete gallery entry     |

**GET /api/galleries/?limit=20&offset=0&sortby=created_at&orderedby=desc**

Query parameters:

| Param       | Default      | Allowed Values                  |
|-------------|-------------|----------------------------------|
| `limit`     | `20`        | Any positive integer             |
| `offset`    | `0`         | Any non-negative integer         |
| `sortby`    | `created_at`| `created_at`, `updated_at`, `id` |
| `orderedby` | `desc`      | `asc`, `desc`                    |

**POST /api/galleries/** (requires Bearer token)
```json
// Request
{
  "picture_url": "https://example.com/photo.jpg",
  "caption": "Village morning view"
}
```

### Dashboard

| Method | Path               | Auth | Description                 |
|--------|--------------------|------|-----------------------------|
| GET    | `/api/dashboard/`  | No   | Aggregated site statistics  |

**Response 200:**
```json
{
  "status": "success",
  "data": {
    "total_article": 15,
    "total_gallery": 42,
    "total_tourism": 8,
    "total_visitor": 1234,
    "total_umkm": 12,
    "article_per_month": [0, 2, 1, 3, 0, 0, 5, 1, 0, 2, 0, 1],
    "gallery_per_month": [1, 0, 3, 2, 1, 0, 4, 2, 1, 0, 3, 1],
    "tourism_per_month": [0, 0, 1, 0, 1, 0, 2, 1, 0, 1, 1, 1],
    "umkm_per_month": [0, 1, 0, 2, 0, 1, 1, 2, 0, 1, 2, 2]
  }
}
```

The `*_per_month` arrays contain 12 elements representing monthly creation counts, with the current month as the last element.

---

## Authentication Flow

```
1. Client sends POST /api/auth/login with username + password
2. Server verifies password against bcrypt hash
3. Server generates:
   - Access token  (RS256, expires in ACCESS_TOKEN_EXPIRED_IN_HOURS)
   - Refresh token (RS256, expires in REFRESH_TOKEN_EXPIRED_IN_DAYS)
4. Client stores tokens and sends:
   Authorization: Bearer <access_token>
   on all authenticated requests

5. When access token expires:
   POST /api/auth/refresh with { "refresh_token": "..." }
   Server returns a new access_token

6. Middleware flow for protected routes:
   Request -> Extract token from Authorization header (or cookie)
           -> Validate RSA signature
           -> Look up user by token's "sub" claim
           -> Inject user into Gin context
           -> Pass to controller
```

**Token claims:**

| Claim      | Description            |
|------------|------------------------|
| `sub`      | User ID                |
| `username` | Username               |
| `roleid`   | User's role ID         |
| `exp`      | Expiration timestamp   |
| `iat`      | Issued-at timestamp    |
| `nbf`      | Not-before timestamp   |

---

## Code Best Practices

### Security

**SQL Injection Prevention**

All database queries use GORM's parameterized queries with `?` placeholders. User-supplied sort parameters are validated against explicit whitelists:

```go
// controllers/gallery.controller.go
allowedSortFields := map[string]bool{
    "created_at": true,
    "updated_at": true,
    "id":         true,
}
if !allowedSortFields[sortBy] {
    sortBy = "created_at" // safe default
}
```

This approach prevents SQL injection through ORDER BY clauses, which parameterized queries alone cannot protect against.

**Password Security**

Passwords are hashed with bcrypt at the default cost factor (10), providing adaptive resistance to brute-force attacks:

```go
// utils/password.go
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedPassword), err
}
```

**JWT with RSA (Asymmetric Signing)**

Using RS256 instead of HS256 means the private key (used for signing) and public key (used for verification) are separate. This enables scenarios where verification happens on different services without exposing the signing key. The middleware also validates the signing method to prevent algorithm confusion attacks:

```go
// utils/token.go - validates RSA signing method
if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
    return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
}
```

**Registration Protection**

The `/api/auth/register` endpoint is behind the auth middleware, so only authenticated users (admins) can create new accounts. This prevents unauthorized user creation.

**CORS Configuration**

CORS is configured with an explicit allowlist of origins rather than a wildcard, and credentials are enabled:

```go
corsConfig.AllowOrigins = []string{
    "http://localhost:3000",
    "https://dashboard-desa-wisata-kebondowo.vercel.app",
    "https://desa-wisata-kebondowo-dashboard.vercel.app",
    "https://dashboard.kebondowo.com",
}
corsConfig.AllowCredentials = true
```

**Generic Error Messages**

Login failures return `"Invalid username or password"` regardless of whether the username exists, preventing user enumeration attacks.

### Database Access

**Atomic Updates for Counters**

Visitor counters use `gorm.Expr` to perform atomic increments at the database level, avoiding read-modify-write race conditions:

```go
// controllers/tourism.controller.go
initializers.DB.Model(&tourism).
    Update("visitor", gorm.Expr("visitor + 1"))
```

**Manual Cascade Deletes**

Since GORM's `AutoMigrate` does not set up foreign key cascade rules by default, parent-child deletes are handled explicitly:

```go
// controllers/umkm.controller.go
initializers.DB.Where("umkm_id = ?", id).Delete(&models.UMKMPicture{})
initializers.DB.Delete(&umkm)
```

**Slug Deduplication**

When creating tourism spots or UMKMs, duplicate slugs are automatically resolved by appending a numeric suffix:

```go
// If "rawa-pening" exists, the next one becomes "rawa-pening-1"
var count int64
initializers.DB.Model(&models.Tourism{}).Where("slug = ?", slug).Count(&count)
if count > 0 {
    slug = fmt.Sprintf("%s-%d", slug, count)
}
```

### Error Handling

- `LoadConfig` returns errors instead of calling `log.Fatal`, allowing callers to handle failures gracefully
- All `strconv.Atoi` results are checked before use
- Database migration errors are checked individually per model
- Type assertions on Gin context values use the safe two-value form (`value, ok := ...`)

### Input Validation

Gin's struct binding tags enforce required fields and constraints at the request parsing layer:

```go
type SignUpInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=8"`
    RoleID   uint32 `json:"roleid"   binding:"required"`
}
```

Invalid requests are rejected with `400 Bad Request` before reaching any business logic.

### Concurrency Safety

- Visitor counters use database-level atomic increments (`SET visitor = visitor + 1`)
- No global mutable state beyond `initializers.DB` (which is set once at startup)
- GORM's connection pool handles concurrent database access

---

## Testing

### Test Architecture

Tests are **end-to-end tests** that exercise the full application stack: HTTP request -> Gin router -> middleware -> controller -> GORM -> PostgreSQL. No mocks are used.

```
Test Function
    |
    v
HTTP Client (net/http)
    |
    v
httptest.Server (real Gin server)
    |
    v
Middleware + Controllers
    |
    v
GORM -> PostgreSQL (kebondowo_test database)
```

### Running Tests

**1. Create the test database:**

```bash
createdb kebondowo_test
```

**2. Ensure your `.env` file is configured** (tests use the same `.env` but override the database name to `kebondowo_test`).

**3. Run all tests:**

```bash
go test ./e2e/ -v -count=1
```

**4. Run with race detection:**

```bash
go test ./e2e/ -v -race -count=1
```

**5. Run a specific test:**

```bash
go test ./e2e/ -v -run TestAuthLogin -count=1
```

**6. Run a test category:**

```bash
# Auth tests only
go test ./e2e/ -v -run "TestAuth|TestRegister|TestRefresh|TestGetMe|TestLogout" -count=1

# Security tests only
go test ./e2e/ -v -run "TestSQL|TestVisitor|TestDelete|TestInvalid|TestSlug" -count=1
```

The `-count=1` flag disables test caching, ensuring tests always hit the real database.

### Test Categories

**Auth Tests** (`e2e/auth_test.go`) -- 10 tests

| Test                            | What it verifies                                     |
|---------------------------------|------------------------------------------------------|
| `TestHealthCheck`               | Health endpoint returns 200                          |
| `TestAuthLogin`                 | Valid credentials return access + refresh tokens     |
| `TestAuthLoginWrongPassword`    | Wrong password returns 400                           |
| `TestAuthLoginNonExistentUser`  | Non-existent user returns 400                        |
| `TestAuthLoginMissingFields`    | Missing required fields returns 400                  |
| `TestRegisterRequiresAuth`      | Register without token returns 401                   |
| `TestRegisterWithAuth`          | Register with token creates user (201)               |
| `TestRefreshToken`              | Valid refresh token returns new access token          |
| `TestRefreshTokenInvalid`       | Invalid refresh token returns 403                    |
| `TestGetMe` / `TestGetMeNoAuth` | Protected profile returns user data or 401           |

**CRUD Tests** (`e2e/crud_test.go`) -- 12 tests

| Test                     | What it verifies                                         |
|--------------------------|----------------------------------------------------------|
| `TestRoleCRUD`           | Create, list, update, delete roles                       |
| `TestGalleryCRUD`        | Create, list with pagination meta, delete galleries      |
| `TestArticleCRUD`        | Full lifecycle with slug lookup and update                |
| `TestTourismCRUD`        | Create with pictures, slug lookup, update, cascade delete|
| `TestUMKMCRUD`           | Create with pictures, slug lookup, update, cascade delete|
| `TestTourismPictureCRUD` | Independent picture create/get/delete                    |
| `TestUMKMPictureCRUD`    | Independent picture create/get/delete                    |
| `TestDashboard`          | Dashboard endpoint returns valid statistics              |

**Security Tests** (`e2e/security_test.go`) -- 11 tests

| Test                              | What it verifies                                    |
|-----------------------------------|-----------------------------------------------------|
| `TestSQLInjectionSortBy`          | SQL in `sortby` param is sanitized                  |
| `TestSQLInjectionOrderedBy`       | SQL in `orderedby` param is sanitized               |
| `TestValidSortFields`             | All whitelisted sort fields work correctly           |
| `TestInvalidLimitOffset`          | Non-numeric pagination params return 400             |
| `TestVisitorCounterArticle`       | Visitor count increments on each article view        |
| `TestVisitorCounterTourism`       | Visitor count increments on each tourism view        |
| `TestDeleteCascadeTourism`        | Deleting tourism removes its pictures                |
| `TestDeleteCascadeUMKM`           | Deleting UMKM removes its pictures                   |
| `TestUMKMPictureDelete`           | UMKM picture delete actually removes the record      |
| `TestInvalidToken`                | Invalid JWT is rejected with 401                     |
| `TestSlugDeduplication`           | Duplicate slugs get auto-suffixed                    |

---

## Deployment

### Docker

```bash
docker build -t kebondowo-api .
docker run -p 8080:8080 --env-file .env kebondowo-api
```

### Production Considerations

- Set `Secure: true` on cookies when serving over HTTPS
- Use environment-specific RSA key pairs (do not reuse dev keys)
- Configure PostgreSQL connection pooling (`max_open_conns`, `max_idle_conns`)
- Add rate limiting middleware to prevent API abuse
- Set up database backups and monitoring
- Use a reverse proxy (nginx, Caddy) for TLS termination
