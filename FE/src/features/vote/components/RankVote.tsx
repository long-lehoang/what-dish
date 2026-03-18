'use client';

import { useState } from 'react';
import { Reorder } from 'framer-motion';
import Image from 'next/image';
import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import type { Dish } from '@features/dish';

interface RankVoteProps {
  dishes: Dish[];
  onComplete: (data: unknown) => void;
  timeRemaining?: number;
}

export function RankVote({ dishes, onComplete, timeRemaining }: RankVoteProps) {
  const [items, setItems] = useState<Dish[]>(dishes);

  function handleConfirm() {
    const ranking = items.map((dish, index) => ({
      dishId: dish.id,
      rank: index + 1,
    }));
    onComplete({ ranking });
  }

  return (
    <div className="flex flex-col gap-4 px-4 py-6">
      <div className="text-center">
        <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100">
          Sắp xếp theo thứ tự yêu thích
        </h3>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Kéo thả để thay đổi vị trí</p>
        {timeRemaining !== undefined && (
          <p className="mt-1 text-xs text-gray-400">Còn {timeRemaining}s</p>
        )}
      </div>

      <Reorder.Group axis="y" values={items} onReorder={setItems} className="space-y-2">
        {items.map((dish, index) => (
          <Reorder.Item
            key={dish.id}
            value={dish}
            className="flex cursor-grab items-center gap-3 rounded-xl bg-white p-3 shadow-sm active:cursor-grabbing active:shadow-md dark:bg-dark-card"
          >
            {/* Drag handle */}
            <div className="flex shrink-0 touch-none items-center text-gray-400" aria-hidden="true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                <circle cx="9" cy="6" r="1.5" />
                <circle cx="15" cy="6" r="1.5" />
                <circle cx="9" cy="12" r="1.5" />
                <circle cx="15" cy="12" r="1.5" />
                <circle cx="9" cy="18" r="1.5" />
                <circle cx="15" cy="18" r="1.5" />
              </svg>
            </div>

            {/* Rank number */}
            <span
              className={cn(
                'flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-bold',
                index === 0
                  ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
                  : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400',
              )}
            >
              {index + 1}
            </span>

            {/* Dish image */}
            <div className="relative h-12 w-12 shrink-0 overflow-hidden rounded-lg">
              {dish.imageUrl ? (
                <Image
                  src={dish.imageUrl}
                  alt={dish.name}
                  fill
                  sizes="48px"
                  className="object-cover"
                />
              ) : (
                <div className="h-full w-full bg-gray-200 dark:bg-gray-700" />
              )}
            </div>

            {/* Dish name */}
            <span className="flex-1 text-sm font-medium text-gray-900 dark:text-gray-100">
              {dish.name}
            </span>
          </Reorder.Item>
        ))}
      </Reorder.Group>

      <Button variant="primary" size="lg" className="mt-4 w-full" onClick={handleConfirm}>
        Xác nhận
      </Button>
    </div>
  );
}
