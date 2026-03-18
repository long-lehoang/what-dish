'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { Button } from '@/shared/ui';

export default function VoteLandingPage() {
  const router = useRouter();
  const [roomCode, setRoomCode] = useState('');

  function handleJoin(e: React.FormEvent) {
    e.preventDefault();
    const code = roomCode.trim().toUpperCase();
    if (code.length >= 4) {
      router.push(`/vote/${code}`);
    }
  }

  return (
    <main className="flex min-h-screen flex-col items-center justify-center px-4">
      <div className="mx-auto w-full max-w-md text-center">
        <div className="mb-6 text-5xl">🗳️</div>
        <h1 className="mb-2 font-heading text-3xl font-bold">Vote Nhóm</h1>
        <p className="mb-10 text-gray-600 dark:text-gray-400">Cùng bạn bè chọn món ăn tối nay!</p>

        <Link href="/vote/create">
          <Button variant="primary" size="lg" className="mb-8 w-full">
            Tạo phòng mới
          </Button>
        </Link>

        <div className="relative mb-8">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-200 dark:border-gray-700" />
          </div>
          <div className="relative flex justify-center">
            <span className="bg-background px-4 text-sm text-gray-500 dark:bg-dark-bg">hoặc</span>
          </div>
        </div>

        <form onSubmit={handleJoin}>
          <label className="mb-2 block text-left text-sm font-medium">Nhập mã phòng</label>
          <div className="flex gap-3">
            <input
              type="text"
              value={roomCode}
              onChange={(e) => setRoomCode(e.target.value.toUpperCase())}
              placeholder="VD: ABC123"
              maxLength={6}
              className="flex-1 rounded-xl border border-gray-200 bg-white px-4 py-3 text-center font-mono text-lg uppercase tracking-widest focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-gray-700 dark:bg-dark-card"
            />
            <Button type="submit" variant="secondary" disabled={roomCode.trim().length < 4}>
              Tham gia
            </Button>
          </div>
        </form>

        <Link
          href="/"
          className="mt-8 inline-block text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400"
        >
          ← Về trang chủ
        </Link>
      </div>
    </main>
  );
}
