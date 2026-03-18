CREATE TABLE engagement_favorites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    recipe_id   UUID NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, recipe_id)
);

CREATE INDEX idx_engagement_favorites_user ON engagement_favorites(user_id, created_at DESC);
CREATE INDEX idx_engagement_favorites_recipe ON engagement_favorites(recipe_id);

CREATE TABLE engagement_views (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID,
    session_id  VARCHAR(100),
    recipe_id   UUID NOT NULL,
    source      VARCHAR(30),
    viewed_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_engagement_views_user ON engagement_views(user_id, viewed_at DESC);
CREATE INDEX idx_engagement_views_recipe ON engagement_views(recipe_id);
CREATE INDEX idx_engagement_views_date ON engagement_views(viewed_at);

CREATE TABLE engagement_ratings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    recipe_id   UUID NOT NULL,
    score       SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment     TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, recipe_id)
);

CREATE INDEX idx_engagement_ratings_recipe ON engagement_ratings(recipe_id);
CREATE INDEX idx_engagement_ratings_user ON engagement_ratings(user_id);
