# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.14.x  | ✅ Active |
| < 0.14  | ❌ No longer supported |

## Reporting a Vulnerability

**Please do NOT open a public issue** for security vulnerabilities.

### How to Report

Send an email to the project maintainers with:

1. **Description** of the vulnerability
2. **Steps to reproduce** (if applicable)
3. **Impact assessment** — what could an attacker achieve?
4. **Suggested fix** (if you have one)

### What to Expect

| Timeframe | Action |
|-----------|--------|
| **48 hours** | Acknowledgement of your report |
| **7 days** | Initial assessment and severity classification |
| **30 days** | Fix deployed or mitigation plan communicated |

### Scope

The following are considered in-scope for security reports:

- Authentication and authorization bypasses
- SQL injection or other injection vulnerabilities
- Exposure of sensitive data (API keys, credentials, PII)
- Insecure direct object references (IDOR)
- Cross-site scripting (XSS) in API responses
- Denial of service through API abuse

### Out of Scope

- Vulnerabilities in third-party dependencies (report to the upstream project)
- Social engineering attacks
- Physical security

## Security Best Practices

This project implements the following security measures:

- **JWT validation** via Keycloak JWKS endpoint
- **Role-based access control** (`RequireRole` middleware)
- **Input sanitization** on all mutated resources (POST/PUT)
- **Secret scanning** — credentials and API keys are excluded via `.gitignore`
- **SonarCloud** analysis on every push for vulnerability detection
- **Dependency auditing** via `go mod tidy` + SonarCloud dependency checks

## Acknowledgments

We appreciate responsible disclosure and will credit reporters (with permission) in our release notes.
