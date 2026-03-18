# CLAUDE.md — WhatDish (Backend)

## Project overview

Go REST API for the **WhatDish** (What to eat today?) web app. Smart food suggestion system: random dish picker, calorie-based recommendations, group meal planning, with detailed Vietnamese cooking recipes. The frontend is a separate Next.js app on Vercel communicating via HTTP.

### Core features (mapped to Use Cases)

| ID | Feature | Description |
|----|---------|-------------|
| UC01 | Random dish | Pick a random dish, exclude recent history |
| UC02 | Calorie-based suggestion | Suggest dishes within a calorie target (±15%) |
| UC03 | Group meal suggestion | Balanced combo for N people (main + soup + sides) |
| UC04 | Filter dishes | Multi-criteria: dish type, region, ingredient, time, difficulty |
| UC05 | Recipe detail | Full recipe with ingredients, steps, nutrition |
| UC06 | Search | Full-text search (PostgreSQL tsvector) |
| UC07 | Auth | Register/login via Supabase Auth |
| UC08 | Nutrition profile | TDEE calculation (Mifflin-St Jeor), personal goals |
| UC09 | Favorites | Save/unsave dishes |
| UC10 | Suggestion history | Track past suggestions, avoid repeats |
| UC11 | Admin: Recipes | Manage recipes in **Notion**, sync to DB via admin endpoint |
| UC12 | Admin: Users | Manage user accounts |
| UC13 | Admin: Categories | Manage categories, tags (via seed data or future admin UI) |

Full use case flows: see `docs/USE_CASES.md`

## Tech stack

- **Language:** Go 1.23+
- **HTTP framework:** Gin
- **Database:** PostgreSQL (hosted on Supabase, single DB)
- **Auth:** Supabase Auth (50K MAU free tier)
- **File storage:** Supabase Storage (1 GB free tier)
- **DB driver:** pgx/v5
- **Query generation:** sqlc (type-safe SQL → Go)
- **Migration:** golang-migrate
- **Config:** envconfig
- **Logging:** slog (stdlib)
- **Validation:** go-playground/validator
- **Content CMS:** Notion (free) — admin manages recipes via Notion database
- **Cache:** In-memory (sync.Map or go-cache) — Redis replacement for free tier
- **Search:** PostgreSQL full-text search (tsvector + unaccent)
- **Testing:** stdlib testing + testify + testcontainers-go
- **Deployment:** Render (free tier, Docker)

## Development workflow

Follow this loop for every feature or change:

```
┌─────────────────────────────────────────────────────────────────┐
│  1. IMPLEMENT                                                   │
│     Build the feature                                           │
│                                                                 │
│  2. QUALITY LOOP (repeat until all green)                       │
│     ┌─────────────────────────────────────────────────────┐     │
│     │  a. Add/update unit tests (service layer)           │     │
│     │  b. Add/update integration tests (repository layer) │     │
│     │  c. Add/update e2e tests (handler layer, full HTTP) │     │
│     │  d. Run & verify all tests pass                     │     │
│     │     → make test                                     │     │
│     │     → make test-integration                         │     │
│     │  e. Review code against best practices              │     │
│     │     → SOLID principles                              │     │
│     │     → interface segregation (small interfaces)      │     │
│     │     → proper error wrapping & handling              │     │
│     │     → no business logic in handlers                 │     │
│     │     → no direct DB access in services               │     │
│     │  f. Run lint & CI tools                             │     │
│     │     → make lint  (golangci-lint)                    │     │
│     │     → make vet   (go vet)                           │     │
│     │     → make fmt-check                                │     │
│     │  g. If anything fails → fix and repeat from (a)     │     │
│     └─────────────────────────────────────────────────────┘     │
│                                                                 │
│  3. UPDATE DOCS                                                 │
│     - Update README.md if public API or setup changed           │
│     - Update this CLAUDE.md if architecture/conventions changed │
│     - Update diagrams (Mermaid in docs/) for any new            │
│       data flows, system interactions, or state machines        │
│     - Update API contract if endpoints changed                  │
│     - Regenerate sqlc if queries changed: make sqlc             │
└─────────────────────────────────────────────────────────────────┘
```

### Test expectations by layer

| Layer | Tool | What to test | Coverage target |
|-------|------|-------------|-----------------|
| Unit | testing + testify | Service logic in `{context}/service.go`, TDEE calc, strategies | 80%+ |
| Integration | testcontainers-go | Repository queries in `{context}/repository.go` against real Postgres | All queries |
| E2E / Handler | httptest + testify | Full HTTP cycle in `{context}/handler.go`, middleware, auth, errors | All endpoints |
| Platform | testing + testify | Notion parser in `platform/notion/parser.go`, every block type | All block types |

### CI checks (must all pass before merge)

```bash
make fmt-check       # gofmt formatting check
make vet             # go vet static analysis
make lint            # golangci-lint (strict config)
make test            # Unit tests
make test-integration # Integration tests (needs Docker)
make test-cover      # All tests + coverage report
make build           # Binary compiles successfully
```

## Architecture

### Pattern: Clean Architecture (Modular Monolith, Microservices-Ready)

The codebase follows Clean Architecture with **Bounded Context modules**. Each module (recipe, suggestion, user, nutrition, engagement) is self-contained and communicates via in-process interfaces. When scaling is needed, modules can be extracted into separate services with minimal code changes.

