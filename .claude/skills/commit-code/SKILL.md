---
name: commit-code
description: Stage and commit code changes with a well-formatted conventional commit message. Detects FE/BE changes and runs appropriate pre-commit checks.
disable-model-invocation: true
allowed-tools: Bash, Read, Glob, Grep
---

# Commit Code

Commit the current changes with a proper conventional commit message.

If `$ARGUMENTS` is provided, use it as guidance for the commit message. Otherwise, analyze the diff to generate an appropriate message.

## Step 1: Understand the changes

Run in parallel:
```bash
git status                           # See all changed/untracked files
git diff                             # Unstaged changes
git diff --cached                    # Staged changes
git log --oneline -5                 # Recent commits for style reference
```

## Step 2: Determine affected areas

Check which parts of the monorepo are affected:
```bash
# Check for FE changes
git diff --name-only HEAD | grep '^FE/' && echo "FE_CHANGED=true" || echo "FE_CHANGED=false"

# Check for BE changes
git diff --name-only HEAD | grep '^BE/' && echo "BE_CHANGED=true" || echo "BE_CHANGED=false"
```

## Step 3: Pre-commit quality checks

### If FE changed — run from the `FE/` directory:
```bash
pnpm type-check
pnpm lint
pnpm format:check
pnpm test
```

If `format:check` fails, auto-fix with `pnpm format` and include the formatted files in the commit.

### If BE changed — run from the `BE/` directory:
```bash
make fmt-check
make vet
make test
```

If `fmt-check` fails, auto-fix with `gofmt -w` on the affected files and include them in the commit.

If any other check fails, **stop and report the error**. Do not commit broken code.

## Step 4: Stage files

- Stage specific files relevant to the change (prefer explicit `git add <file>` over `git add .`)
- Never stage files that likely contain secrets (`.env`, `.env.local`, credentials, keys)
- If unsure about a file, ask before staging

## Step 5: Create the commit

### Message format

```
<type>(<scope>): <subject>

<body>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

### Types
| Type | When |
|------|------|
| `feat` | New feature or functionality |
| `fix` | Bug fix |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or updating tests |
| `style` | Formatting, missing semicolons, etc (no logic change) |
| `docs` | Documentation changes |
| `chore` | Build process, deps, tooling changes |
| `perf` | Performance improvement |

### Scope
Use the area of the codebase affected. Examples:
- `fe`, `be` — for changes scoped to one side
- `fe/recipe`, `be/nutrition` — feature-level scope
- `ci`, `infra` — for workflow/config changes
- Omit scope for changes spanning the whole repo

### Rules
- Subject line: imperative mood, lowercase, no period, under 72 chars
- Body: explain **what** and **why**, not how
- Reference issue numbers if mentioned: `Fixes #123`
- All UI text in Vietnamese, all code/comments/commits in English

### Example
```
feat(fe/vote): add tournament bracket elimination mode

Implement 1v1 bracket voting where dishes are paired and winners
advance to the next round. Uses Framer Motion for slide-in/out
card transitions.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

## Step 6: Verify

After committing, run `git log -1` and `git status` to confirm the commit was created successfully.

## Important
- NEVER amend existing commits unless explicitly asked
- NEVER force push
- NEVER skip hooks (no `--no-verify`)
- NEVER push to remote unless explicitly asked
- Always use HEREDOC for multi-line commit messages to preserve formatting
