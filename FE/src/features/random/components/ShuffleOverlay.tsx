'use client';

import { useEffect, useMemo } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import type { ShufflePhase } from '../types';
import { generateCardPositions, generateConvergePositions } from '../utils/shuffle-algorithm';
import { PHASE_DURATIONS } from '../utils/animation-config';
import { MiniCard } from './MiniCard';

const CARD_COUNT = 12;

interface ShuffleOverlayProps {
  phase: ShufflePhase;
  onPhaseComplete?: () => void;
  className?: string;
}

export function ShuffleOverlay({ phase, onPhaseComplete, className }: ShuffleOverlayProps) {
  const isVisible = phase === 'shuffle' || phase === 'converge';

  const scatterPositions = useMemo(() => generateCardPositions(CARD_COUNT, 320, 400), []);

  const stackPositions = useMemo(() => generateConvergePositions(CARD_COUNT), []);

  useEffect(() => {
    if (phase === 'converge' && onPhaseComplete) {
      const id = setTimeout(onPhaseComplete, PHASE_DURATIONS.converge);
      return () => clearTimeout(id);
    }
  }, [phase, onPhaseComplete]);

  return (
    <AnimatePresence>
      {isVisible && (
        <motion.div
          className={cn(
            'fixed inset-0 z-40 flex items-center justify-center bg-black/60 backdrop-blur-sm',
            className,
          )}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.2 }}
          aria-live="polite"
          aria-label="Đang xáo trộn món ăn"
        >
          <div className="relative h-[400px] w-[320px]">
            {Array.from({ length: CARD_COUNT }).map((_, i) => {
              const scatter = scatterPositions[i]!;
              const stack = stackPositions[i]!;

              const targetPos = phase === 'shuffle' ? scatter : stack;

              return (
                <motion.div
                  key={i}
                  className="absolute left-1/2 top-1/2"
                  initial={{
                    x: 0,
                    y: 0,
                    rotate: 0,
                    opacity: 0,
                  }}
                  animate={{
                    x: targetPos.x,
                    y: targetPos.y,
                    rotate: targetPos.rotation,
                    opacity: 1,
                  }}
                  transition={{
                    duration:
                      phase === 'shuffle'
                        ? PHASE_DURATIONS.shuffle / 1000
                        : PHASE_DURATIONS.converge / 1000,
                    ease: phase === 'shuffle' ? 'easeInOut' : [0.25, 0.46, 0.45, 0.94],
                    delay: i * 0.03,
                  }}
                  style={{
                    marginLeft: -30,
                    marginTop: -40,
                  }}
                >
                  <MiniCard custom={i} />
                </motion.div>
              );
            })}
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
