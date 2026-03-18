# Infrastructure — Free Tier Deployment (WhatDish)

> This document describes the **actual deployment** using free-tier services.
> For the full microservices target architecture, see [ARCHITECTURE.md](./ARCHITECTURE.md).

## 1. Stack Overview

| Layer | Service | Free Tier Limits | Purpose |
|-------|---------|-----------------|---------|
| **Frontend** | Vercel | 100 GB bandwidth, serverless functions | Next.js SSR/SSG |
| **Backend** | Render | 750 hrs/month, sleeps after 15 min inactivity | Go API server |
| **Database** | Supabase | 500 MB storage, unlimited API requests | PostgreSQL + Auth + Storage |
| **Auth** | Supabase Auth | 50K MAU | Email/password, Google, Facebook OAuth |
| **File Storage** | Supabase Storage | 1 GB | Recipe images (synced from Notion) |
| **Content CMS** | Notion | Free for personal use | Recipe management by admins |
| **Search** | PostgreSQL FTS | (included in Supabase) | Full-text search via `tsvector` |
| **Cache** | In-memory | (included in Go process) | `go-cache` or `sync.Map` |

**Total monthly cost: $0**

```
┌─────────────────────┐     ┌─────────────────────┐
│   Vercel (free)      │     │   Notion (free)      │
│   Next.js Frontend   │     │   Recipe CMS         │
└──────────┬──────────┘     └──────────┬──────────┘
           │ HTTPS                      │ Notion API
┌──────────▼──────────────────────────▼─┐
│   Render (free)                        │
│   Go Monolith (Gin)                    │
│   ┌───────────────┐  ┌─────────────┐  │
│   │ recipe module  │  │ sync service│  │
│   │ suggest module │  │ (Notion→PG) │  │
│   │ user module    │  └─────────────┘  │
│   │ nutrition mod  │                   │
│   │ engage module  │                   │
│   └───────┬───────┘                   │
│     EventBus (in-process)             │
│     go-cache (in-memory)              │
└──────────┬────────────────────────────┘
           │ Connection pooling (pgx)
┌──────────▼──────────┐
│   Supabase (free)    │
│   ┌──────────────┐   │
│   │ PostgreSQL   │   │
│   │ (500 MB)     │   │
│   ├──────────────┤   │
│   │ Auth         │   │
│   │ (50K MAU)    │   │
│   ├──────────────┤   │
│   │ Storage      │   │
│   │ (1 GB)       │   │
│   └──────────────┘   │
└──────────────────────┘
```

---

## 2. What Changes from the Target Architecture

### Database: 5 DBs → 1 Database, Separate Schemas

Instead of 5 separate PostgreSQL instances, we use **1 Supabase database** with logical separation via **schemas** (or table prefixes). The table definitions from `DATABASE.md` remain identical — we just put them all in one database.

```sql
-- Option A: Use schemas for bounded context separation
CREATE SCHEMA recipe;
CREATE SCHEMA suggestion;
CREATE SCHEMA auth;       -- or use Supabase's built-in auth schema
CREATE SCHEMA nutrition;
CREATE SCHEMA engagement;

-- Tables live in their respective schemas
CREATE TABLE recipe.recipes ( ... );
CREATE TABLE suggestion.suggestion_sessions ( ... );
CREATE TABLE nutrition.recipe_nutrition ( ... );
CREATE TABLE engagement.favorites ( ... );
```

```sql
-- Option B: Use table prefixes (simpler, works with all ORMs)
-- recipe_recipes, recipe_ingredients, recipe_steps...
-- suggestion_sessions, suggestion_configs...
-- nutrition_recipe_nutrition, nutrition_goals...
-- engagement_favorites, engagement_view_history...
```

**Why this still works for microservices later:** When you extract a service, you migrate its schema/tables to a new database. The code only touches its own tables via repository classes, so the migration is straightforward.

### Auth: Custom JWT → Supabase Auth

Supabase provides built-in authentication with 50K MAU on the free tier. This replaces our entire custom JWT implementation.

**What we get for free:**
- Email/password registration + email verification
- OAuth providers (Google, Facebook, GitHub)
- JWT token management (access + refresh)
- Row Level Security (RLS) integration
- Password reset flow
- Session management

**What changes in our code:**
- `POST /auth/register` → `supabase.auth.signUp()`
- `POST /auth/login` → `supabase.auth.signInWithPassword()`
- `POST /auth/refresh` → handled automatically by Supabase client
- JWT verification → `supabase.auth.getUser(token)`
- `users` table → Supabase's `auth.users` + our custom `user_profiles` table

**What we still build ourselves:**
- `user_profiles` table (nutrition data, TDEE)
- `user_allergies` table
- Admin role management (set role in `user_profiles.role`)
- TDEE calculation logic

