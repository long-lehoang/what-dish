// Types matching BE DTOs from internal/recipe/dto.go

export type { Pagination } from '@shared/lib/api-client';

export type Difficulty = 'EASY' | 'MEDIUM' | 'HARD';
export type DishStatus = 'PUBLISHED' | 'DRAFT';
export type CategoryType = 'DISH_TYPE' | 'REGION' | 'MAIN_INGREDIENT' | 'MEAL_TYPE';

export interface Category {
  id: string;
  name: string;
  slug: string;
  type: CategoryType;
  iconUrl?: string;
  sortOrder: number;
  isActive: boolean;
}

export interface Tag {
  id: string;
  name: string;
  slug: string;
}

export interface Dish {
  id: string;
  externalId?: string;
  name: string;
  slug: string;
  description?: string;
  imageUrl?: string;
  prepTime?: number;
  cookTime?: number;
  totalTime?: number;
  servings?: number;
  difficulty?: Difficulty;
  status: DishStatus;
  dishTypeId?: string;
  regionId?: string;
  mainIngredientId?: string;
  mealTypeId?: string;
  sourceUrl?: string;
  authorNote?: string;
  viewCount: number;
  favoriteCount: number;
  lastSyncedAt?: string;
  createdAt: string;
  updatedAt: string;
}

/** Compute display-ready total time in minutes. */
export function getTotalTime(dish: Dish): number | null {
  return (dish.totalTime ?? (dish.prepTime ?? 0) + (dish.cookTime ?? 0)) || null;
}

export interface Ingredient {
  id: string;
  recipeId: string;
  ingredientId?: string;
  name: string;
  amount?: number;
  unit?: string;
  note?: string;
  isOptional: boolean;
  groupName?: string;
  sortOrder: number;
}

export interface Step {
  id: string;
  recipeId: string;
  stepNumber: number;
  title?: string;
  description: string;
  imageUrl?: string;
  duration?: number; // seconds
  sortOrder: number;
}

export interface DishDetail extends Dish {
  ingredients: Ingredient[];
  steps: Step[];
  tags: Tag[];
  dishType?: Category;
  region?: Category;
  mainIngredient?: Category;
  mealType?: Category;
}

export interface DishFilters {
  dishType?: string; // category UUID
  region?: string;
  mainIngredient?: string;
  mealType?: string;
  difficulty?: Difficulty;
  maxCookTime?: number;
  tags?: string; // comma-separated tag slugs
  search?: string;
  page?: number;
  pageSize?: number;
}
