export type DishCategory =
  | 'com'
  | 'bun_pho'
  | 'lau'
  | 'xao'
  | 'nuong'
  | 'chien'
  | 'hap'
  | 'soup'
  | 'salad'
  | 'do_uong'
  | 'trang_mieng'
  | 'other';

export interface Dish {
  id: string;
  name: string;
  slug: string;
  description: string | null;
  imageUrl: string | null;
  thumbnail: string | null;
  category: DishCategory;
  difficulty: 1 | 2 | 3 | 4 | 5;
  prepTime: number | null;
  cookTime: number | null;
  servings: number;
  costMin: number | null;
  costMax: number | null;
  spiceLevel: 0 | 1 | 2 | 3;
  tags: string[];
  dietary: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Ingredient {
  id: string;
  dishId: string;
  name: string;
  amount: number | null;
  unit: string | null;
  isOptional: boolean;
  sortOrder: number;
}

export interface Step {
  id: string;
  dishId: string;
  stepNumber: number;
  instruction: string;
  imageUrl: string | null;
  timerSecs: number | null;
  tip: string | null;
}

export interface DishDetail extends Dish {
  ingredients: Ingredient[];
  steps: Step[];
}

export interface DishListResponse {
  dishes: Dish[];
  total: number;
  page: number;
  pageSize: number;
}

export interface DishFilters {
  category?: DishCategory;
  difficulty?: number;
  maxTime?: number;
  tags?: string[];
  dietary?: string[];
  search?: string;
  page?: number;
  pageSize?: number;
}
