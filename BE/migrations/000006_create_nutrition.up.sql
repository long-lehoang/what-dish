CREATE TABLE nutrition_recipe (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id       UUID NOT NULL UNIQUE,
    calories        DECIMAL(8,2),
    protein         DECIMAL(8,2),
    carbs           DECIMAL(8,2),
    fat             DECIMAL(8,2),
    fiber           DECIMAL(8,2),
    sugar           DECIMAL(8,2),
    sodium          DECIMAL(8,2),
    serving_size    VARCHAR(100),
    data_source     VARCHAR(50) DEFAULT 'MANUAL',
    is_verified     BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_nutrition_recipe_calories ON nutrition_recipe(calories);
CREATE INDEX idx_nutrition_recipe_recipe ON nutrition_recipe(recipe_id);

CREATE TABLE nutrition_goals (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(100) NOT NULL,
    description         TEXT,
    meal_calories_min   INT,
    meal_calories_max   INT,
    daily_calories_min  INT,
    daily_calories_max  INT,
    protein_pct         INT,
    carbs_pct           INT,
    fat_pct             INT,
    is_active           BOOLEAN DEFAULT TRUE,
    sort_order          INT DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);
