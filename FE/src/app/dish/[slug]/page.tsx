import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import Link from 'next/link';
import { RecipeHero, RecipeInfo, IngredientList, StepCard } from '@/features/recipe';
import { apiClient } from '@/shared/lib/api-client';
import type { DishDetail } from '@/features/dish';

export const revalidate = 3600;

interface RecipePageProps {
  params: Promise<{ slug: string }>;
}

async function getDish(slug: string): Promise<DishDetail | null> {
  try {
    return await apiClient.get<DishDetail>(`/api/dishes/${slug}`);
  } catch {
    return null;
  }
}

export async function generateMetadata({ params }: RecipePageProps): Promise<Metadata> {
  const { slug } = await params;
  const dish = await getDish(slug);

  if (!dish) {
    return { title: 'Không tìm thấy món' };
  }

  return {
    title: dish.name,
    description: dish.description || `Công thức nấu ${dish.name} chi tiết, dễ làm tại nhà.`,
    openGraph: {
      title: dish.name,
      description: dish.description || `Công thức nấu ${dish.name}`,
      images: dish.imageUrl ? [{ url: dish.imageUrl }] : [],
    },
  };
}

export default async function RecipePage({ params }: RecipePageProps) {
  const { slug } = await params;
  const dish = await getDish(slug);

  if (!dish) {
    notFound();
  }

  return (
    <main className="min-h-screen pb-20">
      <RecipeHero dish={dish} />

      <div className="mx-auto max-w-2xl px-4">
        <RecipeInfo dish={dish} />

        <section className="mt-8">
          <h2 className="mb-4 font-heading text-xl font-bold">Nguyên liệu</h2>
          <IngredientList ingredients={dish.ingredients} originalServings={dish.servings} />
        </section>

        <section className="mt-8">
          <h2 className="mb-4 font-heading text-xl font-bold">Cách làm</h2>
          <div className="space-y-4">
            {dish.steps
              .sort((a, b) => a.stepNumber - b.stepNumber)
              .map((step) => (
                <StepCard key={step.id} step={step} />
              ))}
          </div>
        </section>

        <div className="mt-10 text-center">
          <Link
            href={`/dish/${slug}/cook`}
            className="inline-block rounded-full bg-primary px-8 py-3 font-heading font-bold text-white shadow-lg transition-transform hover:scale-105"
          >
            Bật Cook Mode 👨‍🍳
          </Link>
        </div>
      </div>
    </main>
  );
}
