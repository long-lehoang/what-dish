'use client';

import Image from 'next/image';
import { cn } from '@shared/lib/utils';
import type { Step } from '@features/dish';

interface StepCardProps {
  step: Step;
  isActive?: boolean;
  className?: string;
}

export function StepCard({ step, isActive, className }: StepCardProps) {
  const timerMinutes = step.duration ? Math.ceil(step.duration / 60) : null;

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
        <div className="flex-1">
          {step.title && (
            <h4 className="mb-1 text-sm font-semibold text-gray-900 dark:text-gray-100">
              {step.title}
            </h4>
          )}
          <p className="text-sm leading-relaxed text-gray-800 dark:text-gray-200">
            {step.description}
          </p>
        </div>
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
          className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-3 py-1.5 text-xs font-medium text-primary transition-colors hover:bg-primary/20"
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
    </div>
  );
}
