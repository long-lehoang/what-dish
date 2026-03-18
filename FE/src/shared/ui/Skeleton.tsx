import { cn } from '@shared/lib/utils';

interface SkeletonProps {
  className?: string;
}

export function Skeleton({ className }: SkeletonProps) {
  return (
    <div
      className={cn('animate-shimmer rounded-lg bg-gray-200 dark:bg-gray-700', className)}
      aria-hidden="true"
    />
  );
}
