import { cn } from '@shared/lib/utils';
import { DishBadge } from '@features/dish';
import { getTotalTime } from '@features/dish';
import type { Dish } from '@features/dish';

interface RecipeInfoProps {
  dish: Dish;
  className?: string;
}

export function RecipeInfo({ dish, className }: RecipeInfoProps) {
  const totalTime = getTotalTime(dish);

  return (
    <section className={cn('flex flex-wrap items-center gap-2 px-4 py-3', className)}>
      {totalTime !== null && <DishBadge type="time" value={totalTime} />}
      {dish.difficulty && <DishBadge type="difficulty" value={dish.difficulty} />}
      {dish.servings && <DishBadge type="servings" value={dish.servings} />}
    </section>
  );
}