```
┌──────────────────────────────────────────────────────────────┐
│                        ADAPTERS (outer)                      │
│                                                              │
│  ┌──────────┐ ┌───────────┐ ┌───────────┐ ┌──────────────┐ │
│  │ HTTP     │ │ PostgreSQL│ │ Notion API│ │ Supabase     │ │
│  │ (Gin)    │ │ (pgx/sqlc)│ │ (content) │ │ Auth+Storage │ │
│  │ handler/ │ │ postgres/ │ │ notion/   │ │ auth/storage/│ │
│  └────┬─────┘ └─────┬─────┘ └─────┬─────┘ └──────┬───────┘ │
│       │              │             │               │         │
│───────┼──────────────┼─────────────┼───────────────┼─────────│
│       ▼              ▼             ▼               ▼         │
│  ┌─────────────────────────────────────────────────────┐     │
│  │               PORTS (interfaces)                     │     │
│  │  DishRepository     NutritionRepository              │     │
│  │  SuggestionRepo     EngagementRepository             │     │
│  │  UserRepository     AuthProvider                     │     │
│  │  ContentSource      ImageStorage                     │     │
│  │  EventBus                                            │     │
│  └──────────────────────┬──────────────────────────────┘     │
│                         │                                     │
│─────────────────────────┼─────────────────────────────────────│
│                         ▼                                     │
│  ┌─────────────────────────────────────────────────────┐     │
│  │                 CORE (inner)                         │     │
│  │                                                      │     │
│  │  model/        Domain entities                       │     │
│  │  service/      Business logic per bounded context    │     │
│  │  service/      SyncService — Notion → PostgreSQL     │     │
│  │  port/         Interface definitions                 │     │
│  └─────────────────────────────────────────────────────┘     │
└──────────────────────────────────────────────────────────────┘
```

### Bounded Contexts (each = one folder under internal/)

```
internal/
├── recipe/          → model + port + service + handler + repo + dto + queries
├── suggestion/      → model + port + service + strategies + handler + repo + dto
├── user/            → model + port + services + handler + repo + dto
├── nutrition/       → model + port + service + tdee + handler + repo + dto
├── engagement/      → model + port + service + handler + repo + dto
├── shared/          → middleware, event bus, cache, database, errors (cross-cutting)
└── platform/        → notion/, supabase/ (external service adapters)
```

```
┌──────────────────────────────────────────────────────────────┐
│                    Single Supabase DB                         │
│            (table prefix separation per context)             │
│                                                              │
│  ┌─────────┐ ┌──────────┐ ┌──────┐ ┌─────────┐ ┌────────┐ │
│  │ recipe  │ │suggestion│ │ user │ │nutrition│ │engage- │ │
│  │         │ │          │ │      │ │         │ │ment    │ │
│  │ dishes  │ │ sessions │ │profi-│ │recipe_  │ │favor-  │ │
│  │ ingredi-│ │ configs  │ │les   │ │nutrition│ │ites    │ │
│  │ ents    │ │ exclu-   │ │aller-│ │goals    │ │views   │ │
│  │ steps   │ │ sions    │ │gies  │ │         │ │ratings │ │
│  │ categor-│ │          │ │      │ │         │ │        │ │
│  │ ies     │ │          │ │      │ │         │ │        │ │
│  │ tags    │ │          │ │      │ │         │ │        │ │
│  │ sync_log│ │          │ │      │ │         │ │        │ │
│  └─────────┘ └──────────┘ └──────┘ └─────────┘ └────────┘ │
└──────────────────────────────────────────────────────────────┘
```

Each context folder owns its tables. Cross-context data access goes through **port interfaces**, not direct SQL joins or Go imports.

### Dependency rule

Dependencies point INWARD only (within each bounded context):
- `handler.go` depends on `service.go` (via interface in `port.go`)
- `service.go` depends on `port.go` (interfaces) and `model.go`
- `repository.go` implements `port.go` interfaces
- `model.go` depends on NOTHING

Cross-context dependencies go through interfaces:
- `suggestion/port.go` defines `DishReader` interface
- `recipe/service.go` satisfies `DishReader` without importing `suggestion/`
- Wired together in `cmd/server/main.go`

### Design patterns used

| Pattern | Where | Why |
|---------|-------|-----|
| **Repository** | `internal/port/` + `internal/repository/` | Abstracts data access behind interfaces. Services never touch SQL directly. |
| **Strategy** | `internal/service/suggestion/` | Three suggestion strategies (Random, ByCalories, ByGroup) implement `SuggestionStrategy` interface. |
| **Factory** | `internal/service/suggestion/` | `NewStrategy(type)` returns the correct suggestion implementation. |
| **Adapter** | `internal/repository/` | PostgreSQL adapters implement port interfaces. Supabase Auth adapter implements `AuthProvider`. Notion adapter implements `ContentSource`. |
| **Anti-corruption layer** | `internal/repository/notion/` | Translates Notion API responses (blocks, properties, rich text) into clean domain models. Notion's data format never leaks into service/model layers. |
| **Cache-Aside** | `internal/service/recipe/` | Notion is source of truth for recipes, PostgreSQL is the read cache. SyncService periodically pulls from Notion and upserts into PostgreSQL. All reads served from PostgreSQL. |
| **DTO / Data mapping** | `internal/dto/` | Request/response DTOs decouple HTTP from domain models. |
| **Middleware Chain** | `internal/middleware/` | CORS, logging, recovery, auth verification. |
| **Dependency Injection** | `cmd/server/main.go` | Constructor injection — all dependencies wired manually at startup. No DI framework. |
| **Event Bus** | `internal/event/` | In-process EventEmitter for cross-module communication. Swap to Kafka/RabbitMQ when extracting services. |
| **Cache-Aside** | `internal/cache/` | In-memory cache (go-cache) for recipe lists, popular dishes. Swap to Redis later. |

