# Contributing to ENVM

Thank you for your interest in contributing to ENVM! We welcome contributions from the community and are grateful for your support in making environment variable management safer and easier for developers.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.25.5 or higher** - [Download Go](https://golang.org/dl/)
- **Docker and Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **PostgreSQL** (for local development) - [Install PostgreSQL](https://www.postgresql.org/download/)
- **Make** - Usually pre-installed on Unix systems
- **Git** - [Install Git](https://git-scm.com/downloads)

### Optional Tools

- **golangci-lint** - For code linting
- **sqlc** - For database code generation
- **Air** - For hot reloading during development

## Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/envm.git
   cd envm
   ```

3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/envm-org/envm.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Copy the example environment file**:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your local configuration.

6. **Start the development environment**:
   ```bash
   docker compose up --build
   ```
   
   The application will be available at `http://localhost:5000`

7. **Run the application** (alternative to Docker):
   ```bash
   make run
   ```


## How to Contribute

### Reporting Bugs

1. **Check existing issues** to avoid duplicates
2. **Create a new issue** with:
   - Clear, descriptive title
   - Detailed description of the bug
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)
   - Relevant logs or screenshots

### Suggesting Features

1. **Check existing feature requests**
2. **Create a new issue** with:
   - Clear title starting with "[Feature Request]"
   - Detailed description of the feature
   - Use cases and benefits
   - Possible implementation approaches (optional)

### Contributing Code

1. **Find an issue** to work on or create one
2. **Comment on the issue** to let others know you're working on it
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes** following our coding standards
5. **Test your changes** thoroughly
6. **Commit your changes** following our commit message guidelines
7. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
8. **Create a Pull Request** from your fork to `envm-org/envm:main`

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting (automatically applied by most editors)
- Run `golangci-lint` before committing:
  ```bash
  make lint
  ```
- Keep functions small and focused (ideally < 50 lines)
- Add comments for exported functions and complex logic
- Use meaningful variable and function names
- Avoid global variables when possible
- Handle errors explicitly - don't ignore them

### Code Organization

- Place business logic in `internal/` packages
- Keep handlers thin - delegate logic to services
- Use dependency injection for better testability
- Follow the repository pattern for data access
- Use interfaces for external dependencies

### Error Handling

```go
// Good
func GetUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}

// Avoid
func GetUser(id string) *User {
    user, _ := db.FindUser(id)  // Don't ignore errors
    return user
}
```

### Security Best Practices

- Never commit secrets, API keys, or credentials
- Always sanitize user input
- Use prepared statements for database queries
- Validate and sanitize file paths
- Use context for timeouts and cancellation
- Implement proper authentication and authorization
- Follow principle of least privilege

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `build`: Build system or dependency changes
- `ci`: CI/CD configuration changes
- `chore`: Other changes (maintenance, tooling)

### Examples

```
feat(auth): add JWT token refresh endpoint

fix(env): resolve variable interpolation bug in nested groups

docs(readme): update installation instructions for Windows

test(users): add unit tests for user service
```

## Pull Request Process

1. **Update documentation** if you're changing functionality
2. **Add tests** for new features or bug fixes
3. **Ensure all tests pass**:
   ```bash
   make test
   ```
4. **Run linters**:
   ```bash
   make lint
   ```
5. **Update CHANGELOG.md** if applicable
6. **Fill out the PR template** completely
7. **Request review** from maintainers
8. **Address review comments** promptly
9. **Squash commits** if requested before merging

### PR Title Format

Use the same format as commit messages:
```
feat(auth): implement OAuth2 provider support
```

### PR Checklist

- [ ] Tests added/updated and passing
- [ ] Documentation updated
- [ ] Code follows style guidelines
- [ ] No breaking changes (or clearly documented)
- [ ] CHANGELOG.md updated
- [ ] Self-review completed
- [ ] Comments added for complex logic

## Testing Guidelines

### Writing Tests

- Write unit tests for all business logic
- Use table-driven tests where appropriate
- Mock external dependencies
- Test edge cases and error conditions
- Aim for >80% code coverage for critical paths

### Test Structure

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserInput
        want    *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserInput{
                Email: "user@example.com",
                Name:  "Test User",
            },
            want:    &User{Email: "user@example.com"},
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewUserService(mockDB)
            got, err := service.CreateUser(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("CreateUser() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/auth/...

# Run with race detector
go test -race ./...
```

## Documentation

### Code Documentation

- Document all exported functions, types, and packages
- Use godoc format for comments
- Include examples for complex functionality
- Keep documentation up-to-date with code changes

### User Documentation

- Update README.md for user-facing changes
- Add tutorials for new features
- Update CLI reference documentation
- Maintain changelog

## Community

### Getting Help

- **GitHub Discussions**: For questions and discussions
- **GitHub Issues**: For bug reports and feature requests
- **Documentation**: Check our [docs](https://docs.envm.dev) first

### Maintainers

- Check [CODEOWNERS](.github/CODEOWNERS) for area-specific maintainers
- Be patient - maintainers are volunteers

### Recognition

Contributors are recognized in:
- Release notes
- CONTRIBUTORS.md file
- GitHub contributors page

## License

By contributing to ENVM, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to ENVM! Your efforts help make environment variable management better for developers everywhere. 🚀
