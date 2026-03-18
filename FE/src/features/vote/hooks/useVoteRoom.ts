'use client';

import { useCallback, useEffect, useState } from 'react';
import { apiClient } from '@shared/lib/api-client';
import { useVoteStore } from '../stores/vote-store';
import { useSocket } from './useSocket';
import type { VoteRoom, VoteResult, VoteType } from '../types';

interface UseVoteRoomReturn {
  room: VoteRoom | null;
  status: VoteRoom['status'];
  participants: VoteRoom['participants'];
  results: VoteResult[];
  countdown: number;
  isLoading: boolean;
  error: string | null;
  joinRoom: (playerName: string) => Promise<void>;
  startVoting: () => void;
  submitVote: (data: unknown) => void;
  playAgain: () => void;
  changeVoteType: (type: VoteType) => void;
}

/**
 * Generate mock VoteResult[] from room dishes, scored by vote data.
 * Works for all 3 vote types: tournament, swipe, ranking.
 */
function generateMockResults(
  room: VoteRoom,
  voteData: Record<string, unknown>,
): VoteResult[] {
  const dishes = room.dishes;
  if (!dishes.length) return [];

  // Tournament: { winnerId }
  if (voteData.winnerId) {
    const winnerId = String(voteData.winnerId);
    return dishes
      .map((dish, i) => ({
        dishId: dish.id,
        dish,
        score: dish.id === winnerId ? 100 : Math.max(0, 80 - i * 10),
        rank: 0,
      }))
      .sort((a, b) => b.score - a.score)
      .map((r, i) => ({ ...r, rank: i + 1 }));
  }

  // Swipe: { liked: string[] }
  if (Array.isArray(voteData.liked)) {
    const likedIds = voteData.liked as string[];
    return dishes
      .map((dish) => ({
        dishId: dish.id,
        dish,
        score: likedIds.includes(dish.id) ? 50 + Math.floor(Math.random() * 50) : Math.floor(Math.random() * 30),
        rank: 0,
      }))
      .sort((a, b) => b.score - a.score)
      .map((r, i) => ({ ...r, rank: i + 1 }));
  }

  // Ranking: { ranking: { dishId, rank }[] }
  if (Array.isArray(voteData.ranking)) {
    const ranking = voteData.ranking as { dishId: string; rank: number }[];
    return ranking
      .map(({ dishId, rank }) => {
        const dish = dishes.find((d) => d.id === dishId);
        return dish
          ? { dishId, dish, score: Math.max(1, dishes.length - rank + 1) * 10, rank }
          : null;
      })
      .filter((r): r is VoteResult => r !== null);
  }

  // Fallback: random scores
  return dishes
    .map((dish) => ({
      dishId: dish.id,
      dish,
      score: Math.floor(Math.random() * 100),
      rank: 0,
    }))
    .sort((a, b) => b.score - a.score)
    .map((r, i) => ({ ...r, rank: i + 1 }));
}

export function useVoteRoom(roomId: string): UseVoteRoomReturn {
  const [isLoading, setIsLoading] = useState(true);
  const { emit } = useSocket(roomId);

  const room = useVoteStore((s) => s.room);
  const status = useVoteStore((s) => s.status);
  const results = useVoteStore((s) => s.results);
  const countdown = useVoteStore((s) => s.countdown);
  const error = useVoteStore((s) => s.error);
  const setRoom = useVoteStore((s) => s.setRoom);
  const setError = useVoteStore((s) => s.setError);
  const setResults = useVoteStore((s) => s.setResults);
  const storeStartVoting = useVoteStore((s) => s.startVoting);
  const storeSubmitVote = useVoteStore((s) => s.submitVote);
  const reset = useVoteStore((s) => s.reset);

  const participants = room?.participants ?? [];

  // Fetch room data via REST (falls back to mock)
  useEffect(() => {
    let cancelled = false;

    async function fetchRoom() {
      try {
        setIsLoading(true);
        const data = await apiClient.get<VoteRoom>(`/api/rooms/${roomId}`);
        if (!cancelled) {
          setRoom(data);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Không thể tải phòng');
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    void fetchRoom();
    return () => {
      cancelled = true;
    };
  }, [roomId, setRoom, setError]);

  const joinRoomAction = useCallback(
    async (playerName: string) => {
      try {
        await apiClient.post(`/api/rooms/${roomId}/join`, {
          name: playerName,
        });
        emit('room:join', { roomId, name: playerName });
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Không thể tham gia phòng');
      }
    },
    [roomId, emit, setError],
  );

  const startVotingAction = useCallback(() => {
    storeStartVoting();
    emit('voting:start', { roomId });
  }, [roomId, emit, storeStartVoting]);

  const submitVoteAction = useCallback(
    (data: unknown) => {
      storeSubmitVote(room?.hostName ?? '');
      emit('vote:submit', { roomId, ...(data as Record<string, unknown>) });

      // In mock mode (no WS), generate results locally after a short delay
      if (room) {
        setTimeout(() => {
          const mockResults = generateMockResults(room, (data ?? {}) as Record<string, unknown>);
          setResults(mockResults);
        }, 800);
      }
    },
    [roomId, emit, storeSubmitVote, room, setResults],
  );

  const playAgainAction = useCallback(() => {
    reset();
    // Re-fetch the room to go back to waiting state
    async function refetch() {
      try {
        const data = await apiClient.get<VoteRoom>(`/api/rooms/${roomId}`);
        setRoom(data);
      } catch {
        // Ignore
      }
    }
    void refetch();
  }, [roomId, reset, setRoom]);

  const changeVoteTypeAction = useCallback(
    (type: VoteType) => {
      if (!room) return;
      setRoom({ ...room, voteType: type });
    },
    [room, setRoom],
  );

  return {
    room,
    status,
    participants,
    results,
    countdown,
    isLoading,
    error,
    joinRoom: joinRoomAction,
    startVoting: startVotingAction,
    submitVote: submitVoteAction,
    playAgain: playAgainAction,
    changeVoteType: changeVoteTypeAction,
  };
}
