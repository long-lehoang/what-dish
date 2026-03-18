# Database Design — WhatDish

## Design Principles

- **Database per Service** — each Bounded Context owns its own database
- **UUID primary keys** — avoid collisions when merging data, suitable for distributed systems
- **Soft delete** — use `deleted_at` instead of hard deletes
- **Audit trail** — every table includes `created_at`, `updated_at`
- **Normalization** — 3NF for source data, denormalize when performance requires it (read models)

---

## 1. Recipe Database (`recipe_db`)

**Owner:** recipe-service

### ERD Overview

```
categories ──1:N──→ recipes ──1:N──→ recipe_ingredients ←──N:1── ingredients
                       │
                       ├──1:N──→ recipe_steps
                       │
                       └──N:M──→ tags (via recipe_tags)
```

### Tables

#### `categories`
Manage dish classification across multiple dimensions.

```sql
CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    type        VARCHAR(50) NOT NULL,
        -- 'DISH_TYPE'       : Braised, Stir-fried, Steamed, Grilled, Fried...
        -- 'REGION'          : Northern, Central, Southern
        -- 'MAIN_INGREDIENT' : Chicken, Beef, Pork, Seafood, Vegetarian...
        -- 'MEAL_TYPE'       : Breakfast, Lunch, Dinner, Snack
    icon_url    VARCHAR(500),
    sort_order  INT DEFAULT 0,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_categories_type ON categories(type);
CREATE INDEX idx_categories_slug ON categories(slug);
```

#### `recipes`
Core table containing dish information. Synced from Notion — `external_id` stores the Notion page ID for upsert matching.

```sql
CREATE TABLE recipes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id     VARCHAR(255) UNIQUE,  -- Notion page ID (for sync matching)
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    image_url       VARCHAR(500),
    prep_time       INT,              -- minutes for preparation
    cook_time       INT,              -- minutes for cooking
    total_time      INT GENERATED ALWAYS AS (prep_time + cook_time) STORED,
    servings        INT DEFAULT 2,    -- default number of servings
    difficulty      VARCHAR(20) DEFAULT 'EASY',
        -- 'EASY', 'MEDIUM', 'HARD'
    status          VARCHAR(20) DEFAULT 'DRAFT',
        -- 'DRAFT', 'PUBLISHED', 'ARCHIVED'

    -- Foreign keys to categories (multi-dimension)
    dish_type_id        UUID REFERENCES categories(id),
    region_id           UUID REFERENCES categories(id),
    main_ingredient_id  UUID REFERENCES categories(id),
    meal_type_id        UUID REFERENCES categories(id),

    -- Metadata
    source_url      VARCHAR(500),     -- Notion page URL
    author_note     TEXT,
    view_count      INT DEFAULT 0,
    favorite_count  INT DEFAULT 0,    -- denormalized counter

    -- Search
    search_vector   tsvector,         -- auto-populated by trigger

    -- Sync tracking
    last_synced_at  TIMESTAMPTZ,      -- last successful sync from Notion

    -- Audit
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ       -- soft delete
);

CREATE INDEX idx_recipes_status ON recipes(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipes_dish_type ON recipes(dish_type_id);
CREATE INDEX idx_recipes_region ON recipes(region_id);
CREATE INDEX idx_recipes_main_ingredient ON recipes(main_ingredient_id);
CREATE INDEX idx_recipes_meal_type ON recipes(meal_type_id);
CREATE INDEX idx_recipes_difficulty ON recipes(difficulty);
CREATE INDEX idx_recipes_cook_time ON recipes(cook_time);
CREATE INDEX idx_recipes_slug ON recipes(slug);
CREATE UNIQUE INDEX idx_recipes_external_id ON recipes(external_id) WHERE external_id IS NOT NULL;
CREATE INDEX idx_recipes_search ON recipes USING GIN(search_vector);
```

#### `ingredients`
Ingredient master data catalog.

