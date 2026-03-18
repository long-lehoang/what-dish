# System Architecture — WhatDish

## 1. Architecture Overview

The system is designed following **Microservices Architecture** with the principles:

- **Database per Service** — each service owns its own database, no shared DB
- **Bounded Context** — clear business boundaries between services
- **Event-Driven** — asynchronous communication via Message Bus for cross-service concerns
- **API Gateway** — single entry point for all client requests

```
                          ┌──────────────┐
                          │    Client     │
                          │ (Web/Mobile)  │
                          └──────┬───────┘
                                 │
                          ┌──────▼───────┐
                          │  API Gateway  │
                          │ (Kong/Nginx)  │
                          └──────┬───────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
        │              ┌─────────▼──────────┐             │
        │              │    Message Bus      │             │
        │              │  (Kafka/RabbitMQ)   │             │
        │              └─┬───┬───┬───┬───┬──┘             │
        │                │   │   │   │   │                │
   ┌────▼────┐    ┌──────▼┐ ┌▼───▼┐ ┌▼───▼──┐    ┌───────▼──┐
   │ Recipe  │    │Suggest│ │User │ │Nutrit- │    │Engagement│
   │ Service │    │ ion   │ │Svc  │ │  ion   │    │ Service  │
   └────┬────┘    │Service│ └──┬──┘ │Service │    └────┬─────┘
        │         └───┬───┘    │    └───┬────┘         │
   ┌────▼────┐   ┌────▼───┐┌──▼──┐ ┌───▼────┐   ┌────▼─────┐
   │Recipe DB│   │Suggest ││User │ │Nutrit- │   │Engage-   │
   │(Postgre)│   │  DB    ││ DB  │ │ion DB  │   │ment DB   │
   └─────────┘   └────────┘└─────┘ └────────┘   └──────────┘

   ┌──────────┐   ┌───────┐
   │Elastic-  │   │ Redis │
   │ search   │   │ Cache │
   └──────────┘   └───────┘
```

---

## 2. Bounded Contexts

### 2.1 Recipe Context (Primary Domain)

**Responsibility:** Manage all data related to dishes, recipes, ingredients, and categories. Recipe content is sourced from **Notion** (free CMS) and cached in PostgreSQL for fast reads.

**Service:** `recipe-service`  
**Database:** `recipe_db` (PostgreSQL)  
**Port:** 3001

**Content Flow:** Notion (source of truth) → SyncService → PostgreSQL (read cache) → API

**Aggregates:**
- Recipe (root) → RecipeIngredient, RecipeStep
- Category
- Ingredient
- Tag
- SyncLog

**API Endpoints:**
- `GET /recipes` — list (paginated, filterable)
- `GET /recipes/:id` — details
- `GET /recipes/random` — random 1 dish
- `GET /recipes/search` — full-text search
- `GET /categories` — categories
- `POST /admin/sync` — trigger Notion → PostgreSQL sync (admin only)

Note: No POST/PUT/DELETE for recipes — content is managed in Notion.

**Events Published:**
- `RecipeSynced` — when sync completes (with added/updated/deleted counts)

**Events Consumed:**
- `CategoryUpdated` — update reference data

---

### 2.2 Suggestion Context (Core Business Logic)

**Responsibility:** Handle all suggestion logic: smart random, calorie-based, and group-based suggestions.

**Service:** `suggestion-service`  
**Database:** `suggestion_db` (PostgreSQL)  
**Port:** 3002

**Aggregates:**
- SuggestionSession (root) → SuggestionItem
- SuggestionConfig (presets for group and calorie modes)

**Core Logic:**

```
Smart Random:
├─ Input: user_id (optional), filters (optional)
├─ Process:
│   ├─ Fetch recipe IDs from Recipe Service (cached in Redis)
│   ├─ Exclude: 7-day history (if user exists)
│   ├─ Apply filters (dish type, region, etc.)
│   ├─ Random select from remaining set
│   └─ Enrich: call Recipe Service for full details
└─ Output: Recipe detail + nutrition summary

Calorie-Based Suggestion:
├─ Input: target_calories, meal_type, filters
├─ Process:
│   ├─ Fetch recipe IDs + calories from Nutrition Service
│   ├─ Filter: calories within target ± 15%
│   ├─ If meal_type = "full_day": distribute 30/40/30
│   ├─ Random select from filtered set
│   └─ Calculate total nutrition for combo
└─ Output: List<Recipe> + total nutrition

Group Meal Suggestion:
├─ Input: group_size, group_type, meal_type, budget
├─ Process:
│   ├─ Determine dish count: ceil(group_size / 2) + 1, min=3, max=5
│   ├─ Slot allocation:
│   │   ├─ 1 main dish (protein-heavy)
│   │   ├─ 1 soup
│   │   └─ N side dishes (vegetables, stir-fried, braised...)
│   ├─ Constraint solving:
│   │   ├─ No duplicate main ingredients
│   │   ├─ Diverse cooking_method
│   │   ├─ Balanced flavor_profile
│   │   └─ Same region or "fusion" if Custom selected
│   ├─ Scale servings to group_size
│   └─ Merge ingredients → shopping list
└─ Output: Combo + merged shopping list + total nutrition
```

