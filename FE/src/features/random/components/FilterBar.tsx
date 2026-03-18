'use client';

import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import type { DishCategory } from '@features/dish/types';
import type { RandomFilters } from '../types';

interface FilterBarProps {
  filters: RandomFilters;
  onFilterChange: (filters: RandomFilters) => void;
  className?: string;
}

const CATEGORY_OPTIONS: { value: DishCategory; label: string; icon: string }[] = [
  { value: 'com', label: 'Cơm', icon: '🍚' },
  { value: 'bun_pho', label: 'Bún/Phở', icon: '🍜' },
  { value: 'lau', label: 'Lẩu', icon: '🍲' },
  { value: 'xao', label: 'Xào', icon: '🥘' },
  { value: 'nuong', label: 'Nướng', icon: '🍖' },
  { value: 'chien', label: 'Chiên', icon: '🍳' },
  { value: 'hap', label: 'Hấp', icon: '🥟' },
  { value: 'soup', label: 'Canh/Súp', icon: '🥣' },
  { value: 'salad', label: 'Salad', icon: '🥗' },
  { value: 'do_uong', label: 'Đồ uống', icon: '🧃' },
  { value: 'trang_mieng', label: 'Tráng miệng', icon: '🍮' },
];

const DIFFICULTY_OPTIONS: { value: number; label: string }[] = [
  { value: 1, label: 'Dễ' },
  { value: 2, label: 'Dễ-TB' },
  { value: 3, label: 'Trung bình' },
  { value: 4, label: 'Khó' },
  { value: 5, label: 'Rất khó' },
];

const TIME_OPTIONS: { value: number; label: string }[] = [
  { value: 15, label: '< 15 phút' },
  { value: 30, label: '< 30 phút' },
  { value: 60, label: '< 1 giờ' },
  { value: 120, label: '< 2 giờ' },
];

const DIETARY_OPTIONS: string[] = ['Chay', 'Không gluten', 'Ít calo', 'Keto'];

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
  if (filters.category) count++;
  if (filters.difficulty) count++;
  if (filters.maxTime) count++;
  if (filters.dietary?.length) count += filters.dietary.length;
  return count;
}

export function FilterBar({ filters, onFilterChange, className }: FilterBarProps) {
  const [expanded, setExpanded] = useState(false);
  const count = activeFilterCount(filters);

  const toggleCategory = (cat: DishCategory) => {
    onFilterChange({ ...filters, category: filters.category === cat ? undefined : cat });
  };

  const toggleDifficulty = (diff: number) => {
    onFilterChange({ ...filters, difficulty: filters.difficulty === diff ? undefined : diff });
  };

  const toggleTime = (time: number) => {
    onFilterChange({ ...filters, maxTime: filters.maxTime === time ? undefined : time });
  };

  const toggleDietary = (diet: string) => {
    const current = filters.dietary ?? [];
    const next = current.includes(diet) ? current.filter((d) => d !== diet) : [...current, diet];
    onFilterChange({ ...filters, dietary: next.length > 0 ? next : undefined });
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
      <div className="scrollbar-hide flex gap-2 overflow-x-auto pb-1">
        {CATEGORY_OPTIONS.map((opt) => (
          <Chip
            key={opt.value}
            label={opt.label}
            icon={opt.icon}
            isActive={filters.category === opt.value}
            onClick={() => toggleCategory(opt.value)}
          />
        ))}
      </div>

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
                      isActive={filters.maxTime === opt.value}
                      onClick={() => toggleTime(opt.value)}
                    />
                  ))}
                </div>
              </div>

              {/* Dietary */}
              <div>
                <p className="mb-1.5 text-[11px] font-medium uppercase tracking-wider text-gray-400 dark:text-gray-500">
                  Chế độ ăn
                </p>
                <div className="flex flex-wrap gap-2">
                  {DIETARY_OPTIONS.map((opt) => (
                    <Chip
                      key={`diet-${opt}`}
                      label={opt}
                      isActive={filters.dietary?.includes(opt) ?? false}
                      onClick={() => toggleDietary(opt)}
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
