'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/shared/ui';
import { apiClient } from '@/shared/lib/api-client';
import type { VoteType, VoteRoom } from '@/features/vote';

const VOTE_TYPE_OPTIONS: { value: VoteType; label: string; description: string; icon: string }[] = [
  {
    value: 'tournament',
    label: 'Đấu loại',
    description: 'So tay đôi, chọn 1 trong 2',
    icon: '⚔️',
  },
  {
    value: 'swipe',
    label: 'Quẹt thẻ',
    description: 'Quẹt phải = thích, trái = bỏ',
    icon: '👆',
  },
  {
    value: 'ranking',
    label: 'Xếp hạng',
    description: 'Kéo thả xếp thứ tự yêu thích',
    icon: '📊',
  },
];

export default function VoteCreatePage() {
  const router = useRouter();
  const [hostName, setHostName] = useState('');
  const [voteType, setVoteType] = useState<VoteType>('tournament');
  const [timerSecs, setTimerSecs] = useState(60);
  const [isCreating, setIsCreating] = useState(false);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!hostName.trim()) return;

    setIsCreating(true);
    try {
      const room = await apiClient.post<VoteRoom>('/api/rooms', {
        hostName: hostName.trim(),
        voteType,
        timerSecs,
      });
      router.push(`/vote/${room.id}`);
    } catch {
      setIsCreating(false);
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center px-4">
      <form onSubmit={handleCreate} className="mx-auto w-full max-w-md space-y-6">
        <div className="text-center">
          <h1 className="font-heading text-2xl font-bold">Tạo phòng vote</h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Thiết lập phòng bình chọn cho nhóm
          </p>
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium">Tên của bạn</label>
          <input
            type="text"
            value={hostName}
            onChange={(e) => setHostName(e.target.value)}
            placeholder="VD: Minh"
            maxLength={20}
            className="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-gray-700 dark:bg-dark-card"
            required
          />
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium">Kiểu bình chọn</label>
          <div className="space-y-2">
            {VOTE_TYPE_OPTIONS.map((option) => (
              <button
                key={option.value}
                type="button"
                onClick={() => setVoteType(option.value)}
                className={`flex w-full items-center gap-3 rounded-xl border-2 p-4 text-left transition-colors ${
                  voteType === option.value
                    ? 'border-primary bg-primary/5'
                    : 'border-gray-200 hover:border-gray-300 dark:border-gray-700'
                }`}
              >
                <span className="text-2xl">{option.icon}</span>
                <div>
                  <div className="font-medium">{option.label}</div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    {option.description}
                  </div>
                </div>
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium">
            Thời gian bình chọn: {timerSecs}s
          </label>
          <input
            type="range"
            min={30}
            max={180}
            step={15}
            value={timerSecs}
            onChange={(e) => setTimerSecs(Number(e.target.value))}
            className="w-full accent-primary"
          />
          <div className="flex justify-between text-xs text-gray-400">
            <span>30s</span>
            <span>180s</span>
          </div>
        </div>

        <Button
          type="submit"
          variant="primary"
          size="lg"
          className="w-full"
          isLoading={isCreating}
          disabled={!hostName.trim()}
        >
          Tạo phòng
        </Button>
      </form>
    </main>
  );
}
