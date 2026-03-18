'use client';

import { useMemo } from 'react';
import Link from 'next/link';
import { AnimatePresence, motion } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import type { Dish } from '@features/dish/types';
import { DishCard } from '@features/dish';
import type { ShufflePhase } from '../types';
import { PHASE_DURATIONS } from '../utils/animation-config';

interface CardRevealProps {
  dish: Dish | null;
  phase: ShufflePhase;
  onReroll: () => void;
  canReroll: boolean;
  className?: string;
}

const CONFETTI_COLORS = [
  '#FF6B35',
  '#E63946',
  '#FFB703',
  '#2EC4B6',
  '#9B5DE5',
  '#F15BB5',
  '#00BBF9',
  '#00F5D4',
  '#FEE440',
  '#FF6F91',
  '#845EC2',
  '#FF9671',
];

function Confetti() {
  const particles = useMemo(
    () =>
      CONFETTI_COLORS.map((color, i) => ({
        color,
        left: `${8 + ((i * 7.5) % 85)}%`,
        delay: i * 0.1,
        size: 6 + (i % 3) * 2,
      })),
    [],
  );

  return (
    <div className="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
      {particles.map((p, i) => (
        <div
          key={i}
          className="animate-confetti absolute top-0"
          style={{
            left: p.left,
            width: p.size,
            height: p.size,
            backgroundColor: p.color,
            borderRadius: i % 2 === 0 ? '50%' : '2px',
            animationDelay: `${p.delay}s`,
          }}
        />
      ))}
    </div>
  );
}

function CardBack() {
  return (
    <div
      className={cn(
        'flex h-full w-full items-center justify-center rounded-2xl',
        'bg-gradient-to-br from-primary to-secondary',
        'shadow-lg',
      )}
    >
      <div className="flex flex-col items-center gap-2 text-white">
        <span className="text-4xl">{'\uD83C\uDF5C'}</span>
        <span className="text-sm font-medium opacity-80">Tối Nay Ăn Gì?</span>
      </div>
    </div>
  );
}

export function CardReveal({ dish, phase, onReroll, canReroll, className }: CardRevealProps) {
  const isVisible = phase === 'select' || phase === 'reveal' || phase === 'settle';
  const isFlipped = phase === 'reveal' || phase === 'settle';
  const showConfetti = phase === 'reveal';
  const showActions = phase === 'settle';

  return (
    <AnimatePresence>
      {isVisible && (
        <motion.div
          className={cn(
            'fixed inset-0 z-50 flex flex-col items-center justify-center bg-black/60 backdrop-blur-sm',
            className,
          )}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.2 }}
        >
          {showConfetti && <Confetti />}

          {/* 3D flip card container */}
          <motion.div
            className={cn(
              'relative',
              (phase === 'reveal' || phase === 'settle') && 'animate-glow-pulse',
            )}
            style={{ perspective: 1000 }}
            initial={{ scale: 0.8, opacity: 0 }}
            animate={{
              scale: phase === 'select' ? 1.3 : 1.1,
              opacity: 1,
            }}
            transition={{
              duration:
                phase === 'select' ? PHASE_DURATIONS.select / 1000 : PHASE_DURATIONS.reveal / 1000,
              ease: [0.34, 1.56, 0.64, 1],
            }}
          >
            <motion.div
              className="relative h-[280px] w-[200px]"
              style={{ transformStyle: 'preserve-3d' }}
              animate={{ rotateY: isFlipped ? 180 : 0 }}
              transition={{ duration: 0.6, ease: [0.34, 1.56, 0.64, 1] }}
            >
              {/* Front face: card back */}
              <div className="absolute inset-0" style={{ backfaceVisibility: 'hidden' }}>
                <CardBack />
              </div>

              {/* Back face: dish info */}
              <div
                className="absolute inset-0 overflow-hidden rounded-2xl bg-white shadow-lg dark:bg-dark-card"
                style={{
                  backfaceVisibility: 'hidden',
                  transform: 'rotateY(180deg)',
                }}
              >
                {dish ? (
                  <DishCard dish={dish} variant="overlay" className="h-full shadow-none" />
                ) : (
                  <div className="flex h-full items-center justify-center">
                    <span className="text-gray-400">Đang chọn...</span>
                  </div>
                )}
              </div>
            </motion.div>
          </motion.div>

          {/* CTA Buttons */}
          <AnimatePresence>
            {showActions && dish && (
              <motion.div
                className="mt-6 flex gap-3"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.3, delay: 0.2 }}
              >
                <Link href={`/dish/${dish.slug}`}>
                  <Button variant="primary" size="lg">
                    Xem công thức
                  </Button>
                </Link>
                <Button
                  variant="outline"
                  size="lg"
                  onClick={onReroll}
                  disabled={!canReroll}
                  aria-label={canReroll ? 'Chọn lại món khác' : 'Đã hết lượt chọn lại'}
                >
                  Chọn lại
                </Button>
              </motion.div>
            )}
          </AnimatePresence>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
