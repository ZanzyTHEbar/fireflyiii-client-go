---
applyTo: '**'
---


# Firefly Client Go - Cursor Rules

## Project Organization
- The project is organized as a Go package with the main client code in the root directory
- Generated OpenAPI code is in firefly.gen.go
- Custom client implementations extend the generated code with more idiomatic Go interfaces

## Code Style Guidelines
- Follow standard Go style conventions (gofmt)
- Use descriptive variable names
- Include comments for all exported functions, types, and methods
- Organize imports alphabetically with standard library first

## Error Handling
- All functions that can fail should return an error
- Custom error types are used to provide more context
- Errors from the API should preserve the original error message
- Use named return values where appropriate for better documentation

## API Pattern
- The main client object provides access to service-specific methods
- Services are organized by resource type (accounts, transactions, etc.)
- All API methods accept a context.Context as the first parameter
- Optional parameters use the functional options pattern

## Feature Tracking
- Feature implementation status is tracked in memory-bank/progress.md
- API coverage is a priority to ensure all Firefly III endpoints are accessible
- Documentation and examples should be added for each new service implementation

## Development Workflow
- Work from the OpenAPI spec to ensure accuracy
- Add tests for all new functionality
- Update documentation when adding new features
- Follow semantic versioning for releases

## Critical Paths
- Client initialization and authentication are foundational
- Error handling is critical for debugging API interactions
- Pagination handling affects all list operations 