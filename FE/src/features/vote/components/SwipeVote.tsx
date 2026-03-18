'use client';

import { useCallback, useState } from 'react';
import { AnimatePresence, motion, useMotionValue, useTransform } from 'framer-motion';
import { DishCard } from '@features/dish';
import type { Dish } from '@features/dish';

interface SwipeVoteProps {
  dishes: Dish[];
  onComplete: (data: unknown) => void;
  timeRemaining?: number;
}

export function SwipeVote({ dishes, onComplete, timeRemaining }: SwipeVoteProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [liked, setLiked] = useState<string[]>([]);
  const x = useMotionValue(0);
  const rotate = useTransform(x, [-200, 200], [-15, 15]);
  const likeOpacity = useTransform(x, [0, 100], [0, 1]);
  const passOpacity = useTransform(x, [-100, 0], [1, 0]);

  const currentDish = dishes[currentIndex];
  const progress = dishes.length > 0 ? currentIndex / dishes.length : 0;

  const handleSwipeComplete = useCallback(
    (direction: 'left' | 'right') => {
      if (!currentDish) return;

      const newLiked = direction === 'right' ? [...liked, currentDish.id] : liked;

      if (currentIndex + 1 >= dishes.length) {
        onComplete({
          liked: direction === 'right' ? newLiked : liked,
          ranking: newLiked,
        });
      } else {
        setLiked(newLiked);
        setCurrentIndex((prev) => prev + 1);
      }
    },
    [currentDish, currentIndex, dishes.length, liked, onComplete],
  );

  if (!currentDish) return null;

  return (
    <div className="flex flex-col items-center gap-6 px-4 py-6">
      <div className="text-center">
        <p className="text-sm font-medium text-gray-500 dark:text-gray-400">
          {currentIndex + 1}/{dishes.length}
        </p>
        {timeRemaining !== undefined && (
          <p className="mt-1 text-xs text-gray-400">Còn {timeRemaining}s</p>
        )}
      </div>

      {/* Progress bar */}
      <div className="h-1.5 w-full max-w-xs overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700">
        <div
          className="h-full rounded-full bg-primary transition-all duration-300"
          style={{ width: `${progress * 100}%` }}
        />
      </div>

      {/* Card stack */}
      <div className="relative h-[400px] w-full max-w-[300px]">
        <AnimatePresence>
          <motion.div
            key={currentDish.id}
            className="absolute inset-0"
            style={{ x, rotate }}
            drag="x"
            dragConstraints={{ left: 0, right: 0 }}
            dragElastic={0.8}
            onDragEnd={(_, info) => {
              if (info.offset.x > 100) {
                handleSwipeComplete('right');
              } else if (info.offset.x < -100) {
                handleSwipeComplete('left');
              }
            }}
            initial={{ scale: 0.95, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
          >
            {/* Like/Pass overlays */}
            <motion.div
              className="pointer-events-none absolute inset-0 z-10 flex items-center justify-center rounded-2xl border-4 border-green-500"
              style={{ opacity: likeOpacity }}
            >
              <span className="rotate-[-12deg] rounded-lg border-4 border-green-500 px-4 py-2 text-3xl font-bold text-green-500">
                THÍCH
              </span>
            </motion.div>
            <motion.div
              className="pointer-events-none absolute inset-0 z-10 flex items-center justify-center rounded-2xl border-4 border-red-500"
              style={{ opacity: passOpacity }}
            >
              <span className="rotate-[12deg] rounded-lg border-4 border-red-500 px-4 py-2 text-3xl font-bold text-red-500">
                BỎ
              </span>
            </motion.div>

            <DishCard dish={currentDish} className="h-full cursor-grab active:cursor-grabbing" />
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Manual buttons */}
      <div className="flex items-center gap-8">
        <button
          onClick={() => handleSwipeComplete('left')}
          className="flex h-14 w-14 items-center justify-center rounded-full border-2 border-red-300 text-red-500 transition-colors hover:bg-red-50 dark:border-red-700 dark:hover:bg-red-900/20"
          aria-label="Bỏ qua"
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
          >
            <path d="M18 6L6 18M6 6l12 12" />
          </svg>
        </button>
        <button
          onClick={() => handleSwipeComplete('right')}
          className="flex h-14 w-14 items-center justify-center rounded-full border-2 border-green-300 text-green-500 transition-colors hover:bg-green-50 dark:border-green-700 dark:hover:bg-green-900/20"
          aria-label="Thích"
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
          >
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
          </svg>
        </button>
      </div>
    </div>
  );
}
