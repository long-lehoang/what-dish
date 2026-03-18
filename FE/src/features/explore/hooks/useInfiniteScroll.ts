'use client';

import { useCallback, useEffect, useRef } from 'react';

interface UseInfiniteScrollOptions {
  onLoadMore: () => void;
  hasMore: boolean;
  isLoading: boolean;
  rootMargin?: string;
}

interface UseInfiniteScrollReturn {
  sentinelRef: (node: HTMLElement | null) => void;
}

export function useInfiniteScroll({
  onLoadMore,
  hasMore,
  isLoading,
  rootMargin = '200px',
}: UseInfiniteScrollOptions): UseInfiniteScrollReturn {
  const observerRef = useRef<IntersectionObserver | null>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const callbackRef = useRef(onLoadMore);
  callbackRef.current = onLoadMore;

  const hasMoreRef = useRef(hasMore);
  hasMoreRef.current = hasMore;

  const isLoadingRef = useRef(isLoading);
  isLoadingRef.current = isLoading;

  const sentinelRef = useCallback(
    (node: HTMLElement | null) => {
      // Cleanup previous observer
      if (observerRef.current) {
        observerRef.current.disconnect();
        observerRef.current = null;
      }

      if (!node) return;

      observerRef.current = new IntersectionObserver(
        (entries) => {
          const entry = entries[0];
          if (entry?.isIntersecting && hasMoreRef.current && !isLoadingRef.current) {
            // Debounce trigger
            if (debounceRef.current) {
              clearTimeout(debounceRef.current);
            }
            debounceRef.current = setTimeout(() => {
              callbackRef.current();
            }, 100);
          }
        },
        { rootMargin },
      );

      observerRef.current.observe(node);
    },
    [rootMargin],
  );

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, []);

  return { sentinelRef };
}
