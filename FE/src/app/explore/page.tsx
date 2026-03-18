import type { Metadata } from 'next';
import { apiClient } from '@/shared/lib/api-client';
import type { DishListResponse } from '@/features/dish';
import { ExploreClient } from './ExploreClient';

export const metadata: Metadata = {
  title: 'Khám phá món ăn',
  description: 'Tìm kiếm và khám phá các món ăn ngon, dễ nấu tại nhà.',
};

async function getInitialDishes(): Promise<DishListResponse | null> {
  try {
    return await apiClient.get<DishListResponse>('/api/dishes?page=1&pageSize=20');
  } catch {
    return null;
  }
}

export default async function ExplorePage() {
  const initialData = await getInitialDishes();

  return (
    <main className="min-h-screen px-4 py-6">
      <div className="mx-auto max-w-5xl">
        <h1 className="mb-6 font-heading text-2xl font-bold md:text-3xl">Khám phá món ăn</h1>
        <ExploreClient
          initialDishes={initialData?.dishes ?? []}
          initialTotal={initialData?.total ?? 0}
        />
      </div>
    </main>
  );
}
