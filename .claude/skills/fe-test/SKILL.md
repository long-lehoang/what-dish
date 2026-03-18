# FE Test

Run Next.js frontend tests and report results.

If `$ARGUMENTS` is provided, run only that test tier (`unit`, `e2e`, or a specific file/pattern). Otherwise, run all tiers.

## Test tiers

### Unit & Component tests (Vitest + React Testing Library)
```bash
cd FE && pnpm test
```
Tests live in `src/**/*.test.{ts,tsx}`, run in jsdom environment.

### Coverage report
```bash
cd FE && pnpm test:coverage
```

### E2E tests (Playwright)
```bash
cd FE && pnpm test:e2e
```
Tests live in `src/__tests__/*.spec.ts`. Playwright auto-starts dev server on `:3000`.
Runs on Mobile Safari (iPhone SE) and Desktop Chrome.

## Step 1: Run tests

Based on `$ARGUMENTS`:
- `unit` → `pnpm test`
- `e2e` → `pnpm test:e2e`
- Specific file/pattern → `pnpm vitest run $ARGUMENTS`
- No arguments → run unit then e2e

## Step 2: Coverage report (if unit tests pass)

```bash
cd FE && pnpm test:coverage
```

Report per-feature coverage.

## Step 3: Report results

```
## Test Results

### Unit & Component Tests
- Total: X | Passed: Y | Failed: Z | Skipped: W

### E2E Tests (if run)
- Total: X | Passed: Y | Failed: Z
- Browsers: Mobile Safari, Desktop Chrome

## Failed Tests (if any)
- [file] TestName — error message and relevant output

## Coverage
- features/random: XX%
- features/recipe: XX%
- features/vote: XX%
- features/explore: XX%
- features/dish: XX%
- shared: XX%
- Overall: XX%

## Recommendations
- Uncovered areas that should have tests
- Flaky test warnings
```

Target coverage: 80%+ for hooks/utils/stores, key user flows covered by E2E.
