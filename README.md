# Tối Nay Ăn Gì? (WhatDish)

> "Hết phân vân — lật là ăn!" (Stop hesitating — flip and eat!)

A Vietnamese food decision app that helps users pick what to eat through a fun card shuffle experience, calorie-based suggestions, group meal planning, and detailed home cooking recipes.

## Features

- **Random Dish Picker** — Animated card shuffle with category/difficulty/time filters and a 3 re-roll limit
- **Calorie-Based Suggestions** — Dishes matched to your TDEE and health goals (lose/maintain/gain)
- **Group Meal Planning** — Balanced combos (main + soup + sides) for 1–10 people
- **Recipe & Cook Mode** — Step-by-step instructions with serving scaling, built-in timers, and Wake Lock
- **Group Voting** — Real-time rooms with Tournament, Swipe, and Ranking vote strategies
- **Explore** — Browse, search, and filter dishes with infinite scroll

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js 15, TypeScript, Tailwind CSS, Framer Motion, Zustand |
| Backend | Go (Gin), Clean Architecture, PostgreSQL, Supabase Auth |
| Content CMS | Notion (recipes managed by editors, synced to DB) |
| Search | PostgreSQL full-text search (tsvector + unaccent) |
| Deployment | Vercel (FE) + Render (BE) + Supabase (DB/Auth/Storage) |

## Project Structure

```
.
├── FE/          # Next.js frontend (Feature-Sliced Design)
├── BE/          # Go backend (Clean Architecture, modular monolith)
├── .github/     # CI workflows for both FE and BE
└── .claude/     # Claude Code skills (build, test, review, commit)
```

See [FE/README.md](FE/README.md) and [BE/README.md](BE/README.md) for detailed documentation.

## Quick Start

### Prerequisites

- Node.js 18+ & pnpm 8+
- Go 1.23+
- Docker (for local PostgreSQL)

### Frontend

```bash
cd FE
pnpm install
pnpm dev          # http://localhost:3000
```

The frontend works fully offline with mock data when the backend is unavailable.

### Backend

```bash
cd BE
cp .env.example .env
make docker-up    # Start PostgreSQL
make migrate-up   # Run migrations
make seed         # Seed reference data
make dev          # http://localhost:8080 (hot reload)
```

### Environment

See [FE/.env.example](FE/.env.example) and [BE/.env.example](BE/.env.example) for required variables.

## Architecture

```
┌──────────────┐         ┌──────────────┐         ┌─────────┐
│   Next.js    │──REST──▶│   Go API     │──sync──▶│ Notion  │
│   (Vercel)   │◀────────│   (Render)   │◀────────│  (CMS)  │
└──────┬───────┘         └──────┬───────┘         └─────────┘
       │                        │
       │ WebSocket              │ pgx
       │                        ▼
       │                 ┌──────────────┐
       └────────────────▶│  Supabase    │
                         │  (PG + Auth) │
                         └──────────────┘
```

- **Frontend**: Feature-Sliced Design — code organized by domain (random, recipe, vote, explore)
- **Backend**: Clean Architecture with bounded contexts (recipe, suggestion, user, nutrition, engagement)
- **Content**: Recipes managed in Notion, synced to PostgreSQL via the backend's SyncService

## Documentation

| Document | Description |
|----------|-------------|
| [BE/docs/api-design.md](BE/docs/api-design.md) | REST API endpoints |
| [BE/docs/architecture.md](BE/docs/architecture.md) | Backend architecture |
| [BE/docs/database.md](BE/docs/database.md) | Database schema |
| [BE/docs/infrastructure.md](BE/docs/infrastructure.md) | Deployment (free tier stack) |
| [BE/docs/usecases.md](BE/docs/usecases.md) | Detailed use case flows |

## License

Private project.
