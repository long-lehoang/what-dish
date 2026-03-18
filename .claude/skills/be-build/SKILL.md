# BE Build

Build the Go backend and verify it compiles cleanly.

## Step 1: Compile check

From the `BE/` directory:

```bash
go build ./...
```

If compilation fails, report the exact errors with file and line numbers.

## Step 2: Vet and format check

```bash
make fmt-check       # List unformatted files
make vet             # Static analysis
```

## Step 3: Binary build

```bash
make build           # CGO_ENABLED=0 go build -o bin/server ./cmd/server
```

## Step 4: Docker build (if requested via $ARGUMENTS containing "docker")

```bash
make docker-up       # docker-compose up -d --build (postgres + app)
```

Verify the app is healthy:
```bash
curl -s http://localhost:8080/health | jq .
```

Use `make docker-rebuild` to rebuild only the app container after code changes.
Use `make docker-down` to stop all containers.

## Step 5: Report

```
## Build Results
- Compilation: PASS/FAIL
- go vet: PASS/FAIL
- gofmt: PASS/FAIL (X files need formatting)
- Binary: PASS/FAIL
- Docker (if run): PASS/FAIL

## Errors (if any)
- [file:line] Error description
```

If all pass, confirm the build is clean with a one-line summary.