### Project structure

```
.
├── cmd/
│   └── server/
│       └── main.go                     # Entry point: wire deps, start server
├── internal/
│   ├── config/
│   │   └── config.go                   # Env parsing, config struct
│   │
│   ├── model/                          # Domain entities (ZERO external deps)
│   │   ├── dish.go                     #   Dish, Ingredient, Step, Category, Tag
│   │   ├── suggestion.go               #   SuggestionSession, SuggestionConfig
│   │   ├── user.go                     #   User, UserProfile, UserAllergy
│   │   ├── nutrition.go                #   RecipeNutrition, NutritionGoal, TDEE
│   │   └── engagement.go               #   Favorite, ViewHistory, Rating
│   │
│   ├── port/                           # Interface definitions (ports)
│   │   ├── dish_repository.go          #   DishRepository, CategoryRepository
│   │   ├── suggestion_repository.go    #   SuggestionSessionRepo, ExclusionRepo
│   │   ├── user_repository.go          #   UserProfileRepo, UserAllergyRepo
│   │   ├── nutrition_repository.go     #   NutritionRepo, NutritionGoalRepo
│   │   ├── engagement_repository.go    #   FavoriteRepo, ViewHistoryRepo, RatingRepo
│   │   ├── content_source.go           #   ContentSource (Notion abstraction)
│   │   ├── auth_provider.go            #   AuthProvider (Supabase Auth abstraction)
│   │   ├── image_storage.go            #   ImageStorage (Supabase Storage abstraction)
│   │   └── event_bus.go                #   EventBus interface
│   │
│   ├── service/                        # Business logic per bounded context
│   │   ├── recipe/
│   │   │   ├── dish_service.go         #   Read operations, filtering, search (reads from PG cache)
│   │   │   ├── dish_service_test.go
│   │   │   ├── category_service.go     #   Category management
│   │   │   ├── search_service.go       #   Full-text search via tsvector
│   │   │   ├── sync_service.go         #   Notion → PostgreSQL sync orchestration
│   │   │   └── sync_service_test.go
│   │   ├── suggestion/
│   │   │   ├── suggestion_service.go   #   Orchestrates all suggestion types
│   │   │   ├── suggestion_service_test.go
│   │   │   ├── strategy.go             #   SuggestionStrategy interface
│   │   │   ├── random_strategy.go      #   UC01: random with history exclusion
│   │   │   ├── random_strategy_test.go
│   │   │   ├── calorie_strategy.go     #   UC02: filter by calorie range
│   │   │   ├── calorie_strategy_test.go
│   │   │   ├── group_strategy.go       #   UC03: balanced combo for N people
│   │   │   └── group_strategy_test.go
│   │   ├── user/
│   │   │   ├── auth_service.go         #   Register, login, verify (delegates to Supabase Auth)
│   │   │   ├── profile_service.go      #   Nutrition profile CRUD
│   │   │   └── profile_service_test.go
│   │   ├── nutrition/
│   │   │   ├── nutrition_service.go    #   Recipe nutrition CRUD, batch query
│   │   │   ├── tdee_calculator.go      #   Mifflin-St Jeor formula
│   │   │   ├── tdee_calculator_test.go
│   │   │   └── calorie_filter.go       #   Filter recipe IDs by calorie range
│   │   └── engagement/
│   │       ├── favorite_service.go     #   Add/remove/list favorites
│   │       ├── view_service.go         #   Record views, track history
│   │       └── stats_service.go        #   Aggregate stats (view count, fav count)
│   │
│   ├── handler/                        # HTTP handlers (thin adapter layer)
│   │   ├── dish_handler.go             #   GET /dishes, /dishes/:id, /dishes/random, /dishes/search
│   │   ├── dish_handler_test.go
│   │   ├── suggestion_handler.go       #   POST /suggestions/random, /by-calories, /by-group
│   │   ├── suggestion_handler_test.go
│   │   ├── auth_handler.go             #   POST /auth/register, /auth/login (proxy to Supabase)
│   │   ├── user_handler.go             #   GET /users/me, PUT /users/me/profile
│   │   ├── nutrition_handler.go        #   GET /nutrition/recipes/:id, POST /nutrition/calculate-tdee
│   │   ├── favorite_handler.go         #   POST/DELETE /favorites, GET /favorites
│   │   ├── category_handler.go         #   GET /categories
│   │   ├── admin_handler.go            #   Admin CRUD endpoints (protected)
│   │   ├── health_handler.go           #   GET /health
│   │   └── router.go                   #   Route registration, middleware setup
│   │
│   ├── repository/                     # Adapter implementations
│   │   ├── postgres/                   #   PostgreSQL adapter (Supabase)
│   │   │   ├── dish_repo.go            #     Implements DishRepository (reads from cache)
│   │   │   ├── dish_repo_test.go
│   │   │   ├── suggestion_repo.go      #     Implements SuggestionSessionRepo
│   │   │   ├── user_repo.go            #     Implements UserProfileRepo
│   │   │   ├── nutrition_repo.go       #     Implements NutritionRepo
│   │   │   ├── engagement_repo.go      #     Implements FavoriteRepo, ViewHistoryRepo
│   │   │   ├── queries/                #     Raw SQL for sqlc
│   │   │   │   ├── dish.sql
│   │   │   │   ├── category.sql
│   │   │   │   ├── ingredient.sql
│   │   │   │   ├── step.sql
│   │   │   │   ├── suggestion.sql
│   │   │   │   ├── user.sql
│   │   │   │   ├── nutrition.sql
│   │   │   │   └── engagement.sql
│   │   │   └── db.go                   #     Connection pool setup
│   │   └── notion/                     #   Notion API adapter (content source)
│   │       ├── client.go               #     Notion API HTTP client
│   │       ├── content_source.go       #     Implements ContentSource port
│   │       ├── parser.go               #     Parse Notion blocks → domain models
│   │       ├── parser_test.go          #     Unit tests for every block type
│   │       └── types.go               #     Notion API response types
│   │
│   ├── auth/                           # Supabase Auth adapter
│   │   └── supabase_auth.go           #   Implements AuthProvider port
│   │
│   ├── storage/                        # Supabase Storage adapter
│   │   └── supabase_storage.go        #   Implements ImageStorage port
│   │
│   ├── event/                          # In-process event bus
│   │   ├── bus.go                     #   EventBus implementation (channels/callbacks)
│   │   └── events.go                  #   Event type constants + payload structs
│   │
│   ├── cache/                          # In-memory cache (replace with Redis later)
│   │   └── memory.go                  #   go-cache wrapper, TTL management
│   │
│   ├── dto/                            # Request/response data transfer objects
│   │   ├── dish_dto.go
│   │   ├── suggestion_dto.go
│   │   ├── user_dto.go
│   │   ├── nutrition_dto.go
│   │   ├── engagement_dto.go
│   │   └── error_dto.go               #   Standardized error response
│   │
│   └── middleware/
│       ├── cors.go
│       ├── logger.go
│       ├── recovery.go
│       └── auth.go                     #   Verify Supabase JWT, extract user
│
├── migrations/                         # golang-migrate SQL files
│   ├── 000001_create_categories.up.sql
│   ├── 000001_create_categories.down.sql
│   ├── 000002_create_recipes.up.sql
│   ├── 000002_create_recipes.down.sql
│   ├── 000003_create_ingredients_steps.up.sql
│   ├── 000004_create_tags.up.sql
│   ├── 000005_create_user_profiles.up.sql
│   ├── 000006_create_nutrition.up.sql
│   ├── 000007_create_suggestions.up.sql
│   ├── 000008_create_engagement.up.sql
│   ├── 000009_add_search_index.up.sql
│   └── 000010_create_sync_logs.up.sql
│
├── seeds/                              # Reference data only (recipes come from Notion sync)
│   ├── 01_categories.sql
│   ├── 02_tags.sql
│   ├── 03_nutrition_goals.sql
│   └── 04_suggestion_configs.sql
│
├── docs/                               # Project documentation
│   ├── USE_CASES.md                   #   Detailed use case flows
│   ├── ARCHITECTURE.md                #   Target microservices architecture
│   ├── INFRASTRUCTURE.md              #   Actual free-tier deployment
│   ├── DATABASE.md                    #   Full database schema
│   ├── API_DESIGN.md                  #   RESTful API endpoints
│   ├── notion-setup.md                #   How to set up Notion database + integration
│   └── diagrams/
│       ├── system-overview.mmd
│       ├── bounded-contexts.mmd
│       ├── suggestion-flow.mmd
│       ├── notion-sync-flow.mmd
│       └── auth-flow.mmd
│
├── sqlc.yaml
├── Dockerfile
├── docker-compose.yml                  # Local dev: postgres + app
├── Makefile
├── .golangci.yml
├── .env.example
└── go.mod
```

