# BE Test

Run Go backend tests and report results.

If `$ARGUMENTS` is provided, run only that test tier (`unit`, `integration`, `e2e`, or a specific package path). Otherwise, run all tiers.

## Test tiers

### Unit tests (fast, no external deps)
```bash
cd BE && make test
```

### Integration tests (testcontainers — spins up real Postgres)
```bash
cd BE && make test-integration
```

### E2E tests (docker-compose — real HTTP calls to running app)
```bash
# Services must be running first
cd BE && make docker-up
# Wait for healthy, then:
cd BE && make test-e2e
```

E2E tests start a fake Supabase auth server on `:9999`, run migrations + seeds, then call the real app API at `localhost:8080`.

## Step 1: Run tests

Based on `$ARGUMENTS`:
- `unit` → `make test`
- `integration` → `make test-integration`
- `e2e` → `make docker-up && make test-e2e`
- Package path (e.g. `./internal/recipe/...`) → `go test -v -count=1 $ARGUMENTS`
- No arguments → run all three tiers in order

Use `-count=1` to disable test caching.

## Step 2: Coverage report (if unit + integration pass)

```bash
cd BE && make test-cover
```

Report per-package coverage percentages.

## Step 3: Report results

```
## Test Results

### Unit Tests
- Total: X | Passed: Y | Failed: Z | Skipped: W

### Integration Tests
- Total: X | Passed: Y | Failed: Z

### E2E Tests
- Total: X | Passed: Y | Failed: Z

## Failed Tests (if any)
- [package] TestName — error message and relevant output

## Coverage
- recipe: XX%
- suggestion: XX%
- user: XX%
- nutrition: XX%
- engagement: XX%
- platform/notion: XX%
- shared: XX%
- Overall: XX%

## Recommendations
- Uncovered areas that should have tests
- Flaky test warnings
```

Target coverage: 80%+ for service layer, 100% for TDEE calculator and Notion parser.
