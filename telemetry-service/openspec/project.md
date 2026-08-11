# Project Context

## Purpose
A telemetry service built with the Kratos microservices framework for collecting, processing, and managing telemetry data. The service provides both gRPC and HTTP/REST APIs for flexible integration.

## Tech Stack
- Go 1.25
- Kratos v2.9.2 (Go microservices framework)
- Protocol Buffers (protobuf) for API definitions
- gRPC for high-performance RPC
- HTTP/REST with OpenAPI specification
- PostgreSQL (primary database with connection pooling)
- Redis (caching layer)
- Google Wire (dependency injection)
- OpenTelemetry (distributed tracing)
- Docker (containerization)

## Project Conventions

### Code Style
- Follow standard Go conventions (gofmt, go vet)
- Use Protocol Buffers for all API definitions
- Repository pattern for data access
- Use case pattern for business logic
- Dependency injection via Google Wire

### Architecture Patterns
This project follows Clean Architecture with clear separation of concerns:

- **api/** - Protocol Buffer definitions and generated code (HTTP, gRPC, OpenAPI)
- **internal/service/** - Service layer that handles transport protocol (gRPC/HTTP handlers)
- **internal/biz/** - Business logic layer containing use cases and domain models
- **internal/data/** - Data access layer implementing repository interfaces
- **internal/server/** - Server configuration and initialization (HTTP & gRPC servers)
- **internal/conf/** - Configuration structures generated from protobuf
- **cmd/** - Application entry points with Wire dependency injection

Dependencies flow inward: service → biz → data. The biz layer defines repository interfaces, and the data layer implements them.

### Testing Strategy
- Unit tests for business logic in the biz layer
- Integration tests for repository implementations
- API tests for service endpoints
- Use Go's standard testing package

### Git Workflow
- Main branch: `main`
- Feature development on feature branches
- Follow conventional commit messages

## Domain Context
The service is designed to handle telemetry data collection and processing. Currently includes a sample "greeter" service that demonstrates the architecture pattern.

### Key Concepts
- **Use Cases**: Business logic encapsulated in the biz layer
- **Repositories**: Abstract data access defined as interfaces in biz, implemented in data layer
- **Services**: Protocol-aware handlers that delegate to use cases
- **Transport Agnostic**: Same business logic serves both gRPC and HTTP clients

## Important Constraints
- Database connections are pooled (max 100 open, 10 idle connections)
- Server timeouts set to 3 seconds for both HTTP and gRPC
- Connection max lifetime: 1 hour
- Connection max idle time: 10 minutes
- Must maintain backward compatibility with protobuf API changes

## External Dependencies
- **PostgreSQL Database**: Primary data store at Supabase
  - Host: sbp-81n5qyn2id6j45cr.supabase.opentrust.net:5432
  - Connection pooling enabled

- **Redis Cache**: Local Redis instance at 127.0.0.1:6379
  - Read timeout: 0.2s
  - Write timeout: 0.2s

- **Third-party proto definitions**: Google APIs, validation rules (in third_party/)

## Build & Development
- `make init` - Install required tools (protoc plugins, wire, kratos CLI)
- `make api` - Generate API code from proto files
- `make config` - Generate internal config code
- `make generate` - Run go generate and wire
- `make build` - Build the service binary
- `make all` - Generate all code (api + config + generate)
- Wire is used for compile-time dependency injection (run `wire` in cmd/telemetry-service/)
