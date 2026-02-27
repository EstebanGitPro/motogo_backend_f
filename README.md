# MotoGo Backend

[![Go Version](https://img.shields.io/badge/Go-1.25.2-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.14.0-blue)](CHANGELOG.md)
[![SonarCloud](https://img.shields.io/badge/SonarCloud-Analyzed-F3702A?logo=sonarcloud)](https://sonarcloud.io/project/overview?id=EstebanGitPro_motogo_backend_f)

> RESTful API backend for **MotoGo** — a mobile platform that helps motorcyclists find trusted workshops and services tailored to their registered vehicles, compare nearby options, and request quotes or diagnostics before visiting.

---

## 📖 Table of Contents

- [Vision](#-vision)
- [Key Features](#-key-features)
- [Technology Stack](#-technology-stack)
- [Architecture](#-architecture)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Configuration](#-configuration)
- [Available Commands](#-available-commands)
- [Observability](#-observability)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Port Reference](#-port-reference)
- [Contributing](#-contributing)
- [Security](#-security)
- [License](#-license)

---

## 🎯 Vision

For motorcyclists who don't have a reliable way to find trusted workshops offering technical services and consumables that match their needs, **MotoGo** is a mobile application that saves users both time and money when searching for motorcycle services.

Unlike Yelp, Google Maps, or Yellow Pages, MotoGo provides access to a **service catalog filtered by the vehicles each user has registered** in their profile, along with nearby service information and the ability to request quotes or approximate diagnostics before traveling to the location.

---

## ✨ Key Features

- **User Management** — Registration, authentication, email verification, and password reset flows
- **Identity Provider Integration** — Keycloak for secure identity and access management (OIDC/JWT)
- **Branch & Workshop Discovery** — Geospatial search with map-based nearby discovery
- **Motorcycle Management** — Vehicle registration, profile images, and evidence galleries
- **Service Catalog** — Global catalog with branch-specific service associations
- **Diagnostics & Services** — Two-sided diagnostic flow with WhatsApp integration
- **Rating System** — Service reviews and rating ranges
- **Dynamic Messaging** — Centralized, database-driven system messages with in-memory caching
- **HATEOAS API** — Richardson Maturity Model Level 3 hypermedia responses
- **Comprehensive Observability** — Prometheus metrics, Grafana dashboards, Loki log aggregation
- **API Documentation** — Auto-generated OpenAPI/Swagger specifications
- **Load Testing** — k6 scripts for performance and soak testing

---

## 🛠 Technology Stack

| Category | Technology |
|----------|-----------|
| **Language** | Go 1.25.2 |
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **Database** | MySQL 8.0 |
| **Identity Management** | [Keycloak](https://www.keycloak.org/) (OIDC + JWKS validation) |
| **File Storage** | Firebase Storage |
| **Email** | [Resend](https://resend.com/) |
| **Metrics** | Prometheus |
| **Dashboards** | Grafana |
| **Log Aggregation** | Loki + Promtail |
| **Structured Logging** | `log/slog` |
| **API Docs** | Swagger (swaggo) |
| **Code Quality** | SonarCloud, golangci-lint, staticcheck |
| **Load Testing** | k6 |
| **Containerization** | Docker + Docker Compose |

---

## 🏗 Architecture

MotoGo Backend follows a **Clean / Hexagonal Architecture** (Ports & Adapters) with 4 layers:

```
┌─────────────────────────────────────────────────┐
│                   handlers/                      │  ← Presentation Layer (HTTP controllers)
├─────────────────────────────────────────────────┤
│                  middleware/                      │  ← Cross-cutting concerns (auth, CORS, logging)
├─────────────────────────────────────────────────┤
│                    core/                         │
│   ┌──────────────┐  ┌────────────────────────┐  │
│   │    ports/     │  │     interactor/        │  │  ← Domain Layer (business logic + interfaces)
│   │ (interfaces)  │  │  (use cases / services)│  │
│   └──────────────┘  └────────────────────────┘  │
├─────────────────────────────────────────────────┤
│                  platform/                       │  ← Infrastructure (DB, Firebase, Keycloak, etc.)
└─────────────────────────────────────────────────┘
```

**Key principles:**
- Dependencies point **inward** — `platform` implements `core/ports` interfaces
- Business logic in `core/interactor` is **framework-agnostic**
- `handlers` only orchestrate HTTP request/response mapping
- `middleware` provides authentication, request tracing, and CORS

---

## 📁 Project Structure

```
motogo_backend_f/
├── cmd/
│   ├── main.go                 # Application entrypoint
│   └── dependency/             # Dependency injection wiring
├── config/                     # Configuration files (JSON) + config loader
├── core/
│   ├── ports/                  # Interfaces (repository, service contracts)
│   └── interactor/             # Use cases / business logic
├── handlers/                   # HTTP controllers (Gin handlers)
├── middleware/                  # Auth, CORS, request ID, metrics
├── server/                     # Server bootstrap and route registration
├── platform/                   # Infrastructure implementations
│   ├── cache/                  # In-memory caching (system messages)
│   ├── constants/              # Shared constants
│   ├── databases/              # MySQL repositories
│   ├── firebase/               # Firebase Storage adapter
│   ├── geocoding/              # Geocoding service (Google / Mapbox)
│   ├── grafana/                # Grafana dashboards & provisioning
│   ├── identity_provider/      # Keycloak integration
│   ├── jwt/                    # JWT / JWKS validation
│   ├── k6/                     # k6 load test scenarios
│   ├── log-rotator/            # Log rotation service
│   ├── logger/                 # Structured logging (slog)
│   ├── loki/                   # Loki configuration
│   ├── prometheus/             # Prometheus scrape configs
│   ├── promtail/               # Promtail log shipping config
│   ├── schema/                 # JSON Schema validation
│   └── swaggo/                 # Generated Swagger docs
├── mocks/                      # Testify mock implementations
├── public/                     # Static assets (email templates)
├── tools/                      # Utilities, keycloak themes, moderation
├── scripts/                    # Automation scripts (Docker setup, observability)
├── tests/                      # Integration / load tests (k6)
├── deploy/                     # Production deployment configs
│   ├── docker-compose.production.yml
│   ├── deploy.sh
│   ├── nginx-motogo.conf
│   └── .env.production.example
├── docs/                       # Swagger docs + Bruno API collections
├── docker-compose.mysql.yml    # Local MySQL (app + keycloak)
├── docker-compose.keycloak.yml # Local Keycloak
├── docker-compose.grafana.yml  # Observability stack
├── docker-compose.swagger-ui.yml
├── Containerfile               # Keycloak custom image (with MotoGo theme)
├── Makefile                    # Build, test, lint automation
├── .golangci.yml               # Linter configuration
└── sonar-project.properties    # SonarCloud integration
```

---

## 🚀 Getting Started

**Estimated time:** 15–20 minutes

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | ≥ 1.25 | Language runtime |
| **Docker** | ≥ 24.0 | Container runtime |
| **Docker Compose** | ≥ 2.20 | Multi-container orchestration |
| **Make** | Any | Build automation (pre-installed on macOS/Linux) |
| **Git** | Any | Version control |

### Step 1: Clone the Repository

```bash
git clone https://github.com/EstebanGitPro/motogo_backend_f.git
cd motogo_backend_f
```

### Step 2: Start MySQL Databases

The project uses **two MySQL instances**: one for the application and one for Keycloak.

```bash
docker compose -f docker-compose.mysql.yml up -d
```

**Verification:**
```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
# Expected: motogo-mysql-app (healthy), motogo-mysql-keycloak (healthy)
```

> **Wait** for both containers to report `(healthy)` before proceeding (~30 seconds).

### Step 3: Start Keycloak

```bash
docker compose -f docker-compose.keycloak.yml up -d
```

**Verification:** Open [http://localhost:8080](http://localhost:8080) — you should see the Keycloak admin console.

> Default credentials: `motogo-admin` / `12345678` (local development only).

### Step 4: Configure the Application

#### 4a. Environment Variables

The `.env` file at the project root configures all external services and secrets. Copy and edit:

```bash
cp .env.example .env
```

The `.env.example` contains the following sections:

| Section | Key Variables | Purpose |
|---------|-------------|----------|
| **Database** | `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME` | MySQL connection |
| **Server** | `SERVER_HOST`, `SERVER_PORT` | API server binding |
| **Resend** | `RESEND_API_KEY`, `RESEND_FROM_EMAIL` | Email delivery (password reset, verification) |
| **Keycloak** | `KEYCLOAK_SERVER_URL`, `KEYCLOAK_REALM`, `KEYCLOAK_CLIENT_SECRET` | Identity provider |
| **Firebase** | `FIREBASE_CREDENTIALS_PATH` | File storage service account path |
| **ID Encoder** | `ID_ENCODER_SECRET`, `ID_ENCODER_MIN_LENGTH` | HashID encoding |
| **Geocoding** | `GEOCODING_PROVIDER`, `GEOCODING_API_KEY`, `GEOCODING_COUNTRY_CODE` | Address resolution (Google / Mapbox) |
| **Verification** | `VERIFICATION_BASE_URL` | Email verification callback URL |

> **Tip:** Environment variables override JSON config file settings. See `config/config.go` for details.

#### 4b. Application Config (JSON)

The application also supports JSON-based configuration for each environment:

| Environment | Config File | Trigger |
|-------------|------------|---------|
| **Local** (default) | `config/local-config.json` | `APP_ENV` unset or `local` |
| **Production** | `config/prod-config.json` | `APP_ENV=production` |

Edit `config/local-config.json` to set values not already overridden by `.env`.

#### 4c. Firebase Service Account

Place your Firebase service account key in the `config/` directory:

```bash
# Download from Firebase Console → Project Settings → Service Accounts
cp ~/Downloads/your-firebase-adminsdk.json config/serviceAccountKey.json
```

> ⚠️ This file is in `.gitignore` and must **never** be committed.

### Step 5: Install Dependencies

```bash
go mod tidy
```

### Step 6: Install Development Tools (Optional)

```bash
make install-tools
# Installs: staticcheck
# For golangci-lint: brew install golangci-lint
```

### Step 7: Run the Application

```bash
make run
```

Or directly:

```bash
go run ./cmd/api
```

### ✅ Verification

| Check | How to Verify |
|-------|---------------|
| **Server running** | Console shows `Listening on 0.0.0.0:8085` |
| **Health endpoint** | `curl http://localhost:8085/health` returns `200 OK` |
| **Metrics endpoint** | `curl http://localhost:8085/metrics` returns Prometheus data |
| **Swagger UI** | Open [http://localhost:8085/swagger/index.html](http://localhost:8085/swagger/index.html) |

---

## ⚙️ Configuration

The application supports multiple environments via JSON config files and environment variable overrides.

**Precedence:** Environment variables (`.env`) → JSON config file → defaults.

All sections of the configuration can be overridden via environment variables. The most commonly overridden variables are:

```bash
# Database
DB_HOST=localhost
DB_PORT=3309
DB_PASSWORD=your-password

# Server
SERVER_PORT=8085

# Keycloak
KEYCLOAK_SERVER_URL=https://auth.example.com
KEYCLOAK_CLIENT_SECRET=your-secret

# External services
RESEND_API_KEY=re_your_api_key
GEOCODING_API_KEY=your-google-maps-key
```

> See `.env.example` for the full list of available environment variables.

---

## 📋 Available Commands

All project automation is available through the `Makefile`:

| Command | Description |
|---------|-------------|
| `make help` | Show all available targets |
| `make run` | Run the application (`go run ./cmd/api`) |
| `make build` | Build binary to `bin/motogo-api` |
| `make test` | Run all tests with verbose output |
| `make test-short` | Run tests without verbose output |
| `make coverage` | Generate HTML coverage report |
| `make coverage-check` | Verify coverage meets 65% minimum threshold |
| `make lint` | Run `go vet` + `staticcheck` |
| `make lint-full` | Run `golangci-lint` with full rule set |
| `make fmt` | Format all Go source files |
| `make tidy` | Run `go mod tidy` |
| `make clean` | Remove build artifacts and coverage files |
| `make pre-commit` | Run full pre-commit suite (fmt + vet + staticcheck + test) |
| `make pre-push` | Run tests with average coverage threshold enforcement |
| `make setup-hooks` | Install Git pre-commit hook |
| `make setup-hooks-all` | Install both pre-commit and pre-push hooks |
| `make install-tools` | Install `staticcheck` and prompt for `golangci-lint` |

---

## 📊 Observability

### Start the Full Observability Stack

```bash
docker compose -f docker-compose.grafana.yml up -d
```

This starts **5 services**:

| Service | Port | Purpose |
|---------|------|---------|
| **Grafana** | [localhost:3000](http://localhost:3000) | Dashboards and visualization |
| **Prometheus** | [localhost:9091](http://localhost:9091) | Metrics scraping and storage |
| **Loki** | localhost:3100 | Log aggregation |
| **Promtail** | — | Log shipping agent |
| **Log Rotator** | — | Automatic log rotation (7-day retention) |

> Grafana default credentials: `admin` / `admin`

### Quick Start Script

For first-time setup with automatic provisioning:

```bash
./scripts/quick-start-observability.sh
```

---

## 📚 API Documentation

### Swagger UI (Integrated)

The API documentation is served directly from the running application:

```
http://localhost:8085/swagger/index.html
```

### Swagger UI (Standalone Container)

For a standalone Swagger UI with hot-reload:

```bash
docker compose -f docker-compose.swagger-ui.yml up -d
# Open http://localhost:3001
```

---

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Coverage Report

```bash
make coverage
# Generates: coverage.html (open in browser)
# Prints: total coverage percentage
```

### Coverage Threshold

The project enforces a **65% minimum coverage** threshold:

```bash
make coverage-check
# Fails CI if coverage < 65%
```

### Load Testing (k6)

```bash
# Start the k6 + InfluxDB stack
docker compose -f docker-compose.k6-influxdb.yml up -d

# Run load tests (scripts in tests/ directory)
```

---

## 🚢 Deployment

Production deployment files are in the `deploy/` directory:

| File | Purpose |
|------|---------|
| `docker-compose.production.yml` | Production MySQL (app + keycloak) + Keycloak |
| `.env.production.example` | Template for production secrets |
| `deploy.sh` | Automated deployment script |
| `nginx-motogo.conf` | Nginx reverse proxy configuration |
| `motogo-api.service` | systemd service file for the Go binary |
| `ssh.sh` | SSH connection helper |

### Quick Deploy

```bash
# 1. Copy and configure production secrets
cp deploy/.env.production.example deploy/.env.production

# 2. Generate secure passwords
openssl rand -base64 24  # Use for each password field

# 3. Run the deployment script
./deploy/deploy.sh
```

> See `deploy/.env.production.example` for all required variables.

---

## 🔌 Port Reference

| Port | Service | Environment |
|------|---------|-------------|
| `8085` | MotoGo API | Local + Production |
| `3306` | MySQL (Application) | Local |
| `3309` | MySQL (Keycloak) | Local |
| `8080` | Keycloak (HTTP) | Local + Production |
| `8443` | Keycloak (HTTPS) | Production |
| `9000` | Keycloak (Health/Metrics) | Local + Production |
| `3000` | Grafana | Local |
| `9091` | Prometheus | Local |
| `3100` | Loki | Local |
| `3001` | Swagger UI (standalone) | Local |

---

## 🤝 Contributing

We welcome contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on:

- Git Flow branching strategy
- Conventional commit format
- Code quality standards (golangci-lint, pre-commit hooks)
- Pull request process

---

## 🛡️ Security

To report a vulnerability, please read our [Security Policy](SECURITY.md).

**Do not** open a public issue for security vulnerabilities.

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

<p align="center">
  <sub>Built with ❤️ by the MotoGo team</sub>
</p>
