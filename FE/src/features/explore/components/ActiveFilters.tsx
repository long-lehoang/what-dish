'use client';

import { cn } from '@shared/lib/utils';
import type { Category, Difficulty, DishFilters } from '@features/dish';

interface ActiveFiltersProps {
  filters: DishFilters;
  onRemove: (key: keyof DishFilters) => void;
  onClearAll: () => void;
  categories?: Category[];
  className?: string;
}

const difficultyLabels: Record<Difficulty, string> = {
  EASY: 'Dễ',
  MEDIUM: 'Trung bình',
  HARD: 'Khó',
};

interface FilterChip {
  key: keyof DishFilters;
  label: string;
}

function getActiveChips(filters: DishFilters, categories: Category[]): FilterChip[] {
  const chips: FilterChip[] = [];

  if (filters.dishType) {
    const cat = categories.find((c) => c.id === filters.dishType);
    chips.push({
      key: 'dishType',
      label: cat?.name ?? 'Loại món',
    });
  }

  if (filters.difficulty) {
    chips.push({
      key: 'difficulty',
      label: difficultyLabels[filters.difficulty] ?? filters.difficulty,
    });
  }

  if (filters.maxCookTime !== undefined) {
    chips.push({
      key: 'maxCookTime',
      label: `< ${filters.maxCookTime} phút`,
    });
  }

  return chips;
}

export function ActiveFilters({
  filters,
  onRemove,
  onClearAll,
  categories = [],
  className,
}: ActiveFiltersProps) {
  const chips = getActiveChips(filters, categories);

  if (chips.length === 0) return null;

  return (
    <div
      className={cn('flex flex-wrap items-center gap-2', className)}
      role="list"
      aria-label="Bộ lọc đang áp dụng"
    >
      {chips.map((chip) => (
        <button
          key={chip.key}
          onClick={() => onRemove(chip.key)}
          className="inline-flex items-center gap-1 rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary transition-colors hover:bg-primary/20"
          role="listitem"
          aria-label={`Xóa bộ lọc: ${chip.label}`}
        >
          {chip.label}
          <svg
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
          >
            <path d="M18 6L6 18M6 6l12 12" />
          </svg>
        </button>
      ))}

      {chips.length > 1 && (
        <button
          onClick={onClearAll}
          className="text-xs font-medium text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
        >
          Xóa tất cả
        </button>
      )}
    </div>
  );
}
