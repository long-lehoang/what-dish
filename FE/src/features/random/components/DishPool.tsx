'use client';

import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { apiClient } from '@shared/lib/api-client';
import { DishCard } from '@features/dish';
import type { Dish } from '@features/dish/types';
import { Skeleton } from '@shared/ui';
import type { RandomFilters } from '../types';

interface DishPoolProps {
  filters: RandomFilters;
  className?: string;
}

const staggerContainer = {
  hidden: {},
  visible: { transition: { staggerChildren: 0.05 } },
};

const staggerItem = {
  hidden: { opacity: 0, y: 16, scale: 0.96 },
  visible: {
    opacity: 1,
    y: 0,
    scale: 1,
    transition: { type: 'spring', stiffness: 300, damping: 24 },
  },
};

function buildQuery(filters: RandomFilters): string {
  const params = new URLSearchParams();
  if (filters.dishType) params.set('dish_type', filters.dishType);
  if (filters.difficulty) params.set('difficulty', filters.difficulty);
  if (filters.maxCookTime) params.set('max_cook_time', String(filters.maxCookTime));
  params.set('pageSize', '50');
  const qs = params.toString();
  return `/api/v1/recipes${qs ? `?${qs}` : ''}`;
}

export function DishPool({ filters, className }: DishPoolProps) {
  const [dishes, setDishes] = useState<Dish[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);

    apiClient
      .getList<Dish>(buildQuery(filters))
      .then((res) => {
        if (!cancelled) {
          setDishes(res.data);
          setLoading(false);
        }
      })
      .catch(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [filters]);

  return (
    <section className={className}>
      {/* Section header */}
      <div className="mb-4 flex items-center gap-3">
        <div className="h-px flex-1 bg-gradient-to-r from-transparent via-gray-200 to-transparent dark:via-gray-700" />
        <h2 className="flex items-center gap-2 text-sm font-medium text-gray-500 dark:text-gray-400">
          <span>🍽️</span>
          {loading ? (
            'Đang tải...'
          ) : (
            <>
              <span className="font-bold text-gray-800 dark:text-gray-200">{dishes.length}</span>{' '}
              món trong giỏ xoay
            </>
          )}
        </h2>
        <div className="h-px flex-1 bg-gradient-to-r from-transparent via-gray-200 to-transparent dark:via-gray-700" />
      </div>

      {/* Loading skeletons */}
      {loading && (
        <div className="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="space-y-2">
              <Skeleton className="aspect-[4/3] w-full rounded-xl" />
              <Skeleton className="h-4 w-3/4 rounded" />
              <Skeleton className="h-3 w-1/2 rounded" />
            </div>
          ))}
        </div>
      )}

      {/* Empty state */}
      {!loading && dishes.length === 0 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="rounded-2xl bg-white/50 py-16 text-center dark:bg-dark-card/30"
        >
          <span className="mb-3 inline-block text-4xl">🍃</span>
          <p className="text-sm text-gray-400 dark:text-gray-500">
            Không tìm thấy món nào phù hợp. Thử bỏ bớt bộ lọc nhé!
          </p>
        </motion.div>
      )}

      {/* Dish grid */}
      {!loading && dishes.length > 0 && (
        <motion.div
          className="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-4"
          variants={staggerContainer}
          initial="hidden"
          animate="visible"
          key={JSON.stringify(filters)}
        >
          {dishes.map((dish) => (
            <motion.div key={dish.id} variants={staggerItem}>
              <DishCard dish={dish} />
            </motion.div>
          ))}
        </motion.div>
      )}
    </section>
  );
}
