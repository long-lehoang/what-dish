import { cn, formatCurrency } from '@shared/lib/utils';
import { DishBadge } from '@features/dish';
import type { Dish } from '@features/dish';

interface RecipeInfoProps {
  dish: Dish;
  className?: string;
}

export function RecipeInfo({ dish, className }: RecipeInfoProps) {
  const totalTime = (dish.prepTime ?? 0) + (dish.cookTime ?? 0);

  const costRange =
    dish.costMin !== null && dish.costMax !== null
      ? `${formatCurrency(dish.costMin)}-${formatCurrency(dish.costMax)}`
      : null;

  return (
    <section className={cn('flex flex-wrap items-center gap-2 px-4 py-3', className)}>
      {totalTime > 0 && <DishBadge type="time" value={totalTime} />}
      <DishBadge type="difficulty" value={dish.difficulty} />
      <DishBadge type="time" value={`${dish.servings} phần`} />
      {costRange !== null && <DishBadge type="cost" value={costRange} />}
      {dish.spiceLevel > 0 && <DishBadge type="spice" value={dish.spiceLevel} />}
    </section>
  );
}
