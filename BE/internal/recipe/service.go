package recipe

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type DishService struct {
	dishRepo DishRepository
	catRepo  CategoryRepository
	tagRepo  TagRepository
}

func NewDishService(dishRepo DishRepository, catRepo CategoryRepository, tagRepo TagRepository) *DishService {
	return &DishService{
		dishRepo: dishRepo,
		catRepo:  catRepo,
		tagRepo:  tagRepo,
	}
}

func (s *DishService) ListDishes(ctx context.Context, filter DishFilter) ([]Dish, int64, error) {
	dishes, total, err := s.dishRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("recipe.ListDishes: %w", err)
	}
	return dishes, total, nil
}

func (s *DishService) GetDish(ctx context.Context, id uuid.UUID) (*DishDetail, error) {
	detail, err := s.dishRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("recipe.GetDish: %w", err)
	}
	return detail, nil
}

func (s *DishService) GetDishBySlug(ctx context.Context, slug string) (*DishDetail, error) {
	detail, err := s.dishRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("recipe.GetDishBySlug: %w", err)
	}
	return detail, nil
}

func (s *DishService) GetRandomDish(ctx context.Context, filter DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error) {
	detail, err := s.dishRepo.GetRandom(ctx, filter, excludeIDs)
	if err != nil {
		return nil, fmt.Errorf("recipe.GetRandomDish: %w", err)
	}
	return detail, nil
}

func (s *DishService) SearchDishes(ctx context.Context, query string, page, pageSize int) ([]Dish, int64, error) {
	dishes, total, err := s.dishRepo.Search(ctx, query, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("recipe.SearchDishes: %w", err)
	}
	return dishes, total, nil
}

func (s *DishService) ListCategories(ctx context.Context, categoryType string) ([]Category, error) {
	cats, err := s.catRepo.List(ctx, categoryType)
	if err != nil {
		return nil, fmt.Errorf("recipe.ListCategories: %w", err)
	}
	return cats, nil
}

func (s *DishService) ListTags(ctx context.Context) ([]Tag, error) {
	tags, err := s.tagRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("recipe.ListTags: %w", err)
	}
	return tags, nil
}