### Request lifecycle

```
HTTP Request
    │
    ▼
[Gin Router] (cmd/server/main.go registers routes)
    │
    ├── shared/middleware/cors.go        → CORS headers (allow FRONTEND_URL only)
    ├── shared/middleware/logger.go      → Request logging (slog)
    ├── shared/middleware/recovery.go    → Panic recovery
    ├── shared/middleware/auth.go        → Verify Supabase JWT (optional/required per route)
    │
    ▼
[recipe/handler.go]                      → Parse request, validate DTO
    │                                    → All in same package (recipe/)
    ▼
[recipe/service.go]                      → Business logic, orchestration
    │                                    → Calls repository via port interface
    ▼
[recipe/repository.go]                   → SQL query (sqlc generated)
    │                                    → Convert DB row → domain model
    ▼
[PostgreSQL on Supabase]
    │
    ▼
Response bubbles back up:
    repository → model → service → model → handler → DTO → JSON
```

Note: handler, service, repository are all in the same Go package (e.g. `internal/recipe/`).
They communicate via interfaces defined in `port.go`, not via direct struct access.
This keeps them testable — mock the port interface in unit tests.

## Database

### Single Supabase database, schema separation

All tables live in one Supabase PostgreSQL instance (500 MB free tier). Bounded contexts are separated by **table prefix** convention since Supabase free tier works best with the default `public` schema.

