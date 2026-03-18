package recipe

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrors "github.com/lehoanglong/whatdish/internal/shared/errors"
)

// ============================================================
// DishRepository implementation
// ============================================================

type pgDishRepo struct {
	pool *pgxpool.Pool
}

func NewDishRepository(pool *pgxpool.Pool) DishRepository {
	return &pgDishRepo{pool: pool}
}

func (r *pgDishRepo) List(ctx context.Context, filter DishFilter) ([]Dish, int64, error) {
	where, args := buildDishFilterClause(filter)

	// Count query.
	countSQL := "SELECT COUNT(*) FROM recipes r " + where
	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("recipe.List count: %w", err)
	}

	if total == 0 {
		return []Dish{}, 0, nil
	}

	// Data query with pagination.
	offset := (filter.Page - 1) * filter.PageSize
	dataSQL := `
		SELECT r.id, r.external_id, r.name, r.slug, r.description, r.image_url,
		       r.prep_time, r.cook_time, r.total_time, r.servings, r.difficulty,
		       r.status, r.dish_type_id, r.region_id, r.main_ingredient_id, r.meal_type_id,
		       r.source_url, r.author_note, r.view_count, r.favorite_count,
		       r.last_synced_at, r.created_at, r.updated_at, r.deleted_at
		FROM recipes r
		` + where + fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT %d OFFSET %d", filter.PageSize, offset)

	rows, err := r.pool.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("recipe.List query: %w", err)
	}
	defer rows.Close()

	var dishes []Dish
	for rows.Next() {
		d, err := scanDish(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("recipe.List scan: %w", err)
		}
		dishes = append(dishes, d)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("recipe.List rows: %w", err)
	}

	return dishes, total, nil
}

func (r *pgDishRepo) GetByID(ctx context.Context, id uuid.UUID) (*DishDetail, error) {
	return r.getDishDetail(ctx, "r.id = $1", id)
}

func (r *pgDishRepo) GetBySlug(ctx context.Context, slug string) (*DishDetail, error) {
	return r.getDishDetail(ctx, "r.slug = $1", slug)
}

