# iEnergy Template Go

A Go-based template project for building scalable and maintainable web applications.

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Technology Stack](#technology-stack)
5. [Getting Started](#getting-started)
6. [Project Structure](#project-structure)
## Overview

`ienergy-template-go` is a boilerplate project designed to help developers quickly set up a Go-based web application. It includes pre-configured tools and best practices for development, testing, and deployment.

## Features

- Modular project structure
- RESTful API with [Gin](https://gin-gonic.com/)
- Database integration with [GORM](https://gorm.io/)
- Swagger API documentation
- Graceful shutdown support
- Built-in linting and testing tools
- Environment-based configuration
- Database migrations with [golang-migrate](https://github.com/golang-migrate/migrate)

## Architecture

The project follows a clean architecture approach, separating concerns into layers:
- **HTTP Layer**: Handles incoming requests and routes them to the appropriate services.
- **Service Layer**: Contains business logic and orchestrates interactions between repositories and other components.
- **Repository Layer**: Manages data access and persistence, abstracting database operations.
- **Model Layer**: Data models and DTOs
- **Middleware**: Provides reusable components for request handling, such as authentication and logging.
- **Infrastructure**: Database connections, etc.
- **Dependency Injection**: Using Uber's FX library

## Technology Stack

- **Programming Language**: Go (Golang)
- **Web Framework**: [Gin](https://gin-gonic.com/)
- **Dependency Injection**: Uber's FX
- **Authentication**: JWT
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io/)
- **API Documentation**: Swagger
- **Linting**: GolangCI-Lint
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)

## Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Make (for running development commands)

### Installation

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

### Environment Variables

The application uses a `.env` file to manage environment-specific configurations. Below are the key variables:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `PORT`: Application port

Refer to `.env.example` for a complete list of variables.

### API Documentation

1. Generate Swagger documentation:
   ```bash
   make swagger-build
   ```

2. Access the Swagger UI at:
   ```
   http://localhost:8080/swagger/index.html
   ```

## Project Structure

```
.
├── cmd/                    # Application entry points
│   ├── app/                # Main application entry point
│       └── main.go         # Starts the application
├── config/                 # Configuration handling
├── docs/                   # Documentation
│   └── swagger/            # Swagger API documentation
├── internal/               # Private application code
│   ├── app/                # Application core logic
│   ├── http/               # HTTP handlers and routes
│   ├── middleware/         # Custom middleware for HTTP requests
│   ├── model/              # Data models and entities
│   ├── repository/         # Data access layer (e.g., database queries)
│   └── service/            # Business logic and services
├── pkg/                    # Public reusable packages
│   ├── constant/           # Application-wide constants
│   ├── database/           # Database connection utilities
│   ├── errors/             # Error custom
│   ├── ginbuilder/         # Utilities for building Gin applications
│   ├── graceful/           # Graceful shutdown utilities
│   ├── logger/             # Logging utilities
│   ├── swagger/            # Swagger integration helpers
│   ├── tracking/           # Request tracking utilities
│   ├── util/               # General utility functions
│   └── wrapper/            # Response wrappers for standardizing API responses
├── test/                   # Test files
│   └── integration/        # Integration tests
├── .env                    # Environment variables file
├── .env.example            # Example environment variables file
├── .golangci.yml           # Configuration for GolangCI-Lint
├── Makefile                # Build and development commands
└── go.mod                  # Go module definition
```