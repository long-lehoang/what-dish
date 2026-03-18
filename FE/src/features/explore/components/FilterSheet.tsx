'use client';

import { useCallback, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import { useIsMobile } from '@shared/hooks';
import type { DishCategory, DishFilters } from '@features/dish';

interface FilterSheetProps {
  filters: DishFilters;
  onApply: (filters: DishFilters) => void;
  className?: string;
}

const CATEGORIES: { value: DishCategory; label: string }[] = [
  { value: 'com', label: 'Cơm' },
  { value: 'bun_pho', label: 'Bún/Phở' },
  { value: 'lau', label: 'Lẩu' },
  { value: 'xao', label: 'Xào' },
  { value: 'nuong', label: 'Nướng' },
  { value: 'chien', label: 'Chiên' },
  { value: 'hap', label: 'Hấp' },
  { value: 'soup', label: 'Canh/Súp' },
  { value: 'salad', label: 'Salad' },
  { value: 'do_uong', label: 'Đồ uống' },
  { value: 'trang_mieng', label: 'Tráng miệng' },
];

const DIFFICULTIES: { value: number; label: string }[] = [
  { value: 1, label: 'Rất dễ' },
  { value: 2, label: 'Dễ' },
  { value: 3, label: 'Trung bình' },
  { value: 4, label: 'Khó' },
  { value: 5, label: 'Rất khó' },
];

const COOK_TIMES: { value: number; label: string }[] = [
  { value: 15, label: '< 15 phút' },
  { value: 30, label: '< 30 phút' },
  { value: 60, label: '< 1 giờ' },
  { value: 120, label: '< 2 giờ' },
];

function countActiveFilters(filters: DishFilters): number {
  let count = 0;
  if (filters.category) count++;
  if (filters.difficulty !== undefined) count++;
  if (filters.maxTime !== undefined) count++;
  return count;
}

export function FilterSheet({ filters, onApply, className }: FilterSheetProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [draft, setDraft] = useState<DishFilters>(filters);
  const isMobile = useIsMobile();
  const activeCount = countActiveFilters(filters);

  const handleOpen = useCallback(() => {
    setDraft(filters);
    setIsOpen(true);
  }, [filters]);

  const handleApply = useCallback(() => {
    onApply(draft);
    setIsOpen(false);
  }, [draft, onApply]);

  const handleClear = useCallback(() => {
    const cleared: DishFilters = {
      search: filters.search,
      page: 1,
      pageSize: filters.pageSize,
    };
    setDraft(cleared);
    onApply(cleared);
    setIsOpen(false);
  }, [filters.search, filters.pageSize, onApply]);

  return (
    <div className={className}>
      {/* Trigger button */}
      <button
        onClick={handleOpen}
        className="relative inline-flex items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:bg-dark-card dark:text-gray-300 dark:hover:bg-gray-800"
        aria-label="Bộ lọc"
      >
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
        </svg>
        Bộ lọc
        {activeCount > 0 && (
          <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary text-xs font-bold text-white">
            {activeCount}
          </span>
        )}
      </button>

      {/* Sheet overlay + content */}
      <AnimatePresence>
        {isOpen && (
          <>
            <motion.div
              className="fixed inset-0 z-40 bg-black/50"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsOpen(false)}
              aria-hidden="true"
            />
            <motion.div
              className={cn(
                'fixed z-50 overflow-y-auto bg-white dark:bg-gray-900',
                isMobile
                  ? 'inset-x-0 bottom-0 max-h-[85vh] rounded-t-2xl'
                  : 'right-0 top-0 h-full w-80',
              )}
              initial={isMobile ? { y: '100%' } : { x: '100%' }}
              animate={isMobile ? { y: 0 } : { x: 0 }}
              exit={isMobile ? { y: '100%' } : { x: '100%' }}
              transition={{ type: 'spring', damping: 30, stiffness: 300 }}
            >
              <div className="p-6">
                {/* Handle bar (mobile) */}
                {isMobile && (
                  <div className="mb-4 flex justify-center">
                    <div className="h-1 w-10 rounded-full bg-gray-300 dark:bg-gray-600" />
                  </div>
                )}

                <h3 className="mb-6 text-lg font-bold text-gray-900 dark:text-gray-100">Bộ lọc</h3>

                {/* Category */}
                <div className="mb-6">
                  <h4 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
                    Loại món
                  </h4>
                  <div className="flex flex-wrap gap-2">
                    {CATEGORIES.map((cat) => (
                      <button
                        key={cat.value}
                        onClick={() =>
                          setDraft((prev) => ({
                            ...prev,
                            category: prev.category === cat.value ? undefined : cat.value,
                          }))
                        }
                        className={cn(
                          'rounded-full px-3 py-1.5 text-xs font-medium transition-colors',
                          draft.category === cat.value
                            ? 'bg-primary text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300',
                        )}
                      >
                        {cat.label}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Difficulty */}
                <div className="mb-6">
                  <h4 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
                    Độ khó
                  </h4>
                  <div className="space-y-2">
                    {DIFFICULTIES.map((diff) => (
                      <label key={diff.value} className="flex cursor-pointer items-center gap-2">
                        <input
                          type="radio"
                          name="difficulty"
                          checked={draft.difficulty === diff.value}
                          onChange={() =>
                            setDraft((prev) => ({
                              ...prev,
                              difficulty: prev.difficulty === diff.value ? undefined : diff.value,
                            }))
                          }
                          className="h-4 w-4 border-gray-300 text-primary focus:ring-primary"
                        />
                        <span className="text-sm text-gray-700 dark:text-gray-300">
                          {diff.label}
                        </span>
                      </label>
                    ))}
                  </div>
                </div>

                {/* Cook time */}
                <div className="mb-8">
                  <h4 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
                    Thời gian nấu
                  </h4>
                  <div className="flex flex-wrap gap-2">
                    {COOK_TIMES.map((ct) => (
                      <button
                        key={ct.value}
                        onClick={() =>
                          setDraft((prev) => ({
                            ...prev,
                            maxTime: prev.maxTime === ct.value ? undefined : ct.value,
                          }))
                        }
                        className={cn(
                          'rounded-full px-3 py-1.5 text-xs font-medium transition-colors',
                          draft.maxTime === ct.value
                            ? 'bg-primary text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300',
                        )}
                      >
                        {ct.label}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Actions */}
                <div className="flex gap-3">
                  <Button variant="primary" size="md" className="flex-1" onClick={handleApply}>
                    Áp dụng
                  </Button>
                  <Button variant="ghost" size="md" onClick={handleClear}>
                    Xóa bộ lọc
                  </Button>
                </div>
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </div>
  );
}