```
recipe_*        → recipes, recipe_ingredients, recipe_steps, categories, tags, recipe_tags
suggestion_*    → suggestion_sessions, suggestion_configs, suggestion_exclusions
user_*          → user_profiles, user_allergies (auth.users managed by Supabase)
nutrition_*     → nutrition_recipe, nutrition_goals
engagement_*    → engagement_favorites, engagement_views, engagement_ratings
```

Full schema: see `docs/DATABASE.md`

### Key tables

```sql
-- Core recipe data (cached from Notion — external_id = Notion page ID)
recipes              -- dish info (name, slug, image, cook_time, difficulty, servings)
recipe_ingredients   -- ingredient list per recipe (amount, unit, note)
recipe_steps         -- cooking steps (step_number, description, duration)
categories           -- multi-type categories (dish_type, region, main_ingredient, meal_type)
tags                 -- flexible tags (weight_loss, quick, kid_friendly...)
recipe_tags          -- M:N junction
sync_logs            -- tracks Notion sync history (added/updated/deleted counts)

-- Suggestion engine
suggestion_sessions  -- each suggestion request = 1 session (type, input_params, result)
suggestion_configs   -- presets for calorie/group modes
suggestion_exclusions -- recent dishes to exclude per user

-- User data (auth.users table managed by Supabase Auth)
user_profiles        -- nutrition profile (gender, age, height, weight, TDEE)
user_allergies       -- ingredients user is allergic to or dislikes

-- Nutrition
nutrition_recipe     -- per-serving macros (calories, protein, carbs, fat, fiber)
nutrition_goals      -- preset goals (weight loss, muscle gain, maintenance)

-- Engagement
engagement_favorites -- user bookmarks
engagement_views     -- view history (source tracking)
engagement_ratings   -- 1-5 star ratings
```

### Key indexes

```sql
idx_recipes_status          -- filter published only (partial: WHERE deleted_at IS NULL)
idx_recipes_dish_type       -- filter by dish type category
idx_recipes_region          -- filter by region
idx_recipes_difficulty      -- filter by difficulty
idx_recipes_search          -- GIN index on tsvector for full-text search
idx_nutrition_calories      -- range queries for calorie-based suggestions
idx_favorites_user          -- user's favorite list
idx_suggestion_sessions_user -- user's suggestion history
idx_recipes_external_id     -- unique Notion page ID for sync upsert
```

### Full-text search setup

```sql
-- Vietnamese search support via 'simple' config + unaccent extension
CREATE EXTENSION IF NOT EXISTS unaccent;

ALTER TABLE recipes ADD COLUMN search_vector tsvector;
CREATE INDEX idx_recipes_search ON recipes USING GIN(search_vector);

-- Auto-update trigger
CREATE FUNCTION update_recipe_search() RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('simple', unaccent(COALESCE(NEW.name, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(NEW.description, ''))), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

## Notion → PostgreSQL sync (recipe content management)

### Why Notion?

Admins manage recipes in a **Notion database** — it's free, has rich text editing, supports images, and is easy for non-developers to use. The Go backend syncs Notion data into PostgreSQL as a read cache. Users never hit the Notion API; all reads are served from PostgreSQL.

### Sync flow

```
[Startup / Cron every N minutes / Manual trigger POST /admin/sync]
    │
    ▼
[SyncService.Sync()]
    │
    ├── Fetch all pages from Notion database (ContentSource port)
    │       → Filter: Status == "Published"
    │       → Paginate through all results
    │
    ├── For each Notion page:
    │       → Parse properties → Dish model (name, category, difficulty, etc.)
    │       → Fetch page blocks (children)
    │       → Parse blocks → Ingredients[] + Steps[] + Tips
    │       → Extract nutrition data if present
    │
    ├── Upsert into PostgreSQL (DishRepository)
    │       → Match by external_id (Notion page ID)
    │       → Insert new, update changed, soft-delete removed
    │       → Upsert ingredients + steps in transaction
    │       → Update search_vector (tsvector) via trigger
    │
    ├── Publish event: recipe.synced → EventBus
    │       → Cache invalidation (in-memory cache cleared)
    │
    └── Log sync result to sync_logs table