**API Endpoints:**
- `POST /suggestions/random` — random 1 dish
- `POST /suggestions/by-calories` — calorie-based suggestion
- `POST /suggestions/by-group` — group meal suggestion
- `GET /suggestions/history` — history (requires auth)

**Events Published:**
- `SuggestionCreated` — on each successful suggestion

**Events Consumed:**
- `RecipeCreated`, `RecipeUpdated`, `RecipeDeleted` — update cached recipe list
- `UserProfileUpdated` — fetch updated TDEE for user

---

### 2.3 User Context (Identity & Access Management)

**Responsibility:** Authentication, authorization, user profile management.

**Service:** `user-service`  
**Database:** `user_db` (PostgreSQL)  
**Port:** 3003

**Aggregates:**
- User (root) → UserProfile, UserAllergy
- RefreshToken

**API Endpoints:**
- `POST /auth/register` — registration
- `POST /auth/login` — login
- `POST /auth/refresh` — refresh token
- `POST /auth/logout` — logout
- `POST /auth/forgot-password` — forgot password
- `POST /auth/reset-password` — reset password
- `GET /auth/oauth/:provider` — OAuth redirect
- `GET /users/me` — current user info
- `PUT /users/me/profile` — update nutrition profile
- `GET /users` — user list (admin)
- `PUT /users/:id/status` — suspend/reactivate (admin)

**Events Published:**
- `UserCreated` — on successful registration
- `UserProfileUpdated` — on nutrition profile update
- `UserSuspended` — admin suspends account
- `UserDeleted` — account deleted

---

### 2.4 Nutrition Context

**Responsibility:** Manage nutrition data, calculate TDEE, provide nutrition data to other services.

**Service:** `nutrition-service`  
**Database:** `nutrition_db` (PostgreSQL)  
**Port:** 3004

**Aggregates:**
- RecipeNutrition
- NutritionGoal (presets)
- TDEECalculation

**Core Logic — TDEE Calculation:**

```
Mifflin-St Jeor Formula:
  Male:   BMR = 10 × weight(kg) + 6.25 × height(cm) - 5 × age + 5
  Female: BMR = 10 × weight(kg) + 6.25 × height(cm) - 5 × age - 161

Activity Multiplier:
  Sedentary:       1.2
  Light exercise:  1.375
  Moderate:        1.55
  Active:          1.725
  Very active:     1.9

TDEE = BMR × Activity Multiplier

Goal Adjustment:
  Lose weight:  TDEE × 0.80
  Maintain:     TDEE × 1.00
  Gain weight:  TDEE × 1.15
```

**API Endpoints:**
- `GET /nutrition/recipes/:id` — nutrition for 1 recipe
- `GET /nutrition/recipes?ids=1,2,3` — batch query
- `POST /nutrition/calculate-tdee` — calculate TDEE
- `GET /nutrition/recipes/by-calories?min=300&max=500` — find recipe IDs by calorie range
- `POST /nutrition/recipes/:id` — add/update nutrition data (admin)

**Events Consumed:**
- `RecipeCreated` — initialize empty nutrition record
- `UserProfileUpdated` — recalculate TDEE

---

### 2.5 Engagement Context

**Responsibility:** Manage user interactions: favorites, ratings, view history.

**Service:** `engagement-service`  
**Database:** `engagement_db` (PostgreSQL)  
**Port:** 3005

**Aggregates:**
- Favorite
- ViewHistory
- Rating

**API Endpoints:**
- `POST /favorites` — add favorite
- `DELETE /favorites/:recipe_id` — remove favorite
- `GET /favorites` — favorites list (requires auth)
- `GET /favorites/check?recipe_id=X` — check if favorited
- `POST /views` — record view
- `GET /stats/recipes/:id` — views + favorites count

**Events Published:**
- `RecipeFavorited` — added to favorites
- `RecipeUnfavorited` — removed from favorites
- `RecipeViewed` — recipe viewed

**Events Consumed:**
- `UserDeleted` — cleanup favorites, views, ratings for that user

