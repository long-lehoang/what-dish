'use client';

import { useState } from 'react';
import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import { useServingScale } from '../hooks/useServingScale';
import type { Ingredient } from '@features/dish';

interface IngredientListProps {
  ingredients: Ingredient[];
  originalServings: number;
  className?: string;
}

export function IngredientList({ ingredients, originalServings, className }: IngredientListProps) {
  const { servings, setServings, scaleAmount } = useServingScale(originalServings);
  const [checked, setChecked] = useState<Set<string>>(new Set());

  function toggleChecked(id: string) {
    setChecked((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }

  return (
    <section className={cn('px-4 py-4', className)}>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100">
          Nguyên liệu
        </h2>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setServings(servings - 1)}
            disabled={servings <= 1}
            aria-label="Giảm khẩu phần"
          >
            -
          </Button>
          <span className="min-w-[2.5rem] text-center text-sm font-medium text-gray-900 dark:text-gray-100">
            {servings} phần
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setServings(servings + 1)}
            disabled={servings >= 20}
            aria-label="Tăng khẩu phần"
          >
            +
          </Button>
        </div>
      </div>

      <ul className="space-y-2">
        {ingredients.map((ingredient) => {
          const isChecked = checked.has(ingredient.id);
          const scaledAmount = scaleAmount(ingredient.amount);

          return (
            <li key={ingredient.id}>
              <label
                className={cn(
                  'flex cursor-pointer items-center gap-3 rounded-lg px-3 py-2 transition-colors hover:bg-gray-50 dark:hover:bg-gray-800',
                  ingredient.isOptional && 'opacity-70',
                )}
              >
                <input
                  type="checkbox"
                  checked={isChecked}
                  onChange={() => toggleChecked(ingredient.id)}
                  className="h-5 w-5 rounded border-gray-300 text-primary focus:ring-primary"
                />
                <span
                  className={cn(
                    'flex-1 text-sm text-gray-800 dark:text-gray-200',
                    isChecked && 'line-through opacity-50',
                  )}
                >
                  {scaledAmount !== null && (
                    <span className="font-medium">
                      {scaledAmount}
                      {ingredient.unit ? ` ${ingredient.unit}` : ''}{' '}
                    </span>
                  )}
                  {ingredient.name}
                  {ingredient.isOptional && (
                    <span className="ml-1 text-xs text-gray-500">(tùy chọn)</span>
                  )}
                </span>
              </label>
            </li>
          );
        })}
      </ul>
    </section>
  );
}