```

### Content parsing rules (Notion blocks → domain models)

```
Notion block type         → Domain model
─────────────────────────────────────────
bulleted_list_item        → Ingredient (parse "200g thịt bò" → amount=200, unit=g, name=thịt bò)
numbered_list_item        → Step (step_number = list position)
heading_2 / heading_3     → Section separator (ingredients vs steps vs tips)
callout                   → Tip text
image                     → Step image or dish cover image
table                     → Ingredient table (if in ingredients section)
paragraph with "(X phút)" → timer_secs extraction via regex
```

### Notion database properties (expected schema)

```
Name          → text       → dish name (Vietnamese)
Slug          → text       → URL-friendly slug
Category      → select     → maps to categories table (dish_type)
Region        → select     → Bắc / Trung / Nam
Main Ingredient → select   → Gà / Bò / Heo / Hải sản / Chay...
Meal Type     → select     → Sáng / Trưa / Tối / Ăn vặt
Difficulty    → select     → EASY / MEDIUM / HARD
Prep Time     → number     → minutes
Cook Time     → number     → minutes
Servings      → number     → default serving count
Calories      → number     → kcal per serving (if known)
Protein       → number     → grams per serving
Carbs         → number     → grams per serving
Fat           → number     → grams per serving
Tags          → multi_select → quick, healthy, kid_friendly...
Cover         → files      → dish photo
Status        → select     → Draft / Published / Archived
```

**Setup guide:** see `docs/notion-setup.md` for how to create the Notion database and integration.

## API endpoints

Full API contract: see `docs/API_DESIGN.md`

### Auth levels

| Symbol | Meaning |
|--------|---------|
| 🔓 | Public — no auth required |
| 🔐 | User — Supabase JWT required |
| 🔑 | Admin — JWT + role = ADMIN |

### Recipes (Recipe Context — read-only cache from Notion)

```
🔓 GET  /api/v1/recipes
     Query: ?dish_type=braised&region=northern&difficulty=EASY
            &max_cook_time=30&tags=quick&page=1&limit=20
     Response: { data: Recipe[], pagination: {...} }

🔓 GET  /api/v1/recipes/:id
     Response: { data: RecipeDetail }  (includes ingredients, steps, nutrition)

🔓 GET  /api/v1/recipes/random
     Query: ?exclude_ids=uuid1,uuid2&dish_type=...
     Response: { data: RecipeDetail }

🔓 GET  /api/v1/recipes/search
     Query: ?q=pho+bo&page=1&limit=20
     Response: { data: Recipe[], pagination: {...} }

🔓 GET  /api/v1/recipes/autocomplete
     Query: ?q=pho&limit=5
     Response: { data: [{id, name, type}] }

🔓 GET  /api/v1/categories
     Query: ?type=DISH_TYPE
     Response: { data: Category[] }
```

Note: No POST/PUT/DELETE for recipes — content is managed in Notion and synced via admin endpoint.

### Suggestions (Suggestion Context)

```
🔓 POST /api/v1/suggestions/random
     Body: { filters: { dish_type?, region?, max_cook_time? } }
     Response: { data: { session_id, recipe, suggestion_type } }

🔓 POST /api/v1/suggestions/by-calories
     Body: { target_calories, meal_type, tolerance_pct?, filters? }
     Response: { data: { session_id, recipes[], total_calories, calorie_diff } }

🔓 POST /api/v1/suggestions/by-group
     Body: { group_size, group_type, meal_type, budget?, filters? }
     Response: { data: { session_id, dishes[], shopping_list[], total_nutrition } }

🔐 GET  /api/v1/suggestions/history
     Query: ?session_type=BY_CALORIES&page=1&limit=20
     Response: { data: SuggestionSession[], pagination, stats }
```

### Auth & User (User Context)

```
🔓 POST /api/v1/auth/register    — proxy to Supabase Auth signUp
🔓 POST /api/v1/auth/login       — proxy to Supabase Auth signIn
🔓 POST /api/v1/auth/refresh     — proxy to Supabase Auth refreshSession
🔓 POST /api/v1/auth/forgot-password
🔓 POST /api/v1/auth/reset-password

🔐 GET  /api/v1/users/me
🔐 PUT  /api/v1/users/me/profile
     Body: { gender, age, height_cm, weight_kg, activity_level, goal, allergies[] }
     Response: { data: { profile, calculated: { bmr, tdee, daily_target, meal_targets } } }

🔑 GET  /api/v1/users              — list users (admin)
🔑 PUT  /api/v1/users/:id/status   — suspend/activate (admin)
```

### Nutrition (Nutrition Context)

```
🔓 GET  /api/v1/nutrition/recipes/:id
🔓 GET  /api/v1/nutrition/recipes?ids=uuid1,uuid2
🔓 GET  /api/v1/nutrition/goals
🔓 POST /api/v1/nutrition/calculate-tdee
     Body: { gender, age, height_cm, weight_kg, activity_level, goal }
     Response: { data: { bmr, tdee, daily_target, meal_breakdown } }

🔑 POST /api/v1/nutrition/recipes/:id   — add/update (admin)
```

### Engagement (Engagement Context)

```
🔐 POST   /api/v1/favorites          — add favorite
🔐 DELETE /api/v1/favorites/:recipe_id — remove
🔐 GET    /api/v1/favorites           — list
🔐 GET    /api/v1/favorites/check?recipe_ids=uuid1,uuid2

🔓 POST   /api/v1/views              — record view (session-based)
```

### Health & Admin

```
🔓 GET  /health → { status: "ok", timestamp: "..." }

🔑 POST /api/v1/admin/sync
     Response: { added: int, updated: int, deleted: int, duration: "2.3s" }
     Purpose: Manually trigger Notion → PostgreSQL sync

🔑 GET  /api/v1/admin/sync/status
     Response: { last_sync: "...", next_sync: "...", total_recipes: int }
```

## Core business logic

### TDEE calculation (Mifflin-St Jeor)

```go
// internal/service/nutrition/tdee_calculator.go

func CalculateBMR(gender string, weightKg, heightCm float64, age int) float64 {
    if gender == "MALE" {
        return 10*weightKg + 6.25*heightCm - 5*float64(age) + 5
    }
    return 10*weightKg + 6.25*heightCm - 5*float64(age) - 161
}

