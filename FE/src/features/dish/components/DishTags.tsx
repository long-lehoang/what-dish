'use client';

import { Badge } from '@shared/ui';
import { cn } from '@shared/lib/utils';
import type { Tag } from '../types';

interface DishTagsProps {
  tags: Tag[];
  maxVisible?: number;
  className?: string;
}

export function DishTags({ tags, maxVisible = 3, className }: DishTagsProps) {
  const visible = tags.slice(0, maxVisible);
  const overflowCount = tags.length - maxVisible;

  if (tags.length === 0) return null;

  return (
    <div className={cn('scrollbar-none flex items-center gap-1.5 overflow-x-auto', className)}>
      {visible.map((tag) => (
        <Badge key={tag.id} variant="default" className="shrink-0 whitespace-nowrap">
          {tag.name}
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
