'use client';

import { cn } from '@shared/lib/utils';
import { Skeleton } from '@shared/ui';
import { DishCard } from '@features/dish';
import type { Dish } from '@features/dish';
import { useInfiniteScroll } from '../hooks/useInfiniteScroll';

interface DishGridProps {
  dishes: Dish[];
  isLoading: boolean;
  hasMore: boolean;
  onLoadMore: () => void;
  className?: string;
}

function SkeletonGrid() {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: 6 }, (_, i) => (
        <div key={i} className="overflow-hidden rounded-2xl">
          <Skeleton className="aspect-[4/3] w-full" />
          <div className="space-y-2 p-4">
            <Skeleton className="h-5 w-3/4" />
            <div className="flex gap-2">
              <Skeleton className="h-6 w-16 rounded-full" />
              <Skeleton className="h-6 w-12 rounded-full" />
              <Skeleton className="h-6 w-20 rounded-full" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="mb-4 text-5xl" aria-hidden="true">
        &#x1F371;
      </div>
      <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
        Không tìm thấy món nào
      </h3>
      <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
        Thử thay đổi bộ lọc hoặc từ khóa tìm kiếm
      </p>
    </div>
  );
}

export function DishGrid({ dishes, isLoading, hasMore, onLoadMore, className }: DishGridProps) {
  const { sentinelRef } = useInfiniteScroll({
    onLoadMore,
    hasMore,
    isLoading,
  });

  if (isLoading && dishes.length === 0) {
    return <SkeletonGrid />;
  }

  if (!isLoading && dishes.length === 0) {
    return <EmptyState />;
  }

  return (
    <div className={cn('space-y-4', className)}>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {dishes.map((dish) => (
          <DishCard key={dish.id} dish={dish} />
        ))}
      </div>

      {/* Loading more indicator */}
      {isLoading && dishes.length > 0 && (
        <div className="flex justify-center py-4">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-gray-300 border-t-primary" />
        </div>
      )}

      {/* Sentinel for infinite scroll */}
      {hasMore && <div ref={sentinelRef} className="h-1" aria-hidden="true" />}
    </div>
  );
}