func (r *pgDishRepo) getDishDetail(ctx context.Context, condition string, arg any) (*DishDetail, error) {
	dishSQL := `
		SELECT r.id, r.external_id, r.name, r.slug, r.description, r.image_url,
		       r.prep_time, r.cook_time, r.total_time, r.servings, r.difficulty,
		       r.status, r.dish_type_id, r.region_id, r.main_ingredient_id, r.meal_type_id,
		       r.source_url, r.author_note, r.view_count, r.favorite_count,
		       r.last_synced_at, r.created_at, r.updated_at, r.deleted_at
		FROM recipes r
		WHERE ` + condition + ` AND r.deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, dishSQL, arg)
	d, err := scanDishRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("recipe.getDishDetail dish: %w", err)
	}

	detail := &DishDetail{Dish: d}

	// Load ingredients.
	ingSQL := `
		SELECT id, recipe_id, ingredient_id, name, amount, unit, note,
		       is_optional, group_name, sort_order
		FROM recipe_ingredients
		WHERE recipe_id = $1
		ORDER BY sort_order`
	ingRows, err := r.pool.Query(ctx, ingSQL, d.ID)
	if err != nil {
		return nil, fmt.Errorf("recipe.getDishDetail ingredients: %w", err)
	}
	defer ingRows.Close()

	for ingRows.Next() {
		var ing Ingredient
		if err := ingRows.Scan(
			&ing.ID, &ing.RecipeID, &ing.IngredientID, &ing.Name,
			&ing.Amount, &ing.Unit, &ing.Note, &ing.IsOptional,
			&ing.GroupName, &ing.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("recipe.getDishDetail scan ingredient: %w", err)
		}
		detail.Ingredients = append(detail.Ingredients, ing)
	}
	if err := ingRows.Err(); err != nil {
		return nil, fmt.Errorf("recipe.getDishDetail ingredient rows: %w", err)
	}
	if detail.Ingredients == nil {
		detail.Ingredients = []Ingredient{}
	}

	// Load steps.
	stepSQL := `
		SELECT id, recipe_id, step_number, title, description, image_url,
		       duration, sort_order
		FROM recipe_steps
		WHERE recipe_id = $1
		ORDER BY sort_order`
	stepRows, err := r.pool.Query(ctx, stepSQL, d.ID)
	if err != nil {
		return nil, fmt.Errorf("recipe.getDishDetail steps: %w", err)
	}
	defer stepRows.Close()

	for stepRows.Next() {
		var s Step
		if err := stepRows.Scan(
			&s.ID, &s.RecipeID, &s.StepNumber, &s.Title,
			&s.Description, &s.ImageURL, &s.Duration, &s.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("recipe.getDishDetail scan step: %w", err)
		}
		detail.Steps = append(detail.Steps, s)
	}
	if err := stepRows.Err(); err != nil {
		return nil, fmt.Errorf("recipe.getDishDetail step rows: %w", err)
	}
	if detail.Steps == nil {
		detail.Steps = []Step{}
	}

	// Load tags.
	tagSQL := `
		SELECT t.id, t.name, t.slug
		FROM tags t
		INNER JOIN recipe_tags rt ON rt.tag_id = t.id
		WHERE rt.recipe_id = $1
		ORDER BY t.name`
	tagRows, err := r.pool.Query(ctx, tagSQL, d.ID)
	if err != nil {
		return nil, fmt.Errorf("recipe.getDishDetail tags: %w", err)
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var t Tag
		if err := tagRows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, fmt.Errorf("recipe.getDishDetail scan tag: %w", err)
		}
		detail.Tags = append(detail.Tags, t)
	}
	if err := tagRows.Err(); err != nil {
		return nil, fmt.Errorf("recipe.getDishDetail tag rows: %w", err)
	}
	if detail.Tags == nil {
		detail.Tags = []Tag{}
	}

	// Load categories in a single batch query.
	catIDs := make([]uuid.UUID, 0, 4)
	idSlots := []*uuid.UUID{d.DishTypeID, d.RegionID, d.MainIngredientID, d.MealTypeID}
	for _, id := range idSlots {
		if id != nil {
			catIDs = append(catIDs, *id)
		}
	}
	catMap := r.loadCategories(ctx, catIDs)
	if d.DishTypeID != nil {
		if c, ok := catMap[*d.DishTypeID]; ok {
			detail.DishType = &c
		}
	}
	if d.RegionID != nil {
		if c, ok := catMap[*d.RegionID]; ok {
			detail.Region = &c
		}
	}
	if d.MainIngredientID != nil {
		if c, ok := catMap[*d.MainIngredientID]; ok {
			detail.MainIngredient = &c
		}
	}
	if d.MealTypeID != nil {
		if c, ok := catMap[*d.MealTypeID]; ok {
			detail.MealType = &c
		}
	}

	return detail, nil
}

// loadCategories fetches multiple categories in a single query and returns them as a map.
func (r *pgDishRepo) loadCategories(ctx context.Context, ids []uuid.UUID) map[uuid.UUID]Category {
	result := make(map[uuid.UUID]Category, len(ids))
	if len(ids) == 0 {
		return result
	}

	catSQL := `
		SELECT id, name, slug, type, icon_url, sort_order, is_active, created_at, updated_at
		FROM categories
		WHERE id = ANY($1)`
	rows, err := r.pool.Query(ctx, catSQL, ids)
	if err != nil {
		slog.Warn("recipe.loadCategories query failed", "error", err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var c Category
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.Type, &c.IconURL,
			&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			slog.Warn("recipe.loadCategories scan failed", "error", err)
			continue
		}
		result[c.ID] = c
	}
	return result
}

func (r *pgDishRepo) GetRandom(ctx context.Context, filter DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error) {
	where, args := buildDishFilterClause(filter)

	if len(excludeIDs) > 0 {
		placeholders := make([]string, len(excludeIDs))
		for i, eid := range excludeIDs {
			args = append(args, eid)
			placeholders[i] = fmt.Sprintf("$%d", len(args))
		}
		where += " AND r.id NOT IN (" + strings.Join(placeholders, ",") + ")"
	}

	randomSQL := `
		SELECT r.id FROM recipes r
		` + where + `
		ORDER BY RANDOM() LIMIT 1`

	var dishID uuid.UUID
	if err := r.pool.QueryRow(ctx, randomSQL, args...).Scan(&dishID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("recipe.GetRandom: %w", err)
	}

	return r.GetByID(ctx, dishID)
}

func (r *pgDishRepo) Search(ctx context.Context, query string, page, pageSize int) ([]Dish, int64, error) {
	offset := (page - 1) * pageSize

	// Convert the user query into a tsquery-compatible form.
	// Split words and join with '&' for AND semantics.
	words := strings.Fields(query)
	tsTerms := make([]string, len(words))
	for i, w := range words {
		tsTerms[i] = w + ":*"
	}
	tsQuery := strings.Join(tsTerms, " & ")

	countSQL := `
		SELECT COUNT(*)
		FROM recipes r
		WHERE r.deleted_at IS NULL
		  AND r.status = 'PUBLISHED'
		  AND r.search_vector @@ to_tsquery('simple', unaccent($1))`

	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, tsQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("recipe.Search count: %w", err)
	}

	if total == 0 {
		return []Dish{}, 0, nil
	}

	dataSQL := `
		SELECT r.id, r.external_id, r.name, r.slug, r.description, r.image_url,
		       r.prep_time, r.cook_time, r.total_time, r.servings, r.difficulty,
		       r.status, r.dish_type_id, r.region_id, r.main_ingredient_id, r.meal_type_id,
		       r.source_url, r.author_note, r.view_count, r.favorite_count,
		       r.last_synced_at, r.created_at, r.updated_at, r.deleted_at
		FROM recipes r
		WHERE r.deleted_at IS NULL
		  AND r.status = 'PUBLISHED'
		  AND r.search_vector @@ to_tsquery('simple', unaccent($1))
		ORDER BY ts_rank(r.search_vector, to_tsquery('simple', unaccent($1))) DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, dataSQL, tsQuery, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("recipe.Search query: %w", err)
	}
	defer rows.Close()

	var dishes []Dish
	for rows.Next() {
		d, err := scanDish(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("recipe.Search scan: %w", err)
		}
		dishes = append(dishes, d)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("recipe.Search rows: %w", err)
	}

	return dishes, total, nil
}

