import { Skeleton } from '@/shared/ui';

export default function Loading() {
  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="space-y-4 text-center">
        <Skeleton className="mx-auto h-16 w-16 rounded-full" />
        <Skeleton className="mx-auto h-6 w-48" />
        <Skeleton className="mx-auto h-4 w-32" />
      </div>
    </main>
  );
}
