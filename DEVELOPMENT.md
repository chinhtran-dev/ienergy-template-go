# Development Guideline

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Project Structure](#project-structure)
3. [Getting Started](#getting-started)
4. [Development Workflow](#development-workflow)
5. [Database Migrations](#database-migrations)
6. [API Documentation](#api-documentation)
7. [Testing](#testing)
8. [Code Style](#code-style)
9. [Git Workflow](#git-workflow)
10. [Deployment](#deployment)

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Make

### Required Tools

Install all required development tools using:
```bash
make install-tools
```

This will install:
- Swag CLI (for API documentation)
- Migrate CLI (for database migrations)
- GolangCI-Lint (for code linting)
- Mockery (for generating mocks)

## Project Structure

```
.
├── cmd/                    # Application entry points
│   ├── app/               # Main application
│   │   └── main.go        # Application entry point
│   └── migrate/           # Migration commands
├── config/                # Configuration files
├── docs/                  # Documentation
│   └── swagger/          # Swagger API documentation
├── internal/              # Private application code
│   ├── app/              # Application core
│   ├── http/             # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── model/            # Data models
│   ├── repository/       # Data access layer
│   └── service/          # Business logic
├── migrations/            # Database migrations
├── pkg/                  # Public packages
│   ├── constant/         # Constants
│   ├── database/         # Database utilities
│   ├── errormap/         # Error mapping
│   ├── ginbuilder/       # Gin framework utilities
│   ├── graceful/         # Graceful shutdown
│   ├── logger/           # Logging utilities
│   ├── swagger/          # Swagger utilities
│   ├── tracking/         # Request tracking
│   ├── util/             # Utility functions
│   └── wrapper/          # Response wrappers
├── test/                 # Test files
├── .env                  # Environment variables
├── .env.example          # Example environment variables
├── .golangci.yml         # GolangCI-Lint configuration
├── Makefile              # Build and development commands
└── go.mod                # Go module definition
```

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd ienergy-template-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Install required tools:
   ```bash
   make install-tools
   ```

5. Run the application:
   ```bash
   make run
   ```

## Development Workflow

1. Create a new feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the code style guidelines

3. Run tests:
   ```bash
   make test
   ```

4. Update API documentation if needed:
   ```bash
   make swagger-build
   ```

5. Commit your changes:
   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

6. Push your changes:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a pull request

## Database Migrations

1. Create a new migration:
   ```bash
   make migrate-create name=your_migration_name
   ```

2. Edit the generated migration files in `migrations/`:
   - `{version}_your_migration_name.up.sql`
   - `{version}_your_migration_name.down.sql`

3. Apply migrations:
   ```bash
   make migrate-up
   ```

4. Rollback migrations:
   ```bash
   make migrate-down
   ```

## API Documentation

1. Initialize Swagger:
   ```bash
   make swagger-init
   ```

2. Update Swagger documentation:
   ```bash
   make swagger-build
   ```

3. View API documentation:
   - Run the application
   - Visit `http://localhost:8080/swagger/index.html`

## Testing

1. Run all tests:
   ```bash
   make test
   ```

2. Run specific test types:
   ```bash
   make test-unit      # Run unit tests
   make test-http      # Run HTTP tests
   make test-all       # Run all test types
   ```

3. Run linting:
   ```bash
   make lint
   ```

## Code Style

1. Use `gofmt` for code formatting:
   ```bash
   go fmt ./...
   ```

2. Follow these naming conventions:
   - Packages: lowercase, single word
   - Interfaces: `I` prefix (e.g., `IUserService`)
   - Structs: PascalCase
   - Functions: PascalCase for public, camelCase for private
   - Variables: camelCase
   - Constants: UPPER_SNAKE_CASE

3. Comment your code:
   - Use `//` for single-line comments
   - Use `/* */` for multi-line comments
   - Document all public functions and types

## Git Workflow

1. Branch naming:
   - `feature/`: New features
   - `bugfix/`: Bug fixes
   - `hotfix/`: Urgent fixes
   - `release/`: Release preparation

2. Commit messages:
   - Format: `<type>: <description>`
   - Types: feat, fix, docs, style, refactor, test, chore
   - Example: `feat: add user authentication`

3. Pull requests:
   - Create from feature branch to main
   - Include description of changes
   - Link related issues
   - Request review from team members

## Deployment

1. Build the application:
   ```bash
   make build
   ```

2. The binary will be created in `bin/app`

3. Deploy the binary to your server

4. Run database migrations:
   ```bash
   make migrate-up
   ```

5. Start the application:
   ```bash
   ./bin/app api
   ```

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Swagger Documentation](https://swaggo.github.io/swaggo.io/)
- [Migrate Documentation](https://github.com/golang-migrate/migrate)
- [GolangCI-Lint Documentation](https://golangci-lint.run/) 