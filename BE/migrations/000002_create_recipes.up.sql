CREATE TABLE recipes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id         VARCHAR(255) UNIQUE,
    name                VARCHAR(255) NOT NULL,
    slug                VARCHAR(255) NOT NULL UNIQUE,
    description         TEXT,
    image_url           VARCHAR(500),
    prep_time           INT,
    cook_time           INT,
    total_time          INT GENERATED ALWAYS AS (COALESCE(prep_time, 0) + COALESCE(cook_time, 0)) STORED,
    servings            INT DEFAULT 2,
    difficulty          VARCHAR(20) DEFAULT 'EASY',
    status              VARCHAR(20) DEFAULT 'DRAFT',
    dish_type_id        UUID REFERENCES categories(id),
    region_id           UUID REFERENCES categories(id),
    main_ingredient_id  UUID REFERENCES categories(id),
    meal_type_id        UUID REFERENCES categories(id),
    source_url          VARCHAR(500),
    author_note         TEXT,
    view_count          INT DEFAULT 0,
    favorite_count      INT DEFAULT 0,
    search_vector       tsvector,
    last_synced_at      TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
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
