# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.14.0] - 2026-02-20

### Added
- Professional project documentation: README.md, CONTRIBUTING.md, SECURITY.md, CHANGELOG.md, LICENSE
- HATEOAS link audit and cleanup across all endpoints
- Service rating and review system (HU48)
- Completed service detail updates
- Load testing infrastructure with k6 and InfluxDB

### Changed
- Database connection pool configuration optimized for production workloads
- Input sanitization standardized across all POST/PUT endpoints (Release 28)
- Swagger documentation parity — all protected routes now show security requirements

### Fixed
- Branch profile image deletion on update (Firebase Storage orphan cleanup)
- JWT token expiration handling during load tests
- HATEOAS links pointing to non-existent endpoints removed

### Security
- Push protection for secrets enabled
- Production configuration files excluded from version control

---

_For versions prior to 0.14.0, see the [Git history](https://github.com/EstebanGitPro/motogo_backend_f/commits/main)._
