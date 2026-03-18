CREATE TABLE user_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE,
    gender          VARCHAR(10),
    age             INT,
    height_cm       DECIMAL(5,1),
    weight_kg       DECIMAL(5,1),
    activity_level  VARCHAR(30),
    goal            VARCHAR(20),
    bmr             DECIMAL(8,2),
    tdee            DECIMAL(8,2),
    daily_target    DECIMAL(8,2),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_allergies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    ingredient_id   UUID,
    ingredient_name VARCHAR(200) NOT NULL,
    allergy_type    VARCHAR(20) DEFAULT 'ALLERGY',
    UNIQUE(user_id, ingredient_name)
);

CREATE INDEX idx_user_allergies_user ON user_allergies(user_id);
