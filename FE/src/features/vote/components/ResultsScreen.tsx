'use client';

import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import { DishCard } from '@features/dish';
import type { VoteResult } from '../types';

interface ResultsScreenProps {
  results: VoteResult[];
  onPlayAgain?: () => void;
  className?: string;
}

export function ResultsScreen({ results, onPlayAgain, className }: ResultsScreenProps) {
  const [showConfetti, setShowConfetti] = useState(true);
  const winner = results[0];

  useEffect(() => {
    const timer = setTimeout(() => setShowConfetti(false), 4000);
    return () => clearTimeout(timer);
  }, []);

  if (!winner) return null;

  return (
    <div className={cn('flex flex-col items-center gap-6 px-4 py-8', className)}>
      {/* Confetti */}
      {showConfetti && (
        <div className="pointer-events-none fixed inset-0 z-50 overflow-hidden">
          {Array.from({ length: 40 }, (_, i) => (
            <motion.div
              key={i}
              className="absolute h-3 w-3 rounded-full"
              style={{
                backgroundColor: ['#FF6B35', '#E63946', '#FFB703', '#2EC4B6', '#9B5DE5'][i % 5],
                left: `${Math.random() * 100}%`,
              }}
              initial={{ y: -20, opacity: 1 }}
              animate={{
                y: '100vh',
                x: (Math.random() - 0.5) * 200,
                rotate: Math.random() * 720,
                opacity: 0,
              }}
              transition={{
                duration: 2 + Math.random() * 2,
                delay: Math.random() * 0.5,
                ease: 'easeOut',
              }}
            />
          ))}
        </div>
      )}

      {/* Winner announcement */}
      <motion.div
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ delay: 0.3, type: 'spring', stiffness: 200 }}
        className="text-center"
      >
        <p className="text-lg font-medium text-gray-500 dark:text-gray-400">Kết quả</p>
        <h2 className="mt-2 text-2xl font-bold text-gray-900 dark:text-gray-100 md:text-3xl">
          Tối nay ăn <span className="text-primary">{winner.dish.name}</span>!
        </h2>
      </motion.div>

      {/* Winner card */}
      <motion.div
        initial={{ y: 30, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: 0.5 }}
        className="w-full max-w-xs"
      >
        <DishCard dish={winner.dish} />
      </motion.div>

      {/* Rankings */}
      {results.length > 1 && (
        <div className="w-full max-w-sm">
          <h3 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
            Bảng xếp hạng
          </h3>
          <ol className="space-y-2">
            {results.map((result) => (
              <motion.li
                key={result.dishId}
                initial={{ x: -20, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                transition={{ delay: 0.6 + result.rank * 0.1 }}
                className={cn(
                  'flex items-center gap-3 rounded-xl px-4 py-3',
                  result.rank === 1
                    ? 'bg-amber-50 dark:bg-amber-900/20'
                    : 'bg-gray-50 dark:bg-gray-800/50',
                )}
              >
                <span
                  className={cn(
                    'flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-bold',
                    result.rank === 1
                      ? 'bg-amber-200 text-amber-800'
                      : 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-400',
                  )}
                >
                  {result.rank}
                </span>
                <span className="flex-1 text-sm font-medium text-gray-900 dark:text-gray-100">
                  {result.dish.name}
                </span>
                <span className="text-xs text-gray-500">{result.score} điểm</span>
              </motion.li>
            ))}
          </ol>
        </div>
      )}

      {/* CTAs */}
      <div className="flex w-full max-w-sm flex-col gap-3">
        {onPlayAgain && (
          <Button variant="primary" size="lg" className="w-full" onClick={onPlayAgain}>
            Chơi lại
          </Button>
        )}
        <Link href={`/dish/${winner.dish.slug}`} className="w-full">
          <Button variant="outline" size="lg" className="w-full">
            Xem công thức
          </Button>
        </Link>
      </div>
    </div>
  );
}