var activityMultipliers = map[string]float64{
    "SEDENTARY":    1.2,
    "LIGHT":        1.375,
    "MODERATE":     1.55,
    "ACTIVE":       1.725,
    "VERY_ACTIVE":  1.9,
}

var goalAdjustments = map[string]float64{
    "LOSE_WEIGHT": 0.80,
    "MAINTAIN":    1.00,
    "GAIN_WEIGHT": 1.15,
}

func CalculateTDEE(bmr float64, activityLevel, goal string) float64 {
    return bmr * activityMultipliers[activityLevel] * goalAdjustments[goal]
}
```

### Group suggestion balancing (UC03)

```go
// internal/service/suggestion/group_strategy.go

// Constraints for a balanced group meal:
// 1. Determine dish count: min(ceil(groupSize/2)+1, 5), min=3
// 2. Slot allocation: 1 main (protein), 1 soup, rest = sides
// 3. No duplicate main ingredients across dishes
// 4. Diverse cooking methods (no two braised, no two fried, etc.)
// 5. Balanced flavor profiles
// 6. Scale servings to group size
// 7. Merge ingredients → shopping list
```

### Random with exclusion (UC01)

```go
// internal/service/suggestion/random_strategy.go

// 1. Get all published recipe IDs from cache (or DB)
// 2. If user is authenticated:
//    a. Fetch exclusion list (recipe IDs from last 7 days)
//    b. Remove excluded IDs from candidate set
// 3. Apply optional filters (dish_type, region, etc.)
// 4. If no candidates remain → reset exclusions, retry
// 5. crypto/rand to pick one from remaining set
// 6. Save to suggestion_sessions + exclusion_rules (if user)
```

## Auth flow (Supabase Auth)

```
Client                    Backend (Go)              Supabase Auth
  │                           │                          │
  ├─ POST /auth/login ───────►│                          │
  │   {email, password}       │                          │
  │                           ├─ supabase.SignIn() ─────►│
  │                           │                          │
  │                           │◄── {access_token,        │
  │                           │     refresh_token,       │
  │                           │     user}                │
  │◄── {tokens, user} ───────┤                          │
  │                           │                          │
  ├─ GET /api/v1/favorites ──►│                          │
  │   Authorization: Bearer   │                          │
  │                           ├─ supabase.GetUser() ────►│
  │                           │◄── {user} ──────────────┤│
  │                           │                          │
  │                           │ [auth.middleware extracts │
  │                           │  user_id from JWT]       │
  │                           │                          │
  │◄── {data: [...]} ────────┤                          │
```

The Go backend does NOT store passwords or manage JWTs. Supabase Auth handles all of that. The backend only:
1. Proxies auth requests to Supabase
2. Verifies JWTs on protected routes via Supabase's `auth.getUser()`
3. Stores additional profile data (nutrition, allergies) in `user_profiles`

## Coding conventions

### Go style
- Follow standard Go conventions: `gofmt`, `go vet`, Effective Go
- Error handling: ALWAYS check errors, wrap with `fmt.Errorf("context: %w", err)`
- Never use `panic` in business logic — only for unrecoverable startup errors
- Use `slog` for structured logging with key-value pairs
- Pass `context.Context` as the first parameter for any function with I/O
- Interfaces are defined where they are CONSUMED (in `port/`), not where implemented

### Naming
- Packages: lowercase, single word matching bounded context (`recipe`, `suggestion`, `user`)
- Files: lowercase, snake_case (`sync_service.go`, `random_strategy.go`)
- Each BC folder uses short file names: `model.go`, `port.go`, `service.go`, `handler.go`, `repository.go`, `dto.go`
- Specialized files get descriptive names: `tdee.go`, `calorie_strategy.go`, `sync_service.go`
- Interfaces: describe behavior (`DishRepository`, `SuggestionStrategy`) — no `I` prefix
- Exported: PascalCase; unexported: camelCase
- Receivers: 1-2 chars (`s` for service, `h` for handler, `r` for repository)
- Test files: `*_test.go` in same package

### Error handling

```go
// Custom error types for business errors
var (
    ErrNotFound      = errors.New("resource not found")
    ErrConflict      = errors.New("resource conflict")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrForbidden     = errors.New("forbidden")
    ErrValidation    = errors.New("validation failed")
)

// HTTP status mapping in handler:
// ErrValidation   → 400
// ErrUnauthorized → 401
// ErrForbidden    → 403
// ErrNotFound     → 404
// ErrConflict     → 409
// default         → 500
```

### Cross-module communication

**Rule: bounded contexts NEVER import each other directly.** Use interfaces.

```go
// suggestion/port.go — defines what suggestion needs from recipe
type DishReader interface {
    GetByID(ctx context.Context, id uuid.UUID) (*recipe.Dish, error)  // NO! imports recipe
    GetByID(ctx context.Context, id uuid.UUID) (*Dish, error)         // YES! own model or shared type
}

// suggestion/port.go — defines what suggestion needs from nutrition
type CalorieProvider interface {
    GetRecipeIDsByCalorieRange(ctx context.Context, min, max int) ([]uuid.UUID, error)
}
```

```go
// In-process event bus (replace with Kafka/RabbitMQ when extracting services)

