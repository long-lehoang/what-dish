'use client';

import { Badge } from '@shared/ui';
import { cn } from '@shared/lib/utils';

type DishBadgeType = 'time' | 'cost' | 'difficulty' | 'servings';

interface DishBadgeProps {
  type: DishBadgeType;
  value: number | string | null;
  className?: string;
}

function formatTime(minutes: number | string | null): string | null {
  if (minutes === null) return null;
  const numMinutes = typeof minutes === 'string' ? parseInt(minutes, 10) : minutes;
  if (isNaN(numMinutes)) return null;
  if (numMinutes < 60) return `${numMinutes} phút`;
  const hours = Math.floor(numMinutes / 60);
  const remaining = numMinutes % 60;
  return remaining === 0 ? `${hours} giờ` : `${hours}g${remaining}p`;
}

function formatCost(value: number | string | null): string | null {
  if (value === null) return null;
  const num = typeof value === 'string' ? parseInt(value, 10) : value;
  if (isNaN(num)) return null;
  return `~${Math.round(num / 1000)}k`;
}

function formatDifficulty(value: number | string | null): string | null {
  if (value === null) return null;
  if (typeof value === 'string') {
    switch (value) {
      case 'EASY':
        return 'Dễ';
      case 'MEDIUM':
        return 'TB';
      case 'HARD':
        return 'Khó';
      default:
        return value;
    }
  }
  if (value <= 2) return 'Dễ';
  if (value <= 3) return 'TB';
  return 'Khó';
}

function formatServings(value: number | string | null): string | null {
  if (value === null) return null;
  const num = typeof value === 'string' ? parseInt(value, 10) : value;
  if (isNaN(num) || num <= 0) return null;
  return `${num} phần`;
}

const BADGE_CONFIG: Record<
  DishBadgeType,
  {
    icon: string;
    format: (v: number | string | null) => string | null;
    variant: 'default' | 'success' | 'warning' | 'info';
  }
> = {
  time: { icon: '\u23F1', format: formatTime, variant: 'info' },
  cost: { icon: '\uD83D\uDCB0', format: formatCost, variant: 'success' },
  difficulty: {
    icon: '\uD83D\uDC68\u200D\uD83C\uDF73',
    format: formatDifficulty,
    variant: 'warning',
  },
  servings: { icon: '\uD83C\uDF7D\uFE0F', format: formatServings, variant: 'default' },
};

export function DishBadge({ type, value, className }: DishBadgeProps) {
  const config = BADGE_CONFIG[type];
  const formatted = config.format(value);

  if (formatted === null) return null;

  return (
    <Badge variant={config.variant} className={cn('gap-1', className)}>
      {config.icon && <span aria-hidden="true">{config.icon}</span>}
      <span>{formatted}</span>
    </Badge>
  );
}
