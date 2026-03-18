'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { apiClient } from '@shared/lib/api-client';
import { useDebounce } from '@shared/hooks';
import type { Dish, DishFilters } from '@features/dish';

interface UseDishSearchReturn {
  dishes: Dish[];
  filters: DishFilters;
  setFilters: (filters: DishFilters | ((prev: DishFilters) => DishFilters)) => void;
  isLoading: boolean;
  error: string | null;
  loadMore: () => void;
  hasMore: boolean;
}

const PAGE_SIZE = 12;

export function useDishSearch(
  initialDishes?: Dish[],
  initialFilters?: DishFilters,
): UseDishSearchReturn {
  const [dishes, setDishes] = useState<Dish[]>(initialDishes ?? []);
  const [filters, setFilters] = useState<DishFilters>(
    initialFilters ?? { page: 1, pageSize: PAGE_SIZE },
  );
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [page, setPage] = useState(1);
  const abortRef = useRef<AbortController | null>(null);

  const debouncedSearch = useDebounce(filters.search ?? '', 300);

  const fetchDishes = useCallback(
    async (currentPage: number, append: boolean) => {
      // Cancel previous request
      if (abortRef.current) {
        abortRef.current.abort();
      }

      const controller = new AbortController();
      abortRef.current = controller;

      setIsLoading(true);
      setError(null);

      try {
        const params = new URLSearchParams();
        if (filters.dishType) params.set('dish_type', filters.dishType);
        if (filters.difficulty) params.set('difficulty', filters.difficulty);
        if (filters.maxCookTime !== undefined)
          params.set('max_cook_time', String(filters.maxCookTime));
        if (filters.tags) params.set('tags', filters.tags);
        params.set('page', String(currentPage));
        params.set('pageSize', String(PAGE_SIZE));

        // Use search endpoint when search query is present
        let path: string;
        if (debouncedSearch) {
          params.set('q', debouncedSearch);
          path = `/api/v1/recipes/search?${params.toString()}`;
        } else {
          path = `/api/v1/recipes?${params.toString()}`;
        }

        const result = await apiClient.getList<Dish>(path, {
          signal: controller.signal,
        });

        if (append) {
          setDishes((prev) => [...prev, ...result.data]);
        } else {
          setDishes(result.data);
        }

        setHasMore(
          result.data.length === PAGE_SIZE && result.pagination.total > currentPage * PAGE_SIZE,
        );
      } catch (err) {
        if (err instanceof DOMException && err.name === 'AbortError') return;
        setError(err instanceof Error ? err.message : 'Không thể tải danh sách món ăn');
      } finally {
        setIsLoading(false);
      }
    },
    [debouncedSearch, filters.dishType, filters.difficulty, filters.maxCookTime, filters.tags],
  );

  // Reset and fetch on filter change
  useEffect(() => {
    setPage(1);
    void fetchDishes(1, false);
  }, [fetchDishes]);

  const loadMore = useCallback(() => {
    if (isLoading || !hasMore) return;
    const nextPage = page + 1;
    setPage(nextPage);
    void fetchDishes(nextPage, true);
  }, [isLoading, hasMore, page, fetchDishes]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (abortRef.current) {
        abortRef.current.abort();
      }
    };
  }, []);

  return {
    dishes,
    filters,
    setFilters,
    isLoading,
    error,
    loadMore,
    hasMore,
  };
}
