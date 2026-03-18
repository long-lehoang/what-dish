package recipe

import (
	"time"

	"github.com/google/uuid"
)

type Dish struct {
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
	SearchVector     *string    `json:"-"`
	LastSyncedAt     *time.Time `json:"lastSyncedAt,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
}

type DishDetail struct {
	Dish
	Ingredients    []Ingredient `json:"ingredients"`
	Steps          []Step       `json:"steps"`
	Tags           []Tag        `json:"tags"`
	DishType       *Category    `json:"dishType,omitempty"`
	Region         *Category    `json:"region,omitempty"`
	MainIngredient *Category    `json:"mainIngredient,omitempty"`
	MealType       *Category    `json:"mealType,omitempty"`
}

type Ingredient struct {
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

type Step struct {
	ID          uuid.UUID `json:"id"`
	RecipeID    uuid.UUID `json:"recipeId"`
	StepNumber  int       `json:"stepNumber"`
	Title       *string   `json:"title,omitempty"`
	Description string    `json:"description"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
	SortOrder   int       `json:"sortOrder"`
}

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Type      string    `json:"type"`
	IconURL   *string   `json:"iconUrl,omitempty"`
	SortOrder int       `json:"sortOrder"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}