---

## 3. Infrastructure Components

### 3.1 API Gateway

**Technology:** Kong / Nginx  
**Responsibilities:**
- **Routing:** dispatch requests to the correct service
- **Authentication:** verify JWT token, inject user info into headers
- **Rate Limiting:** abuse prevention (100 req/min/IP for guests, 300 for users)
- **CORS:** cross-origin management
- **Load Balancing:** round-robin across instances

**Routing Table:**

```
/api/v1/recipes/*         → recipe-service:3001
/api/v1/suggestions/*     → suggestion-service:3002
/api/v1/auth/*            → user-service:3003
/api/v1/users/*           → user-service:3003
/api/v1/nutrition/*       → nutrition-service:3004
/api/v1/favorites/*       → engagement-service:3005
/api/v1/views/*           → engagement-service:3005
/api/v1/search/*          → recipe-service:3001 (proxy to Elasticsearch)
```

### 3.2 Message Bus (Kafka / RabbitMQ)

**Topics/Queues:**

```
events.recipe.synced
events.user.created
events.user.profile-updated
events.user.suspended
events.user.deleted
events.suggestion.created
events.engagement.favorited
events.engagement.unfavorited
events.engagement.viewed
```

**Event Schema (example):**

```json
{
  "event_id": "uuid",
  "event_type": "RecipeCreated",
  "timestamp": "2025-01-15T10:30:00Z",
  "source": "recipe-service",
  "payload": {
    "recipe_id": "uuid",
    "name": "Pho Bo",
    "category_id": "uuid"
  }
}
```

### 3.3 Redis Cache

**Usage:**
- **Recipe list cache:** recipe IDs + basic info (TTL: 5 min)
- **Session cache:** JWT blacklist for logout
- **Rate limiting counter:** request count per IP
- **Popular recipes:** top 50 by views (TTL: 1 hour)
- **User TDEE cache:** calculated TDEE (TTL: 24 hours, invalidate on profile update)

### 3.4 Elasticsearch

**Index:** `recipes`

**Mapping:**

```json
{
  "mappings": {
    "properties": {
      "name": { "type": "text", "analyzer": "vietnamese" },
      "description": { "type": "text", "analyzer": "vietnamese" },
      "ingredients": { "type": "text", "analyzer": "vietnamese" },
      "tags": { "type": "keyword" },
      "region": { "type": "keyword" },
      "category": { "type": "keyword" },
      "calories": { "type": "integer" },
      "cook_time": { "type": "integer" },
      "difficulty": { "type": "keyword" },
      "popularity": { "type": "float" },
      "created_at": { "type": "date" }
    }
  }
}
```

**Sync:** Recipe Service publishes `RecipeCreated/Updated/Deleted` → Search consumer updates the index.

---

## 4. Inter-Service Communication

### 4.1 Synchronous (REST/gRPC)

Used for: real-time request-response where results are needed immediately.

```
Client → API Gateway → Suggestion Service
                            │
                            ├── GET Recipe Service: fetch recipe details
                            ├── GET Nutrition Service: fetch calories
                            └── GET Engagement Service: check favorites
```

**Circuit Breaker Pattern:** If one service goes down, other services avoid cascade failure. Libraries: `opossum` (Node.js) or `resilience4j` (Java).

**Timeout:** 3 seconds for inter-service calls. On timeout → return partial data + warning.

### 4.2 Asynchronous (Event-Driven)

Used for: side effects, data sync, analytics — no immediate response needed.

```
User Service ──UserCreated──→ Message Bus
                                   │
                      ┌────────────┼────────────┐
                      ▼            ▼            ▼
              Suggestion Svc  Nutrition Svc  Email Svc
              (init prefs)    (init TDEE)    (welcome)
```

### 4.3 Communication Matrix

| Source | Target | Method | Purpose |
|--------|--------|--------|---------|
| Suggestion → Recipe | Sync (REST) | Fetch recipe details for suggestions |
| Suggestion → Nutrition | Sync (REST) | Fetch calories for filtering |
| Suggestion → Engagement | Sync (REST) | Check history, favorites |
| Notion → Recipe (SyncService) | Async (Cron) | Sync recipe content on schedule |
| Recipe → Search Index | Async (DB trigger) | Update tsvector on recipe upsert |
| User → Nutrition | Async (Event) | Calculate TDEE on profile update |
| User → Engagement | Async (Event) | Cleanup on user deletion |
| Engagement → Suggestion | Async (Event) | Update preference model |

---

## 5. Deployment Architecture

### 5.1 Development (Docker Compose)

