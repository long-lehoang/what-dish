CREATE TABLE ingredients (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) NOT NULL UNIQUE,
    category        VARCHAR(50),
    default_unit    VARCHAR(20),
    is_allergen     BOOLEAN DEFAULT FALSE,
    allergen_type   VARCHAR(50),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_ingredients_slug ON ingredients(slug);
CREATE INDEX idx_ingredients_category ON ingredients(category);
CREATE INDEX idx_ingredients_allergen ON ingredients(is_allergen) WHERE is_allergen = TRUE;

CREATE TABLE recipe_ingredients (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id       UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id   UUID REFERENCES ingredients(id),
    name            VARCHAR(200) NOT NULL,
    amount          DECIMAL(10,2),
    unit            VARCHAR(20),
    note            VARCHAR(200),
    is_optional     BOOLEAN DEFAULT FALSE,
    group_name      VARCHAR(100),
    sort_order      INT DEFAULT 0
);

CREATE INDEX idx_recipe_ingredients_recipe ON recipe_ingredients(recipe_id);
CREATE INDEX idx_recipe_ingredients_ingredient ON recipe_ingredients(ingredient_id);

CREATE TABLE recipe_steps (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id   UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_number INT NOT NULL,
    title       VARCHAR(200),
    description TEXT NOT NULL,
    image_url   VARCHAR(500),
    duration    INT,
    sort_order  INT DEFAULT 0,
    UNIQUE(recipe_id, step_number)
);

CREATE INDEX idx_recipe_steps_recipe ON recipe_steps(recipe_id);