```go
// Example: Protecting an endpoint (middleware)
func AuthMiddleware(supabaseURL, anonKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractBearerToken(c.GetHeader("Authorization"))
        user, err := supabaseClient.Auth.GetUser(c, token)
        if err != nil {
            c.AbortWithStatusJSON(401, dto.ErrorResponse{Code: "UNAUTHORIZED"})
            return
        }
        c.Set("user_id", user.ID)
        c.Next()
    }
}
```

### Message Bus: Kafka → In-Process Event Bus

In a monolith, cross-module communication uses a Go in-process event bus (channel-based or callback-based) instead of Kafka.

```go
// shared/event/bus.go
type EventBus struct {
    handlers map[string][]HandlerFunc
    mu       sync.RWMutex
}

// recipe module publishes after sync
bus.Publish("recipe.synced", RecipeSyncedEvent{Added: 5, Updated: 2})

// suggestion module subscribes
bus.Subscribe("recipe.synced", func(e Event) {
    cache.Flush() // invalidate recipe list cache
})
```

**Microservices-ready:** When extracting to microservices, replace `bus.Publish()` with Kafka/RabbitMQ produce, and `bus.Subscribe()` with consumer subscriptions. The event names and payload shapes stay the same.

### Cache: Redis → In-Memory (go-cache)

```go
import "github.com/patrickmn/go-cache"

// 5 min default TTL, cleanup every 10 min
c := cache.New(5*time.Minute, 10*time.Minute)

// Usage
c.Set("recipes:list", recipes, cache.DefaultExpiration)
if val, found := c.Get("recipes:list"); found {
    recipes = val.([]model.Dish)
}
```

**Limitation:** Cache is per-process and lost on restart. Acceptable for MVP since cold cache refills quickly via Notion sync on startup. Cache is also invalidated when `recipe.synced` event fires.

### Search: Elasticsearch → PostgreSQL Full-Text Search

PostgreSQL has built-in full-text search that works well for our scale.

```sql
-- Add search vector column to recipes
ALTER TABLE recipes ADD COLUMN search_vector tsvector;

-- Create index
CREATE INDEX idx_recipes_search ON recipes USING GIN(search_vector);

-- Populate (trigger on insert/update)
CREATE OR REPLACE FUNCTION update_recipe_search() RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('simple', COALESCE(NEW.name, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(NEW.description, '')), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER recipe_search_trigger
  BEFORE INSERT OR UPDATE ON recipes
  FOR EACH ROW EXECUTE FUNCTION update_recipe_search();

-- Search query (supports Vietnamese without diacritics via 'simple' config)
SELECT * FROM recipes
WHERE search_vector @@ plainto_tsquery('simple', 'pho bo')
ORDER BY ts_rank(search_vector, plainto_tsquery('simple', 'pho bo')) DESC
LIMIT 20;
```

**For Vietnamese diacritics support**, also add a `unaccented` column:

```sql
-- Install unaccent extension (available in Supabase)
CREATE EXTENSION IF NOT EXISTS unaccent;

ALTER TABLE recipes ADD COLUMN name_unaccent VARCHAR(255);

-- Trigger to auto-populate
CREATE OR REPLACE FUNCTION update_unaccent() RETURNS trigger AS $$
BEGIN
  NEW.name_unaccent := unaccent(NEW.name);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Search both accented and unaccented
SELECT * FROM recipes
WHERE search_vector @@ plainto_tsquery('simple', 'pho bo')
   OR name_unaccent ILIKE '%pho bo%'
LIMIT 20;
```

### File Storage: S3 → Supabase Storage

Recipe images originate from Notion (uploaded by admins). During sync, the SyncService downloads images from Notion's CDN and re-uploads them to Supabase Storage for stable, permanent URLs (Notion's signed URLs expire).

```go
// During sync: download from Notion, upload to Supabase Storage
notionImageURL := block.Image.File.URL // temporary signed URL
imageBytes := downloadImage(notionImageURL)

path := fmt.Sprintf("recipes/%s/cover.jpg", dish.ID)
supabaseStorage.Upload("recipe-images", path, imageBytes)

// Public URL for frontend
publicURL := supabaseStorage.GetPublicURL("recipe-images", path)
```

---

## 3. Free Tier Limitations & Workarounds

### Render: Cold Starts (15 min inactivity → sleep)

**Problem:** <cite>Free web services spin down after 15 minutes without traffic. Spin-up takes about 1 minute.</cite>

**Workarounds:**
1. **UptimeRobot** (free) — ping `/healthz` endpoint every 5 minutes to keep alive
2. **Cron-job.org** (free) — alternative ping service
3. **Accept it for MVP** — first visit after sleep takes ~30-60s, subsequent requests are fast

```javascript
// healthz endpoint
app.get('/healthz', (req, res) => res.json({ status: 'ok', timestamp: Date.now() }));
```

