package recipe

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lehoanglong/whatdish/internal/platform/notion"
)

// SyncService orchestrates Notion → PostgreSQL sync.
type SyncService struct {
	client   *notion.Client
	dbID     string
	dishRepo DishRepository
	catRepo  CategoryRepository
	tagRepo  TagRepository
}

// NewSyncService creates a new SyncService.
func NewSyncService(client *notion.Client, dbID string, dishRepo DishRepository, catRepo CategoryRepository, tagRepo TagRepository) *SyncService {
	return &SyncService{
		client:   client,
		dbID:     dbID,
		dishRepo: dishRepo,
		catRepo:  catRepo,
		tagRepo:  tagRepo,
	}
}

// SyncResult contains the results of a sync operation.
type SyncResult struct {
	Added    int    `json:"added"`
	Updated  int    `json:"updated"`
	Deleted  int    `json:"deleted"`
	Errors   int    `json:"errors"`
	Duration string `json:"duration"`
}

// Sync fetches all published recipes from Notion and upserts them into PostgreSQL.
func (s *SyncService) Sync(ctx context.Context) (*SyncResult, error) {
	start := time.Now()
	result := &SyncResult{}

	slog.Info("starting Notion sync")

	// Fetch all published pages from Notion
	var allPages []notion.Page
	var cursor *string

	for {
		resp, err := s.client.QueryDatabase(ctx, s.dbID, cursor)
		if err != nil {
			return nil, fmt.Errorf("recipe.Sync: query database: %w", err)
		}

		allPages = append(allPages, resp.Results...)

		if !resp.HasMore || resp.NextCursor == nil {
			break
		}
		cursor = resp.NextCursor
	}

	slog.Info("fetched pages from Notion", "count", len(allPages))

	// Build category lookup (name → UUID)
	catLookup, err := s.buildCategoryLookup(ctx)
	if err != nil {
		return nil, fmt.Errorf("recipe.Sync: build category lookup: %w", err)
	}

	// Build tag lookup (name → UUID)
	tagLookup, err := s.buildTagLookup(ctx)
	if err != nil {
		return nil, fmt.Errorf("recipe.Sync: build tag lookup: %w", err)
	}

	// Process each page
	for _, page := range allPages {
		if err := s.processPage(ctx, page, catLookup, tagLookup, result); err != nil {
			slog.Error("sync page failed", "page_id", page.ID, "error", err)
			result.Errors++
		}
	}

	result.Duration = time.Since(start).String()

	slog.Info("Notion sync completed",
		"added", result.Added,
		"updated", result.Updated,
		"errors", result.Errors,
		"duration", result.Duration,
	)

	return result, nil
}

func (s *SyncService) processPage(ctx context.Context, page notion.Page, catLookup map[string]uuid.UUID, tagLookup map[string]uuid.UUID, result *SyncResult) error {
	// Parse page properties → parsed model
	parsed := notion.ParsePageToDish(page)

	// Convert notion.ParsedDish → recipe.Dish
	servings := parsed.Servings
	difficulty := parsed.Difficulty
	dish := Dish{
		ExternalID:  &parsed.ExternalID,
		Name:        parsed.Name,
		Slug:        parsed.Slug,
		Description: parsed.Description,
		ImageURL:    parsed.ImageURL,
		PrepTime:    parsed.PrepTime,
		CookTime:    parsed.CookTime,
		Servings:    &servings,
		Difficulty:  &difficulty,
		Status:      parsed.Status,
		SourceURL:   parsed.SourceURL,
	}

	// Resolve category IDs
	catNames := notion.GetCategoryNames(page)
	if id, ok := catLookup[catNames["DISH_TYPE"]]; ok {
		dish.DishTypeID = &id
	}
	if id, ok := catLookup[catNames["REGION"]]; ok {
		dish.RegionID = &id
	}
	if id, ok := catLookup[catNames["MAIN_INGREDIENT"]]; ok {
		dish.MainIngredientID = &id
	}
	if id, ok := catLookup[catNames["MEAL_TYPE"]]; ok {
		dish.MealTypeID = &id
	}

	// Fetch page content (blocks → ingredients + steps)
	var allBlocks []notion.Block
	var blockCursor *string
	for {
		blocksResp, err := s.client.GetBlockChildren(ctx, page.ID, blockCursor)
		if err != nil {
			return fmt.Errorf("fetch blocks for %s: %w", page.ID, err)
		}
		allBlocks = append(allBlocks, blocksResp.Results...)
		if !blocksResp.HasMore || blocksResp.NextCursor == nil {
			break
		}
		blockCursor = blocksResp.NextCursor
	}

	parsedIngredients, parsedSteps := notion.ParseBlocksToContent(allBlocks)

	// Convert to recipe models
	ingredients := make([]Ingredient, len(parsedIngredients))
	for i, pi := range parsedIngredients {
		ingredients[i] = Ingredient{
			Name:      pi.Name,
			Amount:    pi.Amount,
			Unit:      pi.Unit,
			SortOrder: pi.SortOrder,
		}
	}

	steps := make([]Step, len(parsedSteps))
	for i, ps := range parsedSteps {
		steps[i] = Step{
			StepNumber:  ps.StepNumber,
			Title:       ps.Title,
			Description: ps.Description,
			ImageURL:    ps.ImageURL,
			Duration:    ps.Duration,
			SortOrder:   ps.SortOrder,
		}
	}

	// Resolve tag IDs
	tagNames := notion.GetTagNames(page)
	var tagIDs []uuid.UUID
	for _, name := range tagNames {
		if id, ok := tagLookup[name]; ok {
			tagIDs = append(tagIDs, id)
		}
	}

	// Set sync timestamp
	now := time.Now().UTC()
	dish.LastSyncedAt = &now

	// Upsert into PostgreSQL
	isNew, err := s.dishRepo.UpsertFromSync(ctx, &dish, ingredients, steps, tagIDs)
	if err != nil {
		return fmt.Errorf("upsert dish %q: %w", dish.Name, err)
	}

	if isNew {
		result.Added++
	} else {
		result.Updated++
	}

	return nil
}

func (s *SyncService) buildCategoryLookup(ctx context.Context) (map[string]uuid.UUID, error) {
	lookup := make(map[string]uuid.UUID)

	for _, catType := range []string{"DISH_TYPE", "REGION", "MAIN_INGREDIENT", "MEAL_TYPE"} {
		cats, err := s.catRepo.List(ctx, catType)
		if err != nil {
			return nil, fmt.Errorf("list categories %s: %w", catType, err)
		}
		for _, c := range cats {
			lookup[c.Name] = c.ID
		}
	}

	return lookup, nil
}

func (s *SyncService) buildTagLookup(ctx context.Context) (map[string]uuid.UUID, error) {
	tags, err := s.tagRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}

	lookup := make(map[string]uuid.UUID, len(tags))
	for _, t := range tags {
		lookup[t.Name] = t.ID
	}
	return lookup, nil
}