```sql
CREATE TABLE ingredients (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) NOT NULL UNIQUE,
    category        VARCHAR(50),
        -- 'MEAT', 'SEAFOOD', 'VEGETABLE', 'SPICE', 'SAUCE', 'GRAIN', 'DAIRY', 'OTHER'
    default_unit    VARCHAR(20),      -- default unit: g, ml, tbsp...
    is_allergen     BOOLEAN DEFAULT FALSE,
    allergen_type   VARCHAR(50),      -- 'GLUTEN', 'DAIRY', 'SHELLFISH', 'NUTS', 'SOY'...
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_ingredients_slug ON ingredients(slug);
CREATE INDEX idx_ingredients_category ON ingredients(category);
CREATE INDEX idx_ingredients_allergen ON ingredients(is_allergen) WHERE is_allergen = TRUE;
```

#### `recipe_ingredients`
Junction table linking recipes to ingredients with measurements.

```sql
CREATE TABLE recipe_ingredients (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id   UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id),
    amount      DECIMAL(10,2),        -- quantity
    unit        VARCHAR(20),          -- g, ml, tbsp, cup, piece, slice...
    note        VARCHAR(200),         -- "finely chopped", "thinly sliced"...
    is_optional BOOLEAN DEFAULT FALSE,
    group_name  VARCHAR(100),         -- "For the broth", "For garnish"
    sort_order  INT DEFAULT 0,

    UNIQUE(recipe_id, ingredient_id)
);

CREATE INDEX idx_recipe_ingredients_recipe ON recipe_ingredients(recipe_id);
CREATE INDEX idx_recipe_ingredients_ingredient ON recipe_ingredients(ingredient_id);
```

#### `recipe_steps`
Cooking instructions.

```sql
CREATE TABLE recipe_steps (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id   UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_number INT NOT NULL,
    title       VARCHAR(200),         -- step title (optional)
    description TEXT NOT NULL,
    image_url   VARCHAR(500),
    duration    INT,                  -- time for this step in minutes (for timer)
    sort_order  INT DEFAULT 0,

    UNIQUE(recipe_id, step_number)
);

CREATE INDEX idx_recipe_steps_recipe ON recipe_steps(recipe_id);
```

#### `tags`
Flexible tags (Weight Loss, Quick, Party, Kid-friendly, Not Spicy...).

```sql
CREATE TABLE tags (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name    VARCHAR(100) NOT NULL UNIQUE,
    slug    VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE recipe_tags (
    recipe_id   UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    tag_id      UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, tag_id)
);

CREATE INDEX idx_recipe_tags_tag ON recipe_tags(tag_id);
```

#### `sync_logs`
Track Notion → PostgreSQL sync history.

```sql
CREATE TABLE sync_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at      TIMESTAMPTZ NOT NULL,
    finished_at     TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'RUNNING',
        -- 'RUNNING', 'SUCCESS', 'PARTIAL', 'FAILED'
    recipes_added   INT DEFAULT 0,
    recipes_updated INT DEFAULT 0,
    recipes_deleted INT DEFAULT 0,
    errors          JSONB DEFAULT '[]',   -- array of error messages
    duration_ms     INT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sync_logs_status ON sync_logs(status, created_at DESC);
```

#### Full-text search trigger

```sql
CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE OR REPLACE FUNCTION update_recipe_search() RETURNS trigger AS $
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('simple', unaccent(COALESCE(NEW.name, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(NEW.description, ''))), 'B');
  RETURN NEW;
END;
$ LANGUAGE plpgsql;

CREATE TRIGGER recipe_search_trigger
  BEFORE INSERT OR UPDATE ON recipes
  FOR EACH ROW EXECUTE FUNCTION update_recipe_search();
```

---

## 2. Suggestion Database (`suggestion_db`)

**Owner:** suggestion-service

### Tables

#### `suggestion_sessions`
Each suggestion request = 1 session.

```sql
CREATE TABLE suggestion_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID,             -- NULL if guest
    session_type    VARCHAR(30) NOT NULL,
        -- 'RANDOM', 'BY_CALORIES', 'BY_GROUP'

    -- Input parameters (flexible JSON)
    input_params    JSONB NOT NULL DEFAULT '{}',
    -- e.g. BY_CALORIES: {"target_calories": 500, "meal_type": "lunch"}
    -- e.g. BY_GROUP:    {"group_size": 4, "group_type": "family_with_kids"}

    -- Output
    result_recipe_ids UUID[] NOT NULL DEFAULT '{}',
    total_calories    INT,

    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_suggestion_sessions_user ON suggestion_sessions(user_id, created_at DESC);
CREATE INDEX idx_suggestion_sessions_type ON suggestion_sessions(session_type);
CREATE INDEX idx_suggestion_sessions_created ON suggestion_sessions(created_at);
```

