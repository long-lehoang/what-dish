# BE Code Review

Review the current Go backend code changes for quality, best practices, and CI compliance.

If `$ARGUMENTS` is provided, focus the review on those files or directories. Otherwise, review all uncommitted changes in the `BE/` directory (`git diff -- BE/` + `git diff --cached -- BE/`).

## Step 1: Identify changes

Run `git diff --name-only -- BE/` and `git diff --cached --name-only -- BE/` to list all modified Go files. If no git changes exist, review the files specified in `$ARGUMENTS`.

## Step 2: Run CI tools (all must pass)

Run these commands from the `BE/` directory and report results:

```bash
make fmt-check       # gofmt formatting check
make vet             # go vet static analysis
make lint            # golangci-lint (strict config)
make test            # Unit tests
make test-integration # Integration tests (testcontainers + real Postgres)
make build           # Binary compiles
```

For full E2E validation (optional, slower):
```bash
make docker-up       # Start postgres + app containers
make test-e2e        # Run E2E tests against running services
make docker-down     # Cleanup
```

If any check fails, report the exact error and suggest a fix.

## Step 3: Code quality review

For each changed file, check:

### Architecture & Clean Architecture
- **Dependency rule**: dependencies point INWARD only (handler → service → port ← repository)
- **Cross-context isolation**: bounded contexts (recipe, suggestion, user, nutrition, engagement) NEVER import each other directly
- **Port/Adapter pattern**: interfaces defined in `port.go` where consumed, implemented by repositories
- **No business logic in handlers** — handlers only parse requests, call services, return responses
- **No direct DB access in services** — always go through repository interfaces
- **Thin DTOs**: request/response structs in `dto.go`, mapped via explicit functions

### Go Best Practices
- **Error handling**: every error checked, wrapped with `fmt.Errorf("context: %w", err)`
- **No naked returns** in functions with named return values
- **Context propagation**: `context.Context` as first parameter for all I/O functions
- **Receiver naming**: 1-2 chars (`s` for service, `h` for handler, `r` for repo)
- **Interface segregation**: small, focused interfaces (not god interfaces)
- **No `panic`** in business logic — only for unrecoverable startup errors
- **Structured logging**: `slog` with key-value pairs, no `fmt.Println` or `log.Println`

### Error Wrapping & Handling
- All repository errors wrapped with context: `fmt.Errorf("repo.Method: %w", err)`
- `rows.Err()` checked after every `rows.Next()` loop
- Transaction rollback errors handled (at minimum logged)
- Custom error types (`ErrNotFound`, `ErrValidation`, etc.) used for business errors
- HTTP status mapping correct in handlers (404 for not found, 400 for validation, etc.)

### Database & Repository
- SQL injection prevention (parameterized queries, never string concatenation)
- Transactions used for multi-table writes (ingredients + steps + tags)
- `rows.Close()` via `defer` immediately after query
- Connection pool not leaked (no unclosed rows or connections)
- Soft delete respected (`WHERE deleted_at IS NULL`)
- COALESCE for nullable columns when scanning into non-pointer types

### Security
- No secrets or credentials in code
- No user input directly interpolated into SQL queries
- Auth middleware on protected routes
- CORS properly configured (only `FRONTEND_URL` allowed)
- Input validation on all handler endpoints

### Testing
- New logic has corresponding tests
- Tests cover edge cases and error paths
- Test names describe behavior: `TestServiceName_MethodName_Scenario`
- Mocks implement port interfaces correctly
- No test pollution (each test independent, no shared mutable state)
- Repository tests use testcontainers with real Postgres
- E2E tests call real HTTP endpoints via docker-compose

### Performance
- No N+1 queries (batch fetches where possible)
- Pagination on list endpoints (default limit=20, max=100)
- Indexes exist for filtered/sorted columns
- No blocking operations in request path without timeout

## Step 4: Report

Present findings as a structured report:

```
## CI Results
- fmt-check: PASS/FAIL
- vet: PASS/FAIL
- lint: PASS/FAIL
- test: PASS/FAIL (X passed, Y failed)
- test-integration: PASS/FAIL (X passed, Y failed)
- build: PASS/FAIL

## Issues Found
### Critical (must fix)
- [file:line] Description and suggested fix

### Warning (should fix)
- [file:line] Description and suggested fix

### Suggestion (nice to have)
- [file:line] Description

## Summary
X files reviewed, Y issues found (Z critical, W warnings)
```

If all CI passes and no issues found, report a clean bill of health.
