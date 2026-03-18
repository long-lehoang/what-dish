# FE Build

Build the Next.js frontend and verify it compiles cleanly.

## Step 1: Type check

```bash
cd FE && pnpm type-check
```

If TypeScript errors found, report exact file, line, and error.

## Step 2: Lint and format check

```bash
cd FE && pnpm lint
cd FE && pnpm format:check
```

## Step 3: Production build

```bash
cd FE && pnpm build
```

## Step 4: Report

```
## Build Results
- type-check: PASS/FAIL
- lint: PASS/FAIL
- format:check: PASS/FAIL
- build: PASS/FAIL

## Errors (if any)
- [file:line] Error description
```

If all pass, confirm the build is clean with a one-line summary.
