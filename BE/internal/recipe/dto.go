package recipe

import (
	"time"

	"github.com/google/uuid"
)

// ---------- Request DTOs ----------

type ListDishesRequest struct {
	DishType       string `form:"dish_type"`
	Region         string `form:"region"`
	MainIngredient string `form:"main_ingredient"`
	MealType       string `form:"meal_type"`
	Difficulty     string `form:"difficulty"`
	MaxCookTime    int    `form:"max_cook_time"`
	Tags           string `form:"tags"`
	Page           int    `form:"page"`
	Limit          int    `form:"limit"`
}

type SearchRequest struct {
	Q     string `form:"q"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}

// ---------- Response DTOs ----------

type DishResponse struct {
	ID               uuid.UUID  `json:"id"`
	ExternalID       *string    `json:"externalId,omitempty"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Description      *string    `json:"description,omitempty"`
	ImageURL         *string    `json:"imageUrl,omitempty"`
	PrepTime         *int       `json:"prepTime,omitempty"`
	CookTime         *int       `json:"cookTime,omitempty"`
	TotalTime        *int       `json:"totalTime,omitempty"`
	Servings         *int       `json:"servings,omitempty"`
	Difficulty       *string    `json:"difficulty,omitempty"`
	Status           string     `json:"status"`
	DishTypeID       *uuid.UUID `json:"dishTypeId,omitempty"`
	RegionID         *uuid.UUID `json:"regionId,omitempty"`
	MainIngredientID *uuid.UUID `json:"mainIngredientId,omitempty"`
	MealTypeID       *uuid.UUID `json:"mealTypeId,omitempty"`
	SourceURL        *string    `json:"sourceUrl,omitempty"`
	AuthorNote       *string    `json:"authorNote,omitempty"`
	ViewCount        int        `json:"viewCount"`
	FavoriteCount    int        `json:"favoriteCount"`
	LastSyncedAt     *time.Time `json:"lastSyncedAt,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type IngredientResponse struct {
	ID           uuid.UUID  `json:"id"`
	RecipeID     uuid.UUID  `json:"recipeId"`
	IngredientID *uuid.UUID `json:"ingredientId,omitempty"`
	Name         string     `json:"name"`
	Amount       *float64   `json:"amount,omitempty"`
	Unit         *string    `json:"unit,omitempty"`
	Note         *string    `json:"note,omitempty"`
	IsOptional   bool       `json:"isOptional"`
	GroupName    *string    `json:"groupName,omitempty"`
	SortOrder    int        `json:"sortOrder"`
}

type StepResponse struct {
	ID          uuid.UUID `json:"id"`
	RecipeID    uuid.UUID `json:"recipeId"`
	StepNumber  int       `json:"stepNumber"`
	Title       *string   `json:"title,omitempty"`
	Description string    `json:"description"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
	SortOrder   int       `json:"sortOrder"`
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Type      string    `json:"type"`
	IconURL   *string   `json:"iconUrl,omitempty"`
	SortOrder int       `json:"sortOrder"`
	IsActive  bool      `json:"isActive"`
}

type DishDetailResponse struct {
	DishResponse
	Ingredients    []IngredientResponse `json:"ingredients"`
	Steps          []StepResponse       `json:"steps"`
	Tags           []TagResponse        `json:"tags"`
	DishType       *CategoryResponse    `json:"dishType,omitempty"`
	Region         *CategoryResponse    `json:"region,omitempty"`
	MainIngredient *CategoryResponse    `json:"mainIngredient,omitempty"`
	MealType       *CategoryResponse    `json:"mealType,omitempty"`
}

// ---------- Mappers ----------

func ToDishResponse(d Dish) DishResponse {
	return DishResponse{
		ID:               d.ID,
		ExternalID:       d.ExternalID,
		Name:             d.Name,
		Slug:             d.Slug,
		Description:      d.Description,
		ImageURL:         d.ImageURL,
		PrepTime:         d.PrepTime,
		CookTime:         d.CookTime,
		TotalTime:        d.TotalTime,
		Servings:         d.Servings,
		Difficulty:       d.Difficulty,
		Status:           d.Status,
		DishTypeID:       d.DishTypeID,
		RegionID:         d.RegionID,
		MainIngredientID: d.MainIngredientID,
		MealTypeID:       d.MealTypeID,
		SourceURL:        d.SourceURL,
		AuthorNote:       d.AuthorNote,
		ViewCount:        d.ViewCount,
		FavoriteCount:    d.FavoriteCount,
		LastSyncedAt:     d.LastSyncedAt,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}

func ToDishListResponse(dishes []Dish) []DishResponse {
	result := make([]DishResponse, len(dishes))
	for i, d := range dishes {
		result[i] = ToDishResponse(d)
	}
	return result
}

func toCategoryResponse(c *Category) *CategoryResponse {
	if c == nil {
		return nil
	}
	return &CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		Type:      c.Type,
		IconURL:   c.IconURL,
		SortOrder: c.SortOrder,
		IsActive:  c.IsActive,
	}
}

func ToCategoryListResponse(cats []Category) []CategoryResponse {
	result := make([]CategoryResponse, len(cats))
	for i, c := range cats {
		result[i] = CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			Type:      c.Type,
			IconURL:   c.IconURL,
			SortOrder: c.SortOrder,
			IsActive:  c.IsActive,
		}
	}
	return result
}

func ToTagListResponse(tags []Tag) []TagResponse {
	result := make([]TagResponse, len(tags))
	for i, t := range tags {
		result[i] = TagResponse{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		}
	}
	return result
}

func ToDishDetailResponse(d *DishDetail) DishDetailResponse {
	ingredients := make([]IngredientResponse, len(d.Ingredients))
	for i, ing := range d.Ingredients {
		ingredients[i] = IngredientResponse{
			ID:           ing.ID,
			RecipeID:     ing.RecipeID,
			IngredientID: ing.IngredientID,
			Name:         ing.Name,
			Amount:       ing.Amount,
			Unit:         ing.Unit,
			Note:         ing.Note,
			IsOptional:   ing.IsOptional,
			GroupName:    ing.GroupName,
			SortOrder:    ing.SortOrder,
		}
	}

	steps := make([]StepResponse, len(d.Steps))
	for i, s := range d.Steps {
		steps[i] = StepResponse{
			ID:          s.ID,
			RecipeID:    s.RecipeID,
			StepNumber:  s.StepNumber,
			Title:       s.Title,
			Description: s.Description,
			ImageURL:    s.ImageURL,
			Duration:    s.Duration,
			SortOrder:   s.SortOrder,
		}
	}

	tags := make([]TagResponse, len(d.Tags))
	for i, t := range d.Tags {
		tags[i] = TagResponse{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		}
	}

	return DishDetailResponse{
		DishResponse:   ToDishResponse(d.Dish),
		Ingredients:    ingredients,
		Steps:          steps,
		Tags:           tags,
		DishType:       toCategoryResponse(d.DishType),
		Region:         toCategoryResponse(d.Region),
		MainIngredient: toCategoryResponse(d.MainIngredient),
		MealType:       toCategoryResponse(d.MealType),
	}
}