func (r *pgDishRepo) GetAllPublishedIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id FROM recipes WHERE status = 'PUBLISHED' AND deleted_at IS NULL")
	if err != nil {
		return nil, fmt.Errorf("recipe.GetAllPublishedIDs: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("recipe.GetAllPublishedIDs scan: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *pgDishRepo) UpsertFromSync(ctx context.Context, dish *Dish, ingredients []Ingredient, steps []Step, tagIDs []uuid.UUID) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("recipe.UpsertFromSync begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err.Error() != "tx is closed" {
			slog.Error("recipe.UpsertFromSync: rollback failed", "error", err)
		}
	}()

	// Upsert the dish by external_id.
	now := time.Now().UTC()
	upsertDishSQL := `
		INSERT INTO recipes (
			id, external_id, name, slug, description, image_url,
			prep_time, cook_time, servings, difficulty, status,
			dish_type_id, region_id, main_ingredient_id, meal_type_id,
			source_url, author_note, last_synced_at, created_at, updated_at
		) VALUES (
			COALESCE((SELECT id FROM recipes WHERE external_id = $1), gen_random_uuid()),
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $17, $17
		)
		ON CONFLICT (external_id) DO UPDATE SET
			name = EXCLUDED.name,
			slug = EXCLUDED.slug,
			description = EXCLUDED.description,
			image_url = EXCLUDED.image_url,
			prep_time = EXCLUDED.prep_time,
			cook_time = EXCLUDED.cook_time,
			servings = EXCLUDED.servings,
			difficulty = EXCLUDED.difficulty,
			status = EXCLUDED.status,
			dish_type_id = EXCLUDED.dish_type_id,
			region_id = EXCLUDED.region_id,
			main_ingredient_id = EXCLUDED.main_ingredient_id,
			meal_type_id = EXCLUDED.meal_type_id,
			source_url = EXCLUDED.source_url,
			author_note = EXCLUDED.author_note,
			last_synced_at = EXCLUDED.last_synced_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id, (xmax = 0) AS is_new`

	var dishID uuid.UUID
	var isNew bool
	err = tx.QueryRow(ctx, upsertDishSQL,
		dish.ExternalID, dish.Name, dish.Slug, dish.Description, dish.ImageURL,
		dish.PrepTime, dish.CookTime, dish.Servings, dish.Difficulty, dish.Status,
		dish.DishTypeID, dish.RegionID, dish.MainIngredientID, dish.MealTypeID,
		dish.SourceURL, dish.AuthorNote, now,
	).Scan(&dishID, &isNew)
	if err != nil {
		return false, fmt.Errorf("recipe.UpsertFromSync upsert dish: %w", err)
	}

	// Replace ingredients: delete old, insert new.
	if _, err := tx.Exec(ctx, "DELETE FROM recipe_ingredients WHERE recipe_id = $1", dishID); err != nil {
		return false, fmt.Errorf("recipe.UpsertFromSync delete ingredients: %w", err)
	}
	for _, ing := range ingredients {
		_, err := tx.Exec(ctx, `
			INSERT INTO recipe_ingredients (id, recipe_id, ingredient_id, name, amount, unit, note, is_optional, group_name, sort_order)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			dishID, ing.IngredientID, ing.Name, ing.Amount, ing.Unit, ing.Note,
			ing.IsOptional, ing.GroupName, ing.SortOrder,
		)
		if err != nil {
			return false, fmt.Errorf("recipe.UpsertFromSync insert ingredient: %w", err)
		}
	}

	// Replace steps: delete old, insert new.
	if _, err := tx.Exec(ctx, "DELETE FROM recipe_steps WHERE recipe_id = $1", dishID); err != nil {
		return false, fmt.Errorf("recipe.UpsertFromSync delete steps: %w", err)
	}
	for _, s := range steps {
		_, err := tx.Exec(ctx, `
			INSERT INTO recipe_steps (id, recipe_id, step_number, title, description, image_url, duration, sort_order)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)`,
			dishID, s.StepNumber, s.Title, s.Description, s.ImageURL, s.Duration, s.SortOrder,
		)
		if err != nil {
			return false, fmt.Errorf("recipe.UpsertFromSync insert step: %w", err)
		}
	}

	// Replace tags: delete old, insert new.
	if _, err := tx.Exec(ctx, "DELETE FROM recipe_tags WHERE recipe_id = $1", dishID); err != nil {
		return false, fmt.Errorf("recipe.UpsertFromSync delete tags: %w", err)
	}
	for _, tagID := range tagIDs {
		_, err := tx.Exec(ctx,
			"INSERT INTO recipe_tags (recipe_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			dishID, tagID,
		)
		if err != nil {
			return false, fmt.Errorf("recipe.UpsertFromSync insert tag: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("recipe.UpsertFromSync commit: %w", err)
	}

	slog.Info("recipe.UpsertFromSync completed", "dish_id", dishID, "external_id", dish.ExternalID, "is_new", isNew)
	return isNew, nil
}

func (r *pgDishRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx,
		"UPDATE recipes SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return fmt.Errorf("recipe.SoftDelete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *pgDishRepo) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE recipes SET view_count = view_count + 1 WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("recipe.IncrementViewCount: %w", err)
	}
	return nil
}

// ============================================================
// CategoryRepository implementation
// ============================================================

type pgCategoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &pgCategoryRepo{pool: pool}
}

func (r *pgCategoryRepo) List(ctx context.Context, categoryType string) ([]Category, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if categoryType != "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, name, slug, type, icon_url, sort_order, is_active, created_at, updated_at
			FROM categories
			WHERE type = $1 AND is_active = true
			ORDER BY sort_order`, categoryType)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, name, slug, type, icon_url, sort_order, is_active, created_at, updated_at
			FROM categories
			WHERE is_active = true
			ORDER BY type, sort_order`)
	}
	if err != nil {
		return nil, fmt.Errorf("recipe.CategoryList: %w", err)
	}
	defer rows.Close()

	var cats []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.Type, &c.IconURL,
			&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("recipe.CategoryList scan: %w", err)
		}
		cats = append(cats, c)
	}
	if cats == nil {
		cats = []Category{}
	}
	return cats, rows.Err()
}

func (r *pgCategoryRepo) GetBySlug(ctx context.Context, slug string) (*Category, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, type, icon_url, sort_order, is_active, created_at, updated_at
		FROM categories
		WHERE slug = $1`, slug)

	var c Category
	if err := row.Scan(
		&c.ID, &c.Name, &c.Slug, &c.Type, &c.IconURL,
		&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("recipe.CategoryGetBySlug: %w", err)
	}
	return &c, nil
}

