import type { Metadata } from 'next';
import { apiClient } from '@/shared/lib/api-client';
import type { Dish, Category } from '@/features/dish';
import { ExploreClient } from './ExploreClient';

export const metadata: Metadata = {
  title: 'Khám phá món ăn',
  description: 'Tìm kiếm và khám phá các món ăn ngon, dễ nấu tại nhà.',
};

async function getInitialDishes() {
  try {
    return await apiClient.getList<Dish>('/api/v1/recipes?page=1&pageSize=20');
  } catch {
    return null;
  }
}

async function getCategories(): Promise<Category[]> {
  try {
    const result = await apiClient.get<Category[]>('/api/v1/categories');
    return result;
  } catch {
    return [];
  }
}

export default async function ExplorePage() {
  const [initialData, categories] = await Promise.all([getInitialDishes(), getCategories()]);

  return (
    <main className="min-h-screen px-4 py-6">
      <div className="mx-auto max-w-5xl">
        <h1 className="mb-6 font-heading text-2xl font-bold md:text-3xl">Khám phá món ăn</h1>
        <ExploreClient
          initialDishes={initialData?.data ?? []}
          initialTotal={initialData?.pagination.total ?? 0}
          categories={categories}
        />
      </div>
    </main>
  );
}
