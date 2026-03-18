'use client';

import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { cn } from '@shared/lib/utils';
import type { DishDetail } from '@features/dish';

interface RecipeHeroProps {
  dish: DishDetail;
  className?: string;
}

export function RecipeHero({ dish, className }: RecipeHeroProps) {
  const router = useRouter();
  const categoryLabel = dish.dishType?.name;

  return (
    <section className={cn('relative aspect-[16/9] w-full', className)}>
      {dish.imageUrl ? (
        <Image
          src={dish.imageUrl}
          alt={dish.name}
          fill
          priority
          sizes="100vw"
          className="object-cover"
        />
      ) : (
        <div className="h-full w-full bg-gray-200 dark:bg-gray-700" />
      )}

      <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />

      <button
        onClick={() => router.back()}
        className="absolute left-4 top-4 flex h-10 w-10 items-center justify-center rounded-full bg-black/30 text-white backdrop-blur-sm transition-colors hover:bg-black/50"
        aria-label="Quay lại"
      >
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
      </button>

      <div className="absolute bottom-0 left-0 right-0 p-6">
        {categoryLabel && (
          <span className="mb-2 inline-block rounded-full bg-primary/90 px-3 py-1 text-xs font-medium text-white">
            {categoryLabel}
          </span>
        )}
        <h1 className="text-2xl font-bold text-white md:text-3xl">{dish.name}</h1>
      </div>
    </section>
  );
}