```yaml
# docker-compose.yml
services:
  api-gateway:
    image: kong:3.4
    ports: ["8000:8000", "8001:8001"]

  recipe-service:
    build: ./services/recipe
    ports: ["3001:3001"]
    depends_on: [recipe-db, kafka, redis]

  suggestion-service:
    build: ./services/suggestion
    ports: ["3002:3002"]
    depends_on: [suggestion-db, kafka, redis]

  user-service:
    build: ./services/user
    ports: ["3003:3003"]
    depends_on: [user-db, kafka, redis]

  nutrition-service:
    build: ./services/nutrition
    ports: ["3004:3004"]
    depends_on: [nutrition-db, kafka]

  engagement-service:
    build: ./services/engagement
    ports: ["3005:3005"]
    depends_on: [engagement-db, kafka]

  recipe-db:
    image: postgres:16
    environment:
      POSTGRES_DB: recipe_db

  suggestion-db:
    image: postgres:16
    environment:
      POSTGRES_DB: suggestion_db

  user-db:
    image: postgres:16
    environment:
      POSTGRES_DB: user_db

  nutrition-db:
    image: postgres:16
    environment:
      POSTGRES_DB: nutrition_db

  engagement-db:
    image: postgres:16
    environment:
      POSTGRES_DB: engagement_db

  kafka:
    image: confluentinc/cp-kafka:7.5.0

  redis:
    image: redis:7-alpine

  elasticsearch:
    image: elasticsearch:8.11.0
```

### 5.2 Production (Kubernetes)

```
┌────────────────────────────────────────┐
│            Kubernetes Cluster          │
│                                        │
│  ┌──────────┐  ┌──────────────────┐   │
│  │ Ingress  │  │  Service Mesh    │   │
│  │ (Nginx)  │  │  (Istio/Linkerd) │   │
│  └────┬─────┘  └────────┬─────────┘   │
│       │                  │             │
│  ┌────▼──────────────────▼──────────┐  │
│  │        Namespace: whatdish        │  │
│  │                                  │  │
│  │  recipe-svc (2 replicas)         │  │
│  │  suggestion-svc (2 replicas)     │  │
│  │  user-svc (2 replicas)           │  │
│  │  nutrition-svc (1 replica)       │  │
│  │  engagement-svc (1 replica)      │  │
│  └──────────────────────────────────┘  │
│                                        │
│  ┌──────────────────────────────────┐  │
│  │     Namespace: data              │  │
│  │  PostgreSQL (managed / RDS)      │  │
│  │  Redis (managed / ElastiCache)   │  │
│  │  Kafka (managed / MSK)           │  │
│  │  Elasticsearch (managed)         │  │
│  └──────────────────────────────────┘  │
└────────────────────────────────────────┘
```

---

## 6. Scaling Strategy

### Phase 1 — Monolith First (MVP)

For small teams (1-3 devs), start with a **Modular Monolith organized by feature (bounded context)**: everything in one codebase, each context in its own folder with model, port, service, handler, repository, and DTO. Use a shared database with table-prefix separation.

```
internal/
├── recipe/                 # Recipe BC (self-contained)
│   ├── model.go            # Domain entities
│   ├── port.go             # Interface definitions
│   ├── service.go          # Business logic
│   ├── handler.go          # HTTP handlers
│   ├── repository.go       # PostgreSQL implementation
│   ├── dto.go              # Request/response DTOs
│   └── queries/            # SQL for sqlc
├── suggestion/             # Suggestion BC (self-contained)
├── user/                   # User BC (self-contained)
├── nutrition/              # Nutrition BC (self-contained)
├── engagement/             # Engagement BC (self-contained)
├── shared/                 # Cross-cutting: middleware, events, cache, errors
└── platform/               # External adapters: notion/, supabase/
```

**Key rule:** Bounded contexts never import each other directly. Cross-context communication uses interfaces (defined in each context's `port.go`) and is wired via dependency injection in `cmd/server/main.go`.

### Phase 2 — Extract Services (when scaling is needed)

As traffic grows or the team expands, extract context folders into separate services. Since each folder is self-contained (model + port + service + handler + repo + queries), extraction = copy folder + add infra. Priority order:

1. **User Service** — extract first (auth is a cross-cutting concern)
2. **Recipe Service** — read-heavy, needs independent scaling
3. **Suggestion Service** — compute-heavy, needs independent scaling
4. **Nutrition + Engagement** — extract last

### Phase 3 — Full Microservices

Each service runs independently, deploys separately, scales independently. Full adoption of: API Gateway, Message Bus, Service Mesh, Observability (logging, metrics, tracing).