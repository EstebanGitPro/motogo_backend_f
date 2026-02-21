# Contributing to MotoGo Backend

Thank you for considering contributing to MotoGo! This document outlines our development workflow and standards.

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Commit Convention](#commit-convention)
- [Pull Request Process](#pull-request-process)
- [Testing Requirements](#testing-requirements)

---

## Getting Started

1. Fork and clone the repository
2. Follow the [Getting Started guide](README.md#-getting-started) to set up your environment
3. Install Git hooks:

```bash
make setup-hooks-all
```

This installs:
- **pre-commit**: Runs `fmt` + `vet` + `staticcheck` + `test-short`
- **pre-push**: Enforces the 65% average coverage threshold

---

## Development Workflow

We follow **Git Flow** with the following branch naming conventions:

| Branch Type | Pattern | Example |
|-------------|---------|---------|
| Feature | `feature/<description>` | `feature/ratings` |
| Bugfix | `fix/<description>` | `fix/branch-image-deletion` |
| Hotfix | `hotfix/<description>` | `hotfix/auth-token-expiry` |
| Release | `release/v<X.Y.Z>` | `release/v0.14.0` |

**Process:**

1. Create your branch from `develop`
2. Make your changes in focused, logical commits
3. Push and open a Pull Request to `develop`
4. Ensure CI checks pass (lint, test, coverage, SonarCloud)
5. Request review

---

## Code Standards

### Architecture

Follow the **Clean Architecture** layering:

- **`core/ports/`** — Define interfaces (repository contracts)
- **`core/interactor/`** — Business logic (must NOT import `gin`, `sql`, or infrastructure)
- **`handlers/`** — HTTP layer only (request parsing → interactor call → response mapping)
- **`platform/`** — Infrastructure implementations (database, Firebase, Keycloak)

### Linting

The project uses `golangci-lint` with a curated rule set (see `.golangci.yml`):

```bash
# Quick lint (vet + staticcheck)
make lint

# Full lint (golangci-lint — recommended)
make lint-full
```

### Formatting

```bash
make fmt
```

All code must pass `gofmt` and `goimports` formatting.

### Key Rules

- **Input Sanitization**: All string inputs in `POST`/`PUT` handlers must be trimmed
- **HATEOAS**: All API responses must include hypermedia links (Richardson Maturity Level 3)
- **System Messages**: Use the centralized `message_cache` — never hardcode user-facing strings
- **Logging**: Use `log/slog` with structured fields. Logs must be in **Spanish**
- **Timezone**: All date operations use `America/Bogota` via `ParseInLocation`

---

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]
[optional footer]
```

### Types

| Type | Purpose |
|------|---------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `style` | Code formatting (no logic changes) |
| `refactor` | Code restructuring (no feature/fix) |
| `test` | Adding or updating tests |
| `chore` | Build, CI, tooling changes |
| `perf` | Performance improvements |

### Examples

```
feat(ratings): add service review endpoint
fix(branch): resolve image deletion on profile update
test(handlers): add unit tests for motorcycle controller
docs(readme): update installation instructions
```

---

## Pull Request Process

1. **Title**: Use the conventional commit format
2. **Description**: Explain **what** and **why**, not just **how**
3. **Size**: Keep PRs focused. Split large changes into logical commits
4. **Checklist** before opening:
   - [ ] Code compiles: `go build ./...`
   - [ ] All tests pass: `make test`
   - [ ] Lint passes: `make lint`
   - [ ] Coverage meets threshold: `make coverage-check`
   - [ ] No new SonarCloud issues introduced

---

## Testing Requirements

### Minimum Coverage

The project enforces a **65% minimum average coverage** across all packages.

### Testing Layers

| Layer | Framework | Pattern |
|-------|-----------|---------|
| **Interactors** | `testify/mock` | Mock port interfaces, test business logic |
| **Handlers** | `testify/mock` + `httptest` | Mock interactors, test HTTP behavior |
| **Repositories** | `go-sqlmock` | Mock database driver, test SQL generation |
| **Middleware** | `httptest` | Test auth, CORS, and request tracing |

### Running Tests

```bash
# All tests (verbose)
make test

# Quick run (no verbose)
make test-short

# With coverage report
make coverage
```

---

## Questions?

If you have questions about contributing, feel free to open an issue for discussion.