### Supabase: 500 MB Storage Limit

**Estimation for VietFood:**
- 500 recipes × ~5 KB per recipe (text data) = ~2.5 MB
- 500 nutrition records × ~0.5 KB = ~0.25 MB
- Categories, tags, users, etc. = ~1 MB
- **Total text data: ~5 MB** (only 1% of limit)

Even with 2000+ recipes, you have plenty of room. Images are stored in Supabase Storage (separate 1 GB limit), not in the database.

### Supabase: Project Pauses After 7 Days Inactivity

**Workaround:** Same as Render — use a ping service to make at least 1 API call per day.

### Render Free Postgres: Expires After 30 Days

**Solution:** Don't use Render's free Postgres. Use **Supabase** as the database instead — it doesn't expire.

---

## 4. Project Structure (Modular Monolith — organized by feature)

The Go backend is organized **by bounded context**, not by layer. Each context folder contains its own model, port, service, handler, repository, and DTO. See `CLAUDE.md` for the full directory tree.

```
whatdish/
├── apps/
│   └── web/                    # Next.js frontend (→ Vercel)
│       ├── app/
│       ├── components/
│       └── package.json
│
├── backend/                    # Go API server (→ Render)
│   ├── cmd/server/main.go      # Wire deps, start server
│   ├── internal/
│   │   ├── recipe/             # BC: dishes, ingredients, steps, categories, sync
│   │   ├── suggestion/         # BC: random, calorie, group strategies
│   │   ├── user/               # BC: auth, profiles, allergies
│   │   ├── nutrition/          # BC: macros, TDEE, calorie filtering
│   │   ├── engagement/         # BC: favorites, views, ratings
│   │   ├── shared/             # Cross-cutting: middleware, events, cache, errors
│   │   └── platform/           # External adapters: notion/, supabase/
│   ├── migrations/
│   ├── seeds/
│   ├── Dockerfile
│   ├── Makefile
│   └── go.mod
│
├── docs/
│   ├── USE_CASES.md
│   ├── ARCHITECTURE.md
│   ├── INFRASTRUCTURE.md       # This file
│   ├── DATABASE.md
│   ├── API_DESIGN.md
│   └── notion-setup.md
│
└── README.md
```

**Why by feature?** When extracting a bounded context to a microservice, you copy `internal/{context}/` into a new repo — it has everything (model, port, service, handler, repo, SQL queries). No untangling imports from 6 layer folders.

---

## 5. Environment Variables

### Frontend (Vercel)

```env
NEXT_PUBLIC_SUPABASE_URL=https://xxxxx.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJ...
NEXT_PUBLIC_API_URL=https://whatdish-api.onrender.com
```

### Backend (Render)

```env
PORT=8080
DATABASE_URL=postgresql://...
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_SERVICE_ROLE_KEY=eyJ...
SUPABASE_ANON_KEY=eyJ...
NOTION_API_KEY=secret_...
NOTION_DATABASE_ID=...
SYNC_INTERVAL_MINUTES=30
FRONTEND_URL=https://whatdish.vercel.app
ENVIRONMENT=production
LOG_LEVEL=info
```

---

## 6. Migration Path: Free → Paid → Microservices

### Stage 1: Free Tier (current)

All services on free tiers. Accept cold starts and limitations.

### Stage 2: Minimal Paid ($25-32/month)

When you get real users:
- Supabase Pro ($25/month): 8 GB DB, no pausing, daily backups
- Render Starter ($7/month): no sleep, better performance
- Total: ~$32/month

### Stage 3: Scale Services ($100-200/month)

When traffic grows:
- Add Redis (Render Key Value or Upstash)
- Add Elasticsearch (Bonsai.io free tier or Elastic Cloud)
- Upgrade Render instance for more RAM/CPU

### Stage 4: Microservices

When team grows (3+ devs) or single service bottlenecks:
- Extract User Service first (auth is cross-cutting)
- Extract Recipe Service (read-heavy, needs independent scaling)
- Add Kafka/RabbitMQ for async events
- Deploy on Kubernetes

The **Bounded Context separation in code** (modules/) ensures this extraction is a deployment change, not a rewrite.

---

## 7. Docs Validity Checklist

| Document | Still Valid? | Notes |
|----------|-------------|-------|
| `USE_CASES.md` | **100% valid** | Use cases are infrastructure-agnostic |
| `ARCHITECTURE.md` | **Valid as target** | Describes the end-state microservices architecture |
| `DATABASE.md` | **95% valid** | Same tables, just in 1 DB instead of 5. Drop `processed_events` table for now |
| `API_DESIGN.md` | **90% valid** | Same endpoints. Auth endpoints change to proxy Supabase Auth. Rate limiting simplified |
| `INFRASTRUCTURE.md` | **This file** | Replaces infra sections of ARCHITECTURE.md for current deployment |