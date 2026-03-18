'use client';

import { useState } from 'react';
import Image from 'next/image';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import type { Step } from '@features/dish';

interface StepCardProps {
  step: Step;
  isActive?: boolean;
  className?: string;
}

export function StepCard({ step, isActive, className }: StepCardProps) {
  const [showTip, setShowTip] = useState(false);
  const timerMinutes = step.timerSecs !== null ? Math.ceil(step.timerSecs / 60) : null;

  return (
    <div
      className={cn(
        'rounded-2xl border bg-white p-4 shadow-sm transition-colors dark:bg-dark-card',
        isActive
          ? 'border-primary/50 ring-2 ring-primary/20'
          : 'border-gray-200 dark:border-gray-700',
        className,
      )}
    >
      <div className="mb-3 flex items-start gap-3">
        <span
          className={cn(
            'flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-bold text-white',
            isActive ? 'bg-primary' : 'bg-gray-400 dark:bg-gray-600',
          )}
        >
          {step.stepNumber}
        </span>
        <p className="flex-1 text-sm leading-relaxed text-gray-800 dark:text-gray-200">
          {step.instruction}
        </p>
      </div>

      {step.imageUrl && (
        <div className="relative mb-3 aspect-video overflow-hidden rounded-xl">
          <Image
            src={step.imageUrl}
            alt={`Bước ${step.stepNumber}`}
            fill
            sizes="(max-width: 768px) 100vw, 50vw"
            className="object-cover"
          />
        </div>
      )}

      {timerMinutes !== null && (
        <button
          className="mb-2 inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-3 py-1.5 text-xs font-medium text-primary transition-colors hover:bg-primary/20"
          aria-label={`Hẹn giờ ${timerMinutes} phút`}
        >
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
          >
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
          </svg>
          Hẹn giờ {timerMinutes} phút
        </button>
      )}

      {step.tip !== null && (
        <div>
          <button
            onClick={() => setShowTip((prev) => !prev)}
            className="text-xs font-medium text-amber-600 transition-colors hover:text-amber-700 dark:text-amber-400"
            aria-expanded={showTip}
          >
            {showTip ? '▲ Ẩn mẹo' : '▼ Xem mẹo'}
          </button>
          <AnimatePresence>
            {showTip && (
              <motion.div
                initial={{ height: 0, opacity: 0 }}
                animate={{ height: 'auto', opacity: 1 }}
                exit={{ height: 0, opacity: 0 }}
                transition={{ duration: 0.2 }}
                className="overflow-hidden"
              >
                <p className="mt-2 rounded-lg bg-amber-50 p-3 text-xs text-amber-800 dark:bg-amber-900/20 dark:text-amber-300">
                  <span aria-hidden="true">&#x1F4A1; </span>
                  Mẹo: {step.tip}
                </p>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      )}
    </div>
  );
}