// recipe/sync_service.go publishes
s.eventBus.Publish("recipe.synced", shared.RecipeSyncedEvent{Added: 5, Updated: 2})

// suggestion/service.go subscribes
s.eventBus.Subscribe("recipe.synced", func(e shared.Event) {
    s.cache.Flush() // invalidate suggestion cache
})

// Event names follow: {context}.{action}
// recipe.synced
// user.profile_updated, user.deleted
// engagement.favorited, engagement.viewed
```

**Shared types** for cross-context data live in `shared/` — not in any BC. Keep them minimal (IDs, event payloads, pagination structs).

### Project conventions
- **Organized by feature (bounded context), not by layer** — each BC has its own model/port/service/handler/repo/dto
- SQL queries co-located with their BC: `internal/recipe/queries/*.sql` → run `make sqlc`
- Migrations in top-level `migrations/`, numbered sequentially: `000001_`, `000002_`...
- Seed data in `seeds/`, run via `make seed`
- External service adapters in `internal/platform/` (notion, supabase) — shared across BCs
- Cross-cutting concerns in `internal/shared/` (middleware, event bus, cache, errors)
- Config from env vars, parsed via `envconfig` into `internal/shared/config/config.go`
- CORS: only allow `FRONTEND_URL` origin
- JSON response keys always camelCase (struct tags)
- Timestamps always ISO 8601 (UTC)
- UUIDs for all primary keys (generated by PostgreSQL `gen_random_uuid()`)
- Pagination: offset-based, default limit=20, max=100
- Soft delete: `deleted_at TIMESTAMPTZ` column, filter with `WHERE deleted_at IS NULL`
- **Cross-context imports are forbidden** — enforced via `depguard` linter rule

## Environment variables

```env
PORT=8080
DATABASE_URL=postgresql://postgres:password@db.xxxxx.supabase.co:5432/postgres
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_ANON_KEY=eyJ...
SUPABASE_SERVICE_ROLE_KEY=eyJ...
NOTION_API_KEY=secret_...                    # Notion integration token
NOTION_DATABASE_ID=...                       # Notion dishes database ID
SYNC_INTERVAL_MINUTES=30                     # Notion → PostgreSQL auto-sync interval
FRONTEND_URL=http://localhost:3000
ENVIRONMENT=development
LOG_LEVEL=debug
```

## Commands (Makefile)

```bash
make dev              # Run with hot reload (air)
make build            # Build binary
make run              # Run built binary
make test             # Unit tests only
make test-integration # Integration tests (needs Docker)
make test-all         # All tests
make test-cover       # All tests + coverage report
make lint             # golangci-lint (strict)
make vet              # go vet
make fmt              # gofmt format
make fmt-check        # gofmt check (CI)
make migrate-up       # Run pending migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration (usage: make migrate-create name=add_tags)
make sqlc             # Regenerate sqlc code
make seed             # Run seed data (categories, tags, nutrition goals, suggestion configs)
make sync             # Trigger Notion → PostgreSQL sync manually
make docker-up        # docker-compose up (local postgres)
make docker-down      # docker-compose down
```

## Dockerfile (Render deployment)

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
COPY migrations ./migrations
EXPOSE 8080
CMD ["./server"]
```

## Render deploy config

- **Service type:** Web Service (Docker)
- **Health check path:** `/health`
- **Auto-deploy:** on push to `main`
- **Free tier note:** sleeps after 15 min inactivity. Set up cron-job.org or UptimeRobot to ping `/health` every 5 minutes.

## Important notes

- All code, comments, variable names, and commit messages in **English**
- All dish content is in **Vietnamese**, managed in **Notion** — never hardcode dish data in code or SQL
- Notion is the **source of truth** for recipes; PostgreSQL is a **read cache** for fast serving
- The Notion parser (`internal/repository/notion/parser.go`) is critical — **unit test every block type**
- Seed data (`seeds/`) is for reference data only: categories, tags, nutrition goals, suggestion configs. Recipe data comes from Notion sync
- The suggestion engine (`internal/service/suggestion/`) is the core differentiator — test thoroughly
- TDEE calculator must be unit tested with known values
- Group strategy balancing logic is complex — test edge cases (1 person, 10 people, empty DB)
- Cross-module reads go through service interfaces, NEVER direct SQL joins across bounded contexts
- Sync runs on startup + every N minutes (configurable) + manual trigger via admin endpoint
- Log all sync operations with slog (added/updated/deleted counts)
- When extracting to microservices: replace EventBus with Kafka, replace in-memory cache with Redis, split DB by moving tables to separate Supabase projects
- Diagrams maintained in `docs/diagrams/` using Mermaid — update whenever architecture or flows change

## Migration path (free → scale)

| Stage | Trigger | Changes |
|-------|---------|---------|
| **Free** | Now | Monolith on Render + Supabase free |
| **Paid** | Real users, cold starts unacceptable | Render Starter ($7), Supabase Pro ($25) |
| **Cache** | Slow repeated queries | Add Redis (Upstash free or Render KV) |
| **Search** | 2000+ recipes, search quality matters | Add Elasticsearch (Bonsai.io free) |
| **Microservices** | Team grows 3+ devs or single bottleneck | Extract User Service → Recipe Service → Suggestion Service |

The by-feature folder structure ensures each extraction is **copy a folder + add infra** — not a rewrite.