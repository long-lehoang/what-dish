'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { CookMode } from '@/features/recipe';
import { apiClient } from '@/shared/lib/api-client';
import type { DishDetail } from '@/features/dish';
import { Skeleton } from '@/shared/ui';

export default function CookModePage() {
  const params = useParams<{ slug: string }>();
  const router = useRouter();
  const [dish, setDish] = useState<DishDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function fetchDish() {
      try {
        const data = await apiClient.get<DishDetail>(`/api/dishes/${params.slug}`);
        setDish(data);
      } catch {
        router.push(`/dish/${params.slug}`);
      } finally {
        setIsLoading(false);
      }
    }
    fetchDish();
  }, [params.slug, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-dark-bg">
        <Skeleton className="h-32 w-32 rounded-full" />
      </div>
    );
  }

  if (!dish) return null;

  const sortedSteps = [...dish.steps].sort((a, b) => a.stepNumber - b.stepNumber);

  return (
    <CookMode.Root steps={sortedSteps} onExit={() => router.push(`/dish/${params.slug}`)}>
      <CookMode.Step />
      <CookMode.Controls />
    </CookMode.Root>
  );
}