#### `suggestion_configs`
Preset configurations for suggestion types.

```sql
CREATE TABLE suggestion_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_type     VARCHAR(30) NOT NULL,
        -- 'CALORIE_PRESET', 'GROUP_PRESET'
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    params          JSONB NOT NULL,
    -- CALORIE_PRESET: {"min_cal": 300, "max_cal": 400, "label": "Weight loss"}
    -- GROUP_PRESET:   {"min_dishes": 3, "max_dishes": 5, "constraints": {...}}
    is_active       BOOLEAN DEFAULT TRUE,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### `exclusion_rules`
Store recently suggested recipe IDs to avoid repetition.

```sql
CREATE TABLE exclusion_rules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    recipe_id   UUID NOT NULL,
    excluded_until TIMESTAMPTZ NOT NULL,  -- default: created_at + 7 days
    created_at  TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, recipe_id)
);

CREATE INDEX idx_exclusion_rules_user ON exclusion_rules(user_id, excluded_until);
```

---

## 3. User Database (`user_db`)

**Owner:** user-service

### Tables

#### `users`

```sql
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255),     -- NULL if OAuth-only
    name            VARCHAR(200) NOT NULL,
    avatar_url      VARCHAR(500),
    role            VARCHAR(20) DEFAULT 'USER',
        -- 'USER', 'ADMIN', 'SUPER_ADMIN'
    status          VARCHAR(20) DEFAULT 'PENDING_VERIFICATION',
        -- 'PENDING_VERIFICATION', 'ACTIVE', 'SUSPENDED', 'DELETED'
    email_verified  BOOLEAN DEFAULT FALSE,

    -- OAuth
    oauth_provider  VARCHAR(20),      -- 'GOOGLE', 'FACEBOOK'
    oauth_id        VARCHAR(255),

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_id);
```

#### `user_profiles`
Nutrition profile.

```sql
CREATE TABLE user_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    gender          VARCHAR(10),      -- 'MALE', 'FEMALE'
    age             INT,
    height_cm       DECIMAL(5,1),
    weight_kg       DECIMAL(5,1),
    activity_level  VARCHAR(30),
        -- 'SEDENTARY', 'LIGHT', 'MODERATE', 'ACTIVE', 'VERY_ACTIVE'
    goal            VARCHAR(20),
        -- 'LOSE_WEIGHT', 'MAINTAIN', 'GAIN_WEIGHT'

    -- Calculated
    bmr             DECIMAL(8,2),     -- Basal Metabolic Rate
    tdee            DECIMAL(8,2),     -- Total Daily Energy Expenditure
    daily_target    DECIMAL(8,2),     -- TDEE adjusted by goal

    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### `user_allergies`
Ingredients the user is allergic to or dislikes.

```sql
CREATE TABLE user_allergies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ingredient_id   UUID NOT NULL,    -- references recipe_db.ingredients
    allergy_type    VARCHAR(20) DEFAULT 'ALLERGY',
        -- 'ALLERGY' (actual allergy), 'DISLIKE' (preference)

    UNIQUE(user_id, ingredient_id)
);

CREATE INDEX idx_user_allergies_user ON user_allergies(user_id);
```

#### `refresh_tokens`

```sql
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
```

---

## 4. Nutrition Database (`nutrition_db`)

**Owner:** nutrition-service

### Tables

#### `recipe_nutrition`
Nutrition information per serving.

```sql
CREATE TABLE recipe_nutrition (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id   UUID NOT NULL UNIQUE,  -- references recipe_db.recipes

    -- Macros (per serving)
    calories    DECIMAL(8,2),         -- kcal
    protein     DECIMAL(8,2),         -- grams
    carbs       DECIMAL(8,2),         -- grams
    fat         DECIMAL(8,2),         -- grams
    fiber       DECIMAL(8,2),         -- grams
    sugar       DECIMAL(8,2),         -- grams
    sodium      DECIMAL(8,2),         -- mg

    -- Meta
    serving_size    VARCHAR(100),     -- "1 bowl", "1 serving", "1 plate"
    data_source     VARCHAR(50),      -- 'MANUAL', 'CALCULATED', 'IMPORTED'
    is_verified     BOOLEAN DEFAULT FALSE,

    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_recipe_nutrition_calories ON recipe_nutrition(calories);
CREATE INDEX idx_recipe_nutrition_recipe ON recipe_nutrition(recipe_id);
```

