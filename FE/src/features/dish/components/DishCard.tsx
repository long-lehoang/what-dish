'use client';

import Image from 'next/image';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import type { Dish } from '../types';
import { DishBadge } from './DishBadge';

type DishCardVariant = 'default' | 'compact' | 'overlay';

interface DishCardProps {
  dish: Dish;
  variant?: DishCardVariant;
  onClick?: () => void;
  className?: string;
}

function ImagePlaceholder({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        'flex items-center justify-center bg-gradient-to-br from-primary/20 to-accent/20',
        className,
      )}
      aria-hidden="true"
    >
      <span className="text-4xl">{'\uD83C\uDF5C'}</span>
    </div>
  );
}

function DishImage({ dish, className }: { dish: Dish; className?: string }) {
  const src = dish.thumbnail ?? dish.imageUrl;

  if (!src) {
    return <ImagePlaceholder className={className} />;
  }

  return (
    <div className={cn('relative overflow-hidden', className)}>
      <Image
        src={src}
        alt={dish.name}
        fill
        sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
        className="object-cover transition-transform duration-300 group-hover:scale-105"
      />
    </div>
  );
}

function DefaultCard({ dish }: { dish: Dish }) {
  const totalTime = (dish.prepTime ?? 0) + (dish.cookTime ?? 0) || null;

  return (
    <>
      <DishImage dish={dish} className="aspect-[4/3] w-full" />
      <div className="flex flex-col gap-2 p-3">
        <h3 className="line-clamp-2 text-sm font-semibold text-gray-900 dark:text-gray-100">
          {dish.name}
        </h3>
        <div className="flex flex-wrap gap-1.5">
          <DishBadge type="time" value={totalTime} />
          <DishBadge type="difficulty" value={dish.difficulty} />
          <DishBadge type="spice" value={dish.spiceLevel} />
        </div>
      </div>
    </>
  );
}

function CompactCard({ dish }: { dish: Dish }) {
  const totalTime = (dish.prepTime ?? 0) + (dish.cookTime ?? 0) || null;

  return (
    <div className="flex items-center gap-3 p-2">
      <DishImage dish={dish} className="h-16 w-16 shrink-0 rounded-lg" />
      <div className="flex min-w-0 flex-col gap-1">
        <h3 className="truncate text-sm font-semibold text-gray-900 dark:text-gray-100">
          {dish.name}
        </h3>
        <div className="flex gap-1.5">
          <DishBadge type="time" value={totalTime} />
          <DishBadge type="difficulty" value={dish.difficulty} />
        </div>
      </div>
    </div>
  );
}

function OverlayCard({ dish }: { dish: Dish }) {
  const totalTime = (dish.prepTime ?? 0) + (dish.cookTime ?? 0) || null;

  return (
    <div className="relative">
      <DishImage dish={dish} className="aspect-[4/3] w-full" />
      <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent" />
      <div className="absolute bottom-0 left-0 right-0 p-3">
        <h3 className="mb-1.5 line-clamp-2 text-sm font-semibold text-white">{dish.name}</h3>
        <div className="flex gap-1.5">
          <DishBadge type="time" value={totalTime} />
          <DishBadge type="difficulty" value={dish.difficulty} />
        </div>
      </div>
    </div>
  );
}

const VARIANT_MAP: Record<DishCardVariant, React.ComponentType<{ dish: Dish }>> = {
  default: DefaultCard,
  compact: CompactCard,
  overlay: OverlayCard,
};

export function DishCard({ dish, variant = 'default', onClick, className }: DishCardProps) {
  const CardContent = VARIANT_MAP[variant];

  return (
    <motion.div
      whileHover={{ scale: 1.02 }}
      transition={{ type: 'spring', stiffness: 300, damping: 20 }}
      className={cn(
        'group cursor-pointer overflow-hidden rounded-xl bg-white shadow-sm transition-shadow hover:shadow-md dark:bg-dark-card',
        className,
      )}
    >
      <Link
        href={`/dish/${dish.slug}`}
        onClick={onClick}
        className="block"
        aria-label={`Xem ${dish.name}`}
      >
        <CardContent dish={dish} />
      </Link>
    </motion.div>
  );
}
