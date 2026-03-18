'use client';

import { CardShuffle } from '@/features/random';

export default function RandomPage() {
  return (
    <main className="mx-auto min-h-screen max-w-6xl pb-12">
      <CardShuffle />
    </main>
  );
}
