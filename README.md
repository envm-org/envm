# ENVM - Environment Variable Management

<p align="center">
  <img src="../assets/logo.png" alt="ENVM Logo" width="200"/>
</p>

<h3 align="center">Secure, Simple, Synchronized Environment Management</h3>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"/></a>
  <a href="https://golang.org/dl/"><img src="https://img.shields.io/badge/Go-1.25.5-blue.svg" alt="Go Version"/></a>
  <a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Docker-Ready-blue.svg" alt="Docker"/></a>
</p>

<p align="center">
  <a href="#-features">Features</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-documentation">Documentation</a> •
  <a href="#-contributing">Contributing</a>
</p>

---

## 📖 About

**ENVM** is a powerful environment variable management and synchronization tool designed to help developers and teams securely manage, share, and sync environment variables across projects and environments. Built with security as a top priority, ENVM provides encryption at rest, role-based access control, and seamless collaboration features.

### Why ENVM?

- **🔒 Secure by Default**: AES-256 encryption, audit logging, and RBAC
- **🤝 Team Collaboration**: Share sensitive environment variables safely without exposure
- **🔄 Multi-Environment Support**: Manage dev, staging, and production environments effortlessly
- **📦 Universal Format**: `.envm` file format with syntax highlighting and IDE support
- **🚀 Developer Friendly**: CLI, Web UI, and API access
- **🌐 Cross-Platform**: Linux, macOS, Windows support

## ✨ Features

### Core Features

- ✅ **Secure Variable Management**: Create, update, delete, and organize environment variables
- ✅ **Encryption**: AES-256 encryption for stored variables
- ✅ **Access Control**: Role-based permissions (read/write/admin)
- ✅ **Audit Logging**: Track who accessed or modified variables
- ✅ **Project Organization**: Group variables by projects and environments
- ✅ **Team Collaboration**: Share variables securely within organizations
- ✅ **API-First**: RESTful API for programmatic access
- ✅ **CLI Tool**: Command-line interface for CI/CD integration
- ✅ **Web Interface**: Modern, responsive web UI

### Security Features

- 🔐 Encryption at rest and in transit (TLS 1.3)
- 🔑 JWT-based authentication with refresh token rotation
- 👥 Role-based access control (RBAC)
- 📊 Comprehensive audit logging
- 🔄 Secret rotation reminders
- 🛡️ Rate limiting and brute force protection
- 🔍 Input validation and sanitization

### Upcoming Features

- [ ] HashiCorp Vault integration
- [ ] AWS Secrets Manager integration
- [ ] Azure Key Vault integration
- [ ] Automatic secret rotation
- [ ] WebSocket real-time sync
- [ ] NPM package wrapper
- [ ] VS Code extension
- [ ] Mobile app

## � Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/envm-org/envm.git
   cd envm
   ```

2. **Start the application**:
   ```bash
   docker compose up --build
   ```

3. **Access ENVM**:
   - API: `http://localhost:5000`
   - Web UI: [https://envm.vercel.app/](https://envm.vercel.app/)

4. **Verify the installation**:
   ```bash
   curl -i http://localhost:5000/health
   ```

### Configuration

Create a `.env` file in the root directory with the following variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URI` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | Secret key for JWT signing |
| `API_PORT` | No | API server port (default: 5000) |
| `ENCRYPTION_KEY` | Yes | Key for encrypting environment variables |

Example `.env`:
```bash
DATABASE_URI=postgres://postgres:password@localhost:5432/envm?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
API_PORT=5000
ENCRYPTION_KEY=your-32-byte-encryption-key
```

## 📚 Documentation

## 🤝 Contributing

We welcome contributions from the community! 

- Read our [Contributing Guidelines](CONTRIBUTING.md) to get started
- Check our [Code of Conduct](CODE_CODUCT.md)
- Browse [open issues](https://github.com/envm-org/envm/issues) or create a new one
- Submit pull requests with improvements

For developers, see [CLI_REFERENCE.md](CLI_REFERENCE.md) for development setup and commands.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔒 Security

Security is paramount for ENVM. If you discover a security vulnerability, please follow our [Security Policy](SECURITY.md) for responsible disclosure.

**Do not report security vulnerabilities through public GitHub issues.**

## 📞 Support & Community

- **Documentation**: [envm-docs.vercel.app](https://envm-docs.vercel.app/)
- **Issues**: [GitHub Issues](https://github.com/envm-org/envm/issues)
- **Discussions**: [GitHub Discussions](https://github.com/envm-org/envm/discussions)
- **Email**: support@envm.dev

## � Project Status

ENVM is currently in active development. We're working towards our MVP release:

- [x] Core API implementation
- [x] Authentication & authorization
- [x] PostgreSQL integration
- [x] Docker setup
- [ ] CLI tool
- [ ] Web UI
- [ ] Documentation site
- [ ] Public beta release

See our [project roadmap](https://github.com/envm-org/envm/projects) for upcoming features.

## 🙏 Acknowledgments

Thanks to all our [contributors](https://github.com/envm-org/envm/graphs/contributors)!

