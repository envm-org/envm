# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v.0.0.1] - 2026-02-24

### Added
- Initial user migration and domain services.
- NPM package setup for CLI distribution.
- Automated releases with GoReleaser and GitHub Actions.
- `/api/v1` API routing and client updates.
- CLI implementation (`push`, `pull`, `load`, `env`, `init`, `login`, `project` commands).
- Docusaurus documentation setup.
- AES-256 encryption package for environment variables.
- Comprehensive database migrations using Postgres and SQLC.
- JWT-based authentication with refresh token functionality.
- Postman collection for the ENVM API.
- Pre-commit hooks for formatting and linting.
- Graceful server shutdown.
- GitHub Issue templates and Contribution documentation.

### Changed
- Refactored CLI build processes and updated CLI documentation in README.
- Updated Go module dependencies (e.g., `golang.org/x/crypto`).
- Standardized environment variables across configurations.
- Restructured application to use Chi router, dedicated handlers, and organized services.
- Removed organization concepts from the CLI.
- Simplified Go build commands for local development.
- Removed CI/CD actions for Docusaurus docs in favor of Vercel deployment.

### Fixed
- Fixed symlink compatibility for relative image paths in README files.
- Added MIT License explicitly to the project.
