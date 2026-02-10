# Agent Guidelines for REST API Go Project

## Build/Test Commands
- **Build**: `go build ./cmd/api`
- **Run**: `go run ./cmd/api`
- **Test all**: `go test ./...`
- **Test package**: `go test ./internal/api/middlewares`
- **Test with coverage**: `go test -cover ./...`
- **Format code**: `go fmt ./...`
- **Vet code**: `go vet ./...`
- **Mod tidy**: `go mod tidy`

## Code Style Guidelines
- **Imports**: Standard library first, then third-party, then local packages (separated by blank lines)
- **Naming**: Use camelCase for functions/variables, PascalCase for exported types/functions
- **Error handling**: Always check errors explicitly; use `log.Fatal()` for critical failures
- **Structs**: Use struct tags for JSON serialization (e.g., `json:"name"`)
- **Comments**: Use `//` for single-line comments; document exported functions
- **Package structure**: `cmd/` for executables, `internal/` for private code
- **HTTP handlers**: Follow pattern `func handlerName(w http.ResponseWriter, r *http.Request)`
- **Middleware**: Return `http.Handler` and use closure pattern with `http.HandlerFunc`
- **Security**: Always set security headers, use TLS configuration, validate input
- **Constants**: Use `const` blocks for related constants