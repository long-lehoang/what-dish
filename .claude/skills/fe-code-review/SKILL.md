# FE Code Review

Review the current Next.js frontend code changes for quality, best practices, and CI compliance.

If `$ARGUMENTS` is provided, focus the review on those files or directories. Otherwise, review all uncommitted changes in the `FE/` directory (`git diff -- FE/` + `git diff --cached -- FE/`).

## Step 1: Identify changes

Run `git diff --name-only -- FE/` and `git diff --cached --name-only -- FE/` to list all modified files. If no git changes exist, review the files specified in `$ARGUMENTS`.

## Step 2: Run CI tools (all must pass)

Run these commands from the `FE/` directory and report results:

```bash
pnpm type-check       # TypeScript strict compilation
pnpm lint             # ESLint (no warnings allowed)
pnpm format:check     # Prettier formatting
pnpm test             # Vitest unit + component tests
pnpm build            # Production build succeeds
```

If any check fails, report the exact error and suggest a fix.

## Step 3: Code quality review

For each changed file, check:

### Architecture & Design
- SOLID principles followed
- DRY — no copy-pasted logic that should be abstracted
- Separation of concerns — components under 150 lines
- Feature-Sliced Design boundaries respected (features don't import from other features' internals)
- Barrel exports used correctly (import from `index.ts`, not internal paths)

### TypeScript
- No `any` types — use `unknown` and narrow
- Strict mode compliance
- Proper null handling (optional chaining, nullish coalescing)
- Interface for shapes, Type for unions

### React & Next.js
- Proper `'use client'` directives where needed
- No unnecessary re-renders (stable refs, memoization only when measured)
- Server vs client component split is correct
- `next/image` used for images, `next/link` for navigation

### Accessibility
- Semantic HTML elements
- `aria-label` on interactive elements without visible text
- Keyboard navigation support
- Focus management in modals/overlays

### Security
- No secrets or credentials in code
- No `dangerouslySetInnerHTML` without sanitization
- No user input directly interpolated into URLs or queries

### Performance
- Only `transform` and `opacity` for animations (GPU-accelerated)
- Images use `next/image` with proper `sizes` attribute
- No blocking operations in render path

### Testing
- New logic has corresponding tests
- Tests cover edge cases and error paths
- Test names describe behavior, not implementation

## Step 4: Report

Present findings as a structured report:

```
## CI Results
- type-check: PASS/FAIL
- lint: PASS/FAIL
- format:check: PASS/FAIL
- test: PASS/FAIL (X passed, Y failed)
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