// ============================================================
// TagRepository implementation
// ============================================================

type pgTagRepo struct {
	pool *pgxpool.Pool
}

func NewTagRepository(pool *pgxpool.Pool) TagRepository {
	return &pgTagRepo{pool: pool}
}

func (r *pgTagRepo) List(ctx context.Context) ([]Tag, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name, slug FROM tags ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("recipe.TagList: %w", err)
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, fmt.Errorf("recipe.TagList scan: %w", err)
		}
		tags = append(tags, t)
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, rows.Err()
}

func (r *pgTagRepo) GetBySlug(ctx context.Context, slug string) (*Tag, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, name, slug FROM tags WHERE slug = $1", slug)
	var t Tag
	if err := row.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("recipe.TagGetBySlug: %w", err)
	}
	return &t, nil
}

// ============================================================
// Helpers
// ============================================================

// buildDishFilterClause builds a WHERE clause and args slice from DishFilter.
// It always includes the base conditions (status = PUBLISHED, deleted_at IS NULL).
func buildDishFilterClause(filter DishFilter) (string, []any) {
	conditions := []string{
		"r.status = 'PUBLISHED'",
		"r.deleted_at IS NULL",
	}
	var args []any
	argIdx := 1

	if filter.DishTypeID != nil {
		conditions = append(conditions, fmt.Sprintf("r.dish_type_id = $%d", argIdx))
		args = append(args, *filter.DishTypeID)
		argIdx++
	}
	if filter.RegionID != nil {
		conditions = append(conditions, fmt.Sprintf("r.region_id = $%d", argIdx))
		args = append(args, *filter.RegionID)
		argIdx++
	}
	if filter.MainIngredientID != nil {
		conditions = append(conditions, fmt.Sprintf("r.main_ingredient_id = $%d", argIdx))
		args = append(args, *filter.MainIngredientID)
		argIdx++
	}
	if filter.MealTypeID != nil {
		conditions = append(conditions, fmt.Sprintf("r.meal_type_id = $%d", argIdx))
		args = append(args, *filter.MealTypeID)
		argIdx++
	}
	if filter.Difficulty != nil {
		conditions = append(conditions, fmt.Sprintf("r.difficulty = $%d", argIdx))
		args = append(args, *filter.Difficulty)
		argIdx++
	}
	if filter.MaxCookTime != nil {
		conditions = append(conditions, fmt.Sprintf("r.cook_time <= $%d", argIdx))
		args = append(args, *filter.MaxCookTime)
		argIdx++
	}
	if len(filter.Tags) > 0 {
		// Subquery: recipe must have ALL specified tags (by slug).
		placeholders := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, tag)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf(`
			r.id IN (
				SELECT rt.recipe_id
				FROM recipe_tags rt
				INNER JOIN tags t ON t.id = rt.tag_id
				WHERE t.slug IN (%s)
				GROUP BY rt.recipe_id
				HAVING COUNT(DISTINCT t.slug) = %d
			)`, strings.Join(placeholders, ","), len(filter.Tags)))
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

