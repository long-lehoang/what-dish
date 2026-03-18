'use client';

import { Badge } from '@shared/ui';
import { cn } from '@shared/lib/utils';

interface DishTagsProps {
  tags: string[];
  dietary?: string[];
  maxVisible?: number;
  className?: string;
}

export function DishTags({ tags, dietary = [], maxVisible = 3, className }: DishTagsProps) {
  const allItems = [
    ...dietary.map((d) => ({ label: d, isDietary: true })),
    ...tags.map((t) => ({ label: t, isDietary: false })),
  ];

  const visible = allItems.slice(0, maxVisible);
  const overflowCount = allItems.length - maxVisible;

  if (allItems.length === 0) return null;

  return (
    <div className={cn('scrollbar-none flex items-center gap-1.5 overflow-x-auto', className)}>
      {visible.map((item) => (
        <Badge
          key={item.label}
          variant={item.isDietary ? 'success' : 'default'}
          className="shrink-0 whitespace-nowrap"
        >
          {item.label}
        </Badge>
      ))}
      {overflowCount > 0 && (
        <Badge variant="default" className="shrink-0 whitespace-nowrap">
          +{overflowCount}
        </Badge>
      )}
    </div>
  );
}
