package recipe

import (
	"context"

	"github.com/google/uuid"
)

type DishFilter struct {
	DishTypeID       *uuid.UUID `json:"dishTypeId,omitempty"`
	RegionID         *uuid.UUID `json:"regionId,omitempty"`
	MainIngredientID *uuid.UUID `json:"mainIngredientId,omitempty"`
	MealTypeID       *uuid.UUID `json:"mealTypeId,omitempty"`
	Difficulty       *string    `json:"difficulty,omitempty"`
	MaxCookTime      *int       `json:"maxCookTime,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
	Page             int        `json:"page"`
	PageSize         int        `json:"pageSize"`
}

type DishRepository interface {
	List(ctx context.Context, filter DishFilter) ([]Dish, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*DishDetail, error)
	GetBySlug(ctx context.Context, slug string) (*DishDetail, error)
	GetRandom(ctx context.Context, filter DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error)
	Search(ctx context.Context, query string, page, pageSize int) ([]Dish, int64, error)
	GetAllPublishedIDs(ctx context.Context) ([]uuid.UUID, error)
	UpsertFromSync(ctx context.Context, dish *Dish, ingredients []Ingredient, steps []Step, tagIDs []uuid.UUID) (isNew bool, err error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
}

type CategoryRepository interface {
	List(ctx context.Context, categoryType string) ([]Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
}

type TagRepository interface {
	List(ctx context.Context) ([]Tag, error)
	GetBySlug(ctx context.Context, slug string) (*Tag, error)
}

type DishServicePort interface {
	ListDishes(ctx context.Context, filter DishFilter) ([]Dish, int64, error)
	GetDish(ctx context.Context, id uuid.UUID) (*DishDetail, error)
	GetDishBySlug(ctx context.Context, slug string) (*DishDetail, error)
	GetRandomDish(ctx context.Context, filter DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error)
	SearchDishes(ctx context.Context, query string, page, pageSize int) ([]Dish, int64, error)
	ListCategories(ctx context.Context, categoryType string) ([]Category, error)
	ListTags(ctx context.Context) ([]Tag, error)
}

type ContentSource interface {
	FetchPublishedRecipes(ctx context.Context) ([]Dish, error)
	FetchRecipeContent(ctx context.Context, externalID string) ([]Ingredient, []Step, error)
}
