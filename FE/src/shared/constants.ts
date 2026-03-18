export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

export const WS_BASE_URL = process.env.NEXT_PUBLIC_WS_URL ?? 'ws://localhost:8080';

export const MAX_REROLLS = 3;

export const HISTORY_DAYS = 7;

export const REROLL_EXHAUSTED_MESSAGE = 'Thôi ăn cái này đi, đừng kén nữa 😄';

export const VOTE_TYPES = {
  TOURNAMENT: 'tournament',
  SWIPE: 'swipe',
  RANKING: 'ranking',
} as const;

export type VoteType = (typeof VOTE_TYPES)[keyof typeof VOTE_TYPES];

export const ROOM_STATUSES = {
  IDLE: 'idle',
  WAITING: 'waiting',
  VOTING: 'voting',
  FINISHED: 'finished',
} as const;

export type RoomStatus = (typeof ROOM_STATUSES)[keyof typeof ROOM_STATUSES];

export const DEFAULT_SERVINGS = 2;