#### `nutrition_goals`
Preset nutrition goals.

```sql
CREATE TABLE nutrition_goals (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,    -- "Weight Loss", "Muscle Gain", "Eat Clean"
    description TEXT,

    -- Target ranges per meal
    meal_calories_min   INT,
    meal_calories_max   INT,
    daily_calories_min  INT,
    daily_calories_max  INT,

    -- Macro ratios (percentage)
    protein_pct     INT,              -- 30
    carbs_pct       INT,              -- 40
    fat_pct         INT,              -- 30

    is_active       BOOLEAN DEFAULT TRUE,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 5. Engagement Database (`engagement_db`)

**Owner:** engagement-service

### Tables

#### `favorites`

```sql
CREATE TABLE favorites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    recipe_id   UUID NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, recipe_id)
);

CREATE INDEX idx_favorites_user ON favorites(user_id, created_at DESC);
CREATE INDEX idx_favorites_recipe ON favorites(recipe_id);
```

#### `view_history`

```sql
CREATE TABLE view_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID,                 -- NULL if guest (track by session)
    session_id  VARCHAR(100),         -- anonymous session tracking
    recipe_id   UUID NOT NULL,
    source      VARCHAR(30),          -- 'RANDOM', 'SEARCH', 'FILTER', 'DIRECT', 'SUGGESTION'
    viewed_at   TIMESTAMPTZ DEFAULT NOW()
);

-- Consider partitioning by month for performance
CREATE INDEX idx_view_history_user ON view_history(user_id, viewed_at DESC);
CREATE INDEX idx_view_history_recipe ON view_history(recipe_id);
CREATE INDEX idx_view_history_date ON view_history(viewed_at);
```

#### `ratings`

```sql
CREATE TABLE ratings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    recipe_id   UUID NOT NULL,
    score       SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment     TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, recipe_id)
);

CREATE INDEX idx_ratings_recipe ON ratings(recipe_id);
CREATE INDEX idx_ratings_user ON ratings(user_id);
```

---

## 6. Cross-Service Data Consistency

### The Problem

Each service has its own DB, so there are no cross-database foreign keys. For example, `suggestion_sessions.result_recipe_ids` contains UUIDs from `recipe_db.recipes`, but there is no FK constraint.

### Solutions

**Eventual Consistency via Events:**
- When a recipe is deleted, Recipe Service publishes `RecipeDeleted`
- Suggestion Service listens → removes recipe_id from exclusion_rules
- Engagement Service listens → soft-deletes related favorites, views
- Nutrition Service listens → archives nutrition data

**Application-Layer Validation:**
- Suggestion Service calls Recipe Service to verify recipe_id before returning results
- If recipe no longer exists → skip it, random another dish

**Idempotent Event Handlers:**
- Each event has a unique `event_id`
- Consumers store `processed_events` table to avoid duplicate processing

```sql
-- Each service has this table
CREATE TABLE processed_events (
    event_id    UUID PRIMARY KEY,
    event_type  VARCHAR(100),
    processed_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 7. Migration Strategy

### Recommended Tools

- **Node.js:** `knex` migrations or `prisma migrate`
- **Java:** Flyway or Liquibase

### Naming Convention

```
migrations/
├── V001__create_recipes_table.sql
├── V002__create_ingredients_table.sql
├── V003__create_recipe_ingredients_table.sql
├── V004__create_recipe_steps_table.sql
├── V005__add_tags_system.sql
└── V006__add_recipe_indexes.sql
```

### Seed Data

Recipe data comes from Notion sync — seed files are for reference/config data only.

```
seeds/
├── 01_categories.sql       -- Dish types, regions, main ingredients
├── 02_tags.sql              -- Common tags
├── 03_nutrition_goals.sql   -- Preset goals (weight loss, muscle gain, etc.)
└── 04_suggestion_configs.sql -- Suggestion presets
```

Note: Categories and tags must match the select options in the Notion database.
When adding new categories, also add them in Notion.