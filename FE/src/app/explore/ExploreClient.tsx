'use client';

import { useCallback } from 'react';
import { SearchBar, FilterSheet, DishGrid, ActiveFilters, useDishSearch } from '@/features/explore';
import type { Category, Dish, DishFilters } from '@/features/dish';

interface ExploreClientProps {
  initialDishes: Dish[];
  initialTotal: number;
  categories?: Category[];
}

export function ExploreClient({
  initialDishes,
  initialTotal,
  categories = [],
}: ExploreClientProps) {
  const { dishes, filters, setFilters, isLoading, loadMore, hasMore } = useDishSearch(
    initialDishes,
    {},
  );

  const handleRemoveFilter = useCallback(
    (key: keyof DishFilters) => {
      setFilters((prev) => {
        const next = { ...prev };
        delete next[key];
        return { ...next, page: 1 };
      });
    },
    [setFilters],
  );

  const handleClearAll = useCallback(() => {
    setFilters({ search: filters.search, page: 1 });
  }, [setFilters, filters.search]);

  return (
    <div>
      <div className="mb-4 flex items-center gap-3">
        <div className="flex-1">
          <SearchBar
            value={filters.search ?? ''}
            onChange={(search) => setFilters({ ...filters, search, page: 1 })}
          />
        </div>
        <FilterSheet filters={filters} onApply={setFilters} categories={categories} />
      </div>

      <ActiveFilters
        filters={filters}
        onRemove={handleRemoveFilter}
        onClearAll={handleClearAll}
        categories={categories}
      />

      <p className="mb-4 text-sm text-gray-500 dark:text-gray-400">
        {initialTotal > 0 ? `${initialTotal} món` : ''}
      </p>

      <DishGrid
        dishes={dishes.length > 0 ? dishes : initialDishes}
        isLoading={isLoading}
        hasMore={hasMore}
        onLoadMore={loadMore}
      />
    </div>
  );
}
