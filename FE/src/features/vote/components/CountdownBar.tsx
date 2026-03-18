'use client';

import { cn } from '@shared/lib/utils';

interface CountdownBarProps {
  total: number;
  remaining: number;
  className?: string;
}

export function CountdownBar({ total, remaining, className }: CountdownBarProps) {
  const fraction = total > 0 ? remaining / total : 0;
  const percentage = Math.max(0, Math.min(100, fraction * 100));

  const isFlashing = remaining <= 10 && remaining > 0;

  let barColor: string;
  if (fraction <= 0.25) {
    barColor = 'bg-red-500';
  } else if (fraction <= 0.5) {
    barColor = 'bg-yellow-500';
  } else {
    barColor = 'bg-green-500';
  }

  return (
    <div
      className={cn(
        'h-2 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700',
        className,
      )}
      role="progressbar"
      aria-valuenow={remaining}
      aria-valuemin={0}
      aria-valuemax={total}
      aria-label={`Còn lại ${remaining} giây`}
    >
      <div
        className={cn(
          'h-full rounded-full transition-all duration-1000 ease-linear',
          barColor,
          isFlashing && 'animate-pulse',
        )}
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
}
