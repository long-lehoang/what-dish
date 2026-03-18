'use client';

import { useCallback, useState } from 'react';
import { motion } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import { ParticipantList } from './ParticipantList';
import type { Participant, VoteType } from '../types';

interface RoomLobbyProps {
  roomCode: string;
  hostName: string;
  voteType: VoteType;
  participants: Participant[];
  isHost: boolean;
  onStart: () => void;
  onChangeVoteType?: (type: VoteType) => void;
  className?: string;
}

const VOTE_TYPE_CONFIG: { value: VoteType; label: string; icon: string; desc: string }[] = [
  { value: 'tournament', label: 'Đấu loại', icon: '⚔️', desc: '1 vs 1, chọn món thắng' },
  { value: 'swipe', label: 'Vuốt chọn', icon: '👆', desc: 'Vuốt phải = thích, trái = bỏ' },
  { value: 'ranking', label: 'Xếp hạng', icon: '📊', desc: 'Kéo thả sắp xếp yêu thích' },
];

export function RoomLobby({
  roomCode,
  hostName,
  voteType,
  participants,
  isHost,
  onStart,
  onChangeVoteType,
  className,
}: RoomLobbyProps) {
  const [copied, setCopied] = useState(false);

  const copyCode = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(roomCode);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API not available
    }
  }, [roomCode]);

  const shareRoom = useCallback(async () => {
    const url = `${window.location.origin}/vote/${roomCode}`;
    if (navigator.share) {
      try {
        await navigator.share({
          title: 'Tối Nay Ăn Gì - Tham gia bỏ phiếu',
          text: `Tham gia phòng bỏ phiếu: ${roomCode}`,
          url,
        });
      } catch {
        // User cancelled or share not supported
      }
    } else {
      try {
        await navigator.clipboard.writeText(url);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch {
        // Clipboard API not available
      }
    }
  }, [roomCode]);

  return (
    <div className={cn('flex flex-col items-center gap-6 px-4 py-8', className)}>
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center"
      >
        <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Mã phòng</p>
        <button
          onClick={copyCode}
          className="mt-2 text-4xl font-bold tracking-[0.3em] text-gray-900 transition-colors hover:text-primary dark:text-gray-100"
          style={{ fontFamily: 'monospace' }}
          aria-label={`Mã phòng ${roomCode}. Nhấn để sao chép`}
        >
          {roomCode}
        </button>
        {copied && <p className="mt-1 text-xs text-green-600 dark:text-green-400">Đã sao chép!</p>}
      </motion.div>

      {/* Room info */}
      <div className="flex items-center gap-3">
        <span className="text-xs text-gray-500 dark:text-gray-400">Chủ phòng: {hostName}</span>
      </div>

      <Button variant="outline" size="sm" onClick={shareRoom}>
        Chia sẻ phòng
      </Button>

      {/* Vote type selector (host only) */}
      {isHost && onChangeVoteType && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.1 }}
          className="w-full max-w-sm"
        >
          <h3 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
            Kiểu bình chọn
          </h3>
          <div className="space-y-2">
            {VOTE_TYPE_CONFIG.map((opt) => (
              <button
                key={opt.value}
                type="button"
                onClick={() => onChangeVoteType(opt.value)}
                className={cn(
                  'flex w-full items-center gap-3 rounded-xl border-2 p-3 text-left transition-all',
                  voteType === opt.value
                    ? 'border-primary bg-primary/5 shadow-sm'
                    : 'border-gray-200 hover:border-gray-300 dark:border-gray-700 dark:hover:border-gray-600',
                )}
              >
                <span className="text-2xl">{opt.icon}</span>
                <div className="flex-1">
                  <div className="text-sm font-medium">{opt.label}</div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">{opt.desc}</div>
                </div>
                {voteType === opt.value && (
                  <span className="text-xs font-bold text-primary">✓</span>
                )}
              </button>
            ))}
          </div>
        </motion.div>
      )}

      {/* Participants */}
      <div className="w-full max-w-sm">
        <h3 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
          Người tham gia ({participants.length})
        </h3>
        <ParticipantList participants={participants} />
      </div>

      {/* Start button */}
      {isHost && (
        <Button
          variant="primary"
          size="lg"
          className="w-full max-w-sm"
          onClick={onStart}
          disabled={participants.length < 2}
        >
          Bắt đầu bình chọn
        </Button>
      )}

      {!isHost && (
        <p className="text-sm text-gray-500 dark:text-gray-400">Đang chờ chủ phòng bắt đầu...</p>
      )}
    </div>
  );
}