// scanDish scans a Dish from a pgx.Rows row.
func scanDish(rows pgx.Rows) (Dish, error) {
	var d Dish
	err := rows.Scan(
		&d.ID, &d.ExternalID, &d.Name, &d.Slug, &d.Description, &d.ImageURL,
		&d.PrepTime, &d.CookTime, &d.TotalTime, &d.Servings, &d.Difficulty,
		&d.Status, &d.DishTypeID, &d.RegionID, &d.MainIngredientID, &d.MealTypeID,
		&d.SourceURL, &d.AuthorNote, &d.ViewCount, &d.FavoriteCount,
		&d.LastSyncedAt, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)
	return d, err
}

// scanDishRow scans a Dish from a single pgx.Row.
func scanDishRow(row pgx.Row) (Dish, error) {
	var d Dish
	err := row.Scan(
		&d.ID, &d.ExternalID, &d.Name, &d.Slug, &d.Description, &d.ImageURL,
		&d.PrepTime, &d.CookTime, &d.TotalTime, &d.Servings, &d.Difficulty,
		&d.Status, &d.DishTypeID, &d.RegionID, &d.MainIngredientID, &d.MealTypeID,
		&d.SourceURL, &d.AuthorNote, &d.ViewCount, &d.FavoriteCount,
		&d.LastSyncedAt, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)
	return d, err
}
