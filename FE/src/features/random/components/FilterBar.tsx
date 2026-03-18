'use client';

import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import type { Category, Difficulty } from '@features/dish/types';
import type { RandomFilters } from '../types';

interface FilterBarProps {
  filters: RandomFilters;
  onFilterChange: (filters: RandomFilters) => void;
  categories?: Category[];
  className?: string;
}

const DIFFICULTY_OPTIONS: { value: Difficulty; label: string }[] = [
  { value: 'EASY', label: 'Dễ' },
  { value: 'MEDIUM', label: 'Trung bình' },
  { value: 'HARD', label: 'Khó' },
];

const TIME_OPTIONS: { value: number; label: string }[] = [
  { value: 15, label: '< 15 phút' },
  { value: 30, label: '< 30 phút' },
  { value: 60, label: '< 1 giờ' },
  { value: 120, label: '< 2 giờ' },
];

interface ChipProps {
  label: string;
  icon?: string;
  isActive: boolean;
  onClick: () => void;
}

function Chip({ label, icon, isActive, onClick }: ChipProps) {
  return (
    <motion.button
      type="button"
      onClick={onClick}
      whileTap={{ scale: 0.95 }}
      className={cn(
        'shrink-0 whitespace-nowrap rounded-full px-3.5 py-2 text-xs font-medium transition-all duration-200',
        isActive
          ? 'bg-gradient-to-r from-primary to-secondary text-white shadow-sm shadow-primary/20'
          : 'bg-white/80 text-gray-600 hover:bg-white hover:shadow-sm dark:bg-dark-card dark:text-gray-300 dark:hover:bg-dark-card/80',
      )}
      aria-pressed={isActive}
    >
      {icon && <span className="mr-1">{icon}</span>}
      {label}
    </motion.button>
  );
}

function activeFilterCount(filters: RandomFilters): number {
  let count = 0;
  if (filters.dishType) count++;
  if (filters.difficulty) count++;
  if (filters.maxCookTime) count++;
  return count;
}

export function FilterBar({ filters, onFilterChange, categories = [], className }: FilterBarProps) {
  const [expanded, setExpanded] = useState(false);
  const count = activeFilterCount(filters);

  const dishTypeCategories = categories.filter((c) => c.type === 'DISH_TYPE');

  const toggleDishType = (id: string) => {
    onFilterChange({ ...filters, dishType: filters.dishType === id ? undefined : id });
  };

  const toggleDifficulty = (diff: Difficulty) => {
    onFilterChange({ ...filters, difficulty: filters.difficulty === diff ? undefined : diff });
  };

  const toggleTime = (time: number) => {
    onFilterChange({ ...filters, maxCookTime: filters.maxCookTime === time ? undefined : time });
  };

  const clearAll = () => {
    onFilterChange({});
  };

  return (
    <div className={cn('space-y-3', className)} role="group" aria-label="Bộ lọc">
      {/* Toggle row */}
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={() => setExpanded(!expanded)}
          className={cn(
            'inline-flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium transition-all',
            expanded
              ? 'bg-primary/10 text-primary dark:bg-primary/20'
              : 'bg-white/80 text-gray-600 hover:bg-white dark:bg-dark-card dark:text-gray-300',
          )}
        >
          <span>🔍</span>
          <span>Bộ lọc</span>
          {count > 0 && (
            <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-white">
              {count}
            </span>
          )}
          <motion.span
            animate={{ rotate: expanded ? 180 : 0 }}
            transition={{ duration: 0.2 }}
            className="text-xs"
          >
            ▼
          </motion.span>
        </button>

        {count > 0 && (
          <button
            type="button"
            onClick={clearAll}
            className="text-xs text-gray-400 underline-offset-2 transition-colors hover:text-secondary hover:underline dark:text-gray-500"
          >
            Xóa bộ lọc
          </button>
        )}
      </div>

      {/* Category row — always visible as horizontal scroll */}
      {dishTypeCategories.length > 0 && (
        <div className="scrollbar-hide flex gap-2 overflow-x-auto pb-1">
          {dishTypeCategories.map((cat) => (
            <Chip
              key={cat.id}
              label={cat.name}
              isActive={filters.dishType === cat.id}
              onClick={() => toggleDishType(cat.id)}
            />
          ))}
        </div>
      )}

      {/* Expandable filters */}
      <AnimatePresence>
        {expanded && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25, ease: 'easeInOut' }}
            className="overflow-hidden"
          >
            <div className="space-y-3 pt-1">
              {/* Difficulty */}
              <div>
                <p className="mb-1.5 text-[11px] font-medium uppercase tracking-wider text-gray-400 dark:text-gray-500">
                  Độ khó
                </p>
                <div className="flex flex-wrap gap-2">
                  {DIFFICULTY_OPTIONS.map((opt) => (
                    <Chip
                      key={`diff-${opt.value}`}
                      label={opt.label}
                      isActive={filters.difficulty === opt.value}
                      onClick={() => toggleDifficulty(opt.value)}
                    />
                  ))}
                </div>
              </div>

              {/* Time */}
              <div>
                <p className="mb-1.5 text-[11px] font-medium uppercase tracking-wider text-gray-400 dark:text-gray-500">
                  Thời gian nấu
                </p>
                <div className="flex flex-wrap gap-2">
                  {TIME_OPTIONS.map((opt) => (
                    <Chip
                      key={`time-${opt.value}`}
                      label={opt.label}
                      isActive={filters.maxCookTime === opt.value}
                      onClick={() => toggleTime(opt.value)}
                    />
                  ))}
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
