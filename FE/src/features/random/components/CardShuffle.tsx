'use client';

import { useCallback, useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import { REROLL_EXHAUSTED_MESSAGE } from '@shared/constants';
import type { RandomFilters } from '../types';
import { useShuffle } from '../hooks/useShuffle';
import { useRerollLimit } from '../hooks/useRerollLimit';
import { FilterBar } from './FilterBar';
import { ShuffleOverlay } from './ShuffleOverlay';
import { CardReveal } from './CardReveal';
import { DishPool } from './DishPool';

interface CardShuffleProps {
  className?: string;
}

export function CardShuffle({ className }: CardShuffleProps) {
  const [filters, setFilters] = useState<RandomFilters>({});
  const { phase, selectedDish, triggerShuffle, reset, isAnimating } = useShuffle();
  const { canReroll, incrementReroll, resetRerolls } = useRerollLimit();

  const handleShuffle = useCallback(() => {
    if (isAnimating) return;
    triggerShuffle(filters);
  }, [isAnimating, triggerShuffle, filters]);

  const handleReroll = useCallback(() => {
    if (!canReroll || isAnimating) return;
    incrementReroll();
    triggerShuffle(filters);
  }, [canReroll, isAnimating, incrementReroll, triggerShuffle, filters]);

  const handleNewSession = useCallback(() => {
    resetRerolls();
    reset();
  }, [resetRerolls, reset]);

  // Keyboard navigation: Enter/Space triggers shuffle
  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Enter' || event.key === ' ') {
        if (phase === 'idle') {
          event.preventDefault();
          handleShuffle();
        }
      }
    }

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [phase, handleShuffle]);

  return (
    <div className={cn('flex flex-col', className)}>
      {/* Hero section with CTA */}
      <section className="relative overflow-hidden px-4 pb-8 pt-8 md:pt-12">
        {/* Gradient bg */}
        <div className="from-primary/8 absolute inset-0 bg-gradient-to-b via-transparent to-transparent dark:from-primary/15" />
        <div className="bg-primary/8 absolute left-1/2 top-0 h-[300px] w-[400px] -translate-x-1/2 rounded-full blur-3xl dark:bg-primary/10" />

        <div className="relative mx-auto max-w-2xl text-center">
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
          >
            <motion.span
              className="mb-3 inline-block text-5xl"
              animate={{ rotate: [0, -6, 6, -3, 0] }}
              transition={{ duration: 2, repeat: Infinity, repeatDelay: 4, ease: 'easeInOut' }}
            >
              🎲
            </motion.span>
            <h1 className="mb-2 font-heading text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">
              <span className="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                Tối Nay
              </span>{' '}
              Ăn Gì?
            </h1>
            <p className="mb-6 text-sm text-gray-500 dark:text-gray-400">
              Chọn bộ lọc, bấm lật bài — để số phận quyết định!
            </p>
          </motion.div>

          {/* Shuffle CTA */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.15, type: 'spring', stiffness: 260, damping: 20 }}
          >
            <button
              onClick={handleShuffle}
              disabled={isAnimating}
              className={cn(
                'group relative inline-flex items-center gap-3 overflow-hidden rounded-full',
                'bg-gradient-to-r from-primary to-secondary',
                'px-10 py-4 font-heading text-lg font-bold text-white',
                'shadow-xl shadow-primary/20 transition-all duration-300',
                'hover:shadow-2xl hover:shadow-primary/30',
                'active:scale-95 disabled:opacity-60 disabled:shadow-none',
              )}
              aria-label="Bắt đầu chọn món ăn ngẫu nhiên"
            >
              <span className="relative z-10">Lật bài!</span>
              <motion.span
                className="relative z-10 text-2xl"
                animate={{ rotateY: [0, 180, 360] }}
                transition={{ duration: 2, repeat: Infinity, repeatDelay: 3 }}
              >
                🎴
              </motion.span>
              {/* Glow hover */}
              <span className="absolute inset-0 bg-gradient-to-r from-secondary to-primary opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
              {/* Pulse ring */}
              <span className="absolute inset-0 animate-ping rounded-full bg-primary/20 [animation-duration:2s]" />
            </button>
          </motion.div>
        </div>
      </section>

      {/* Filter bar */}
      <FilterBar filters={filters} onFilterChange={setFilters} className="px-4 pb-4" />

      {/* Dish pool grid */}
      <DishPool filters={filters} className="px-4 pb-8" />

      {/* Exhaustion message when out of rerolls */}
      <AnimatePresence>
        {phase === 'settle' && !canReroll && (
          <motion.div
            className="fixed bottom-8 left-1/2 z-40 -translate-x-1/2 rounded-2xl bg-white px-6 py-4 text-center shadow-2xl dark:bg-dark-card"
            initial={{ opacity: 0, y: 40, scale: 0.9 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 20 }}
            transition={{ type: 'spring', stiffness: 260, damping: 20 }}
          >
            <p className="mb-2 text-sm font-medium text-primary">{REROLL_EXHAUSTED_MESSAGE}</p>
            <Button variant="ghost" size="sm" onClick={handleNewSession}>
              Bắt đầu lại
            </Button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Overlay animations */}
      <ShuffleOverlay phase={phase} />
      <CardReveal dish={selectedDish} phase={phase} onReroll={handleReroll} canReroll={canReroll} />
    </div>
  );
}
