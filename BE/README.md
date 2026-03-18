# 🍜 WhatDish — Vietnamese Food Suggestion System

> Smart food suggestion app: random dish picker, calorie-based recommendations, group meal planning, with detailed Vietnamese cooking recipes.

## Overview

VietFood solves the eternal question **"What should I eat today?"** by combining:

- **Smart Random** — not just random dish names, but full recipes with nutrition info
- **Calorie-Based Suggestions** — tailored to health goals (weight loss, muscle gain, maintenance)
- **Group Meal Planning** — curated meal combos for families, friend groups, couples
- **Rich Recipe Library** — step-by-step instructions, ingredient measurements, nutrition data

## System Architecture

The project follows **Microservices Architecture** with clear Bounded Contexts, ready to scale:

```
┌─────────────┐
│  API Gateway │
└──────┬──────┘
       │
┌──────┴──────────────────────────────────┐
│            Message Bus (Kafka)           │
└──┬──────┬──────┬──────┬──────┬─────────┘
   │      │      │      │      │
┌──┴──┐┌──┴──┐┌──┴──┐┌──┴──┐┌──┴──┐
│Recip││Sugge││User ││Nutri││Engag│
│  e  ││stion││     ││tion ││ment │
└──┬──┘└──┬──┘└──┬──┘└──┬──┘└──┬──┘
   │      │      │      │      │
┌──┴──┐┌──┴──┐┌──┴──┐┌──┴──┐┌──┴──┐
│ DB  ││ DB  ││ DB  ││ DB  ││ DB  │
└─────┘└─────┘└─────┘└─────┘└─────┘
```

## Documentation

| File | Description |
|------|-------------|
| [`docs/USE_CASES.md`](docs/USE_CASES.md) | Detailed flow for each use case |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Target microservices architecture (future state) |
| [`docs/INFRASTRUCTURE.md`](docs/INFRASTRUCTURE.md) | **Actual deployment** — free tier stack (Vercel + Render + Supabase) |
| [`docs/DATABASE.md`](docs/DATABASE.md) | Database schema |
| [`docs/API_DESIGN.md`](docs/API_DESIGN.md) | RESTful API endpoints |

## Tech Stack

### Current (Free Tier — $0/month)

| Layer | Technology | Free Tier |
|-------|------------|-----------|
| Frontend | Next.js on **Vercel** | 100 GB bandwidth |
| Backend | Go (Gin) on **Render** | 750 hrs/month |
| Database | **Supabase** PostgreSQL | 500 MB storage |
| Auth | **Supabase Auth** | 50K MAU |
| File Storage | **Supabase Storage** | 1 GB |
| Content CMS | **Notion** | Free — admins manage recipes here |
| Search | PostgreSQL full-text (`tsvector`) | included |
| Cache | In-memory (`go-cache`) | included |

### Target (Microservices — when scaling)

| Layer | Technology |
|-------|------------|
| API Gateway | Kong / Nginx |
| Backend Services | Go microservices |
| Database | PostgreSQL per service |
| Cache | Redis |
| Message Queue | Apache Kafka / RabbitMQ |
| Search | Elasticsearch |
| Orchestration | Kubernetes |

> Architecture is designed as **Modular Monolith organized by feature** — each bounded context (recipe, suggestion, user, nutrition, engagement) lives in its own self-contained folder. Extracting to microservices = copy the folder. See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for the target state and [INFRASTRUCTURE.md](docs/INFRASTRUCTURE.md) for the current deployment.

## Development Roadmap

| Phase | Features | Priority |
|-------|----------|----------|
| **1 — MVP** | Random dish, view recipe, search, filter | High |
| **2 — Core** | Calorie-based suggestion, group suggestion | High |
| **3 — User** | Auth, favorites, nutrition profile, history | Medium |
| **4 — Admin** | Manage recipes, users, categories | Medium |
| **5 — Scale** | Elasticsearch, caching, recommendation engine | Low |

## Getting Started

### Prerequisites

- Node.js 20+
- A free [Supabase](https://supabase.com) account
- A free [Render](https://render.com) account
- A free [Vercel](https://vercel.com) account

### Local Development

```bash
# Clone repo
git clone https://github.com/your-username/whatdish.git
cd whatdish

# Install dependencies
npm install

# Copy env files
cp .env.example .env.local

# Set your Supabase credentials in .env.local

# Run database migrations
npx supabase db push

# Seed initial data
npx supabase db seed

# Start development server
npm run dev
```

### Deploy

```bash
# Frontend → Vercel (connect GitHub repo, auto-deploy)
# Backend  → Render (connect GitHub repo, set env vars)
# Database → Supabase (already hosted)
```

## License

MIT