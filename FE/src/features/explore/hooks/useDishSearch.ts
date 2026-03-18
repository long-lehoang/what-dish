'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { apiClient } from '@shared/lib/api-client';
import { useDebounce } from '@shared/hooks';
import type { Dish, DishFilters, DishListResponse } from '@features/dish';

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
        if (debouncedSearch) params.set('search', debouncedSearch);
        if (filters.category) params.set('category', String(filters.category));
        if (filters.difficulty !== undefined) params.set('difficulty', String(filters.difficulty));
        if (filters.maxTime !== undefined) params.set('maxTime', String(filters.maxTime));
        if (filters.tags?.length) params.set('tags', filters.tags.join(','));
        params.set('page', String(currentPage));
        params.set('pageSize', String(PAGE_SIZE));

        const data = await apiClient.get<DishListResponse>(`/api/dishes?${params.toString()}`, {
          signal: controller.signal,
        });

        if (append) {
          setDishes((prev) => [...prev, ...data.dishes]);
        } else {
          setDishes(data.dishes);
        }

        setHasMore(data.dishes.length === PAGE_SIZE && data.total > currentPage * PAGE_SIZE);
      } catch (err) {
        if (err instanceof DOMException && err.name === 'AbortError') return;
        setError(
          err instanceof Error
            ? err.message
            : 'Không thể tải danh sách món ăn',
        );
      } finally {
        setIsLoading(false);
      }
    },
    [debouncedSearch, filters.category, filters.difficulty, filters.maxTime, filters.tags],
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
