import type { Dish } from '@features/dish';

export type VoteType = 'tournament' | 'swipe' | 'ranking';

export type RoomStatus = 'idle' | 'waiting' | 'voting' | 'finished';

export interface VoteRoom {
  id: string;
  code: string;
  hostName: string;
  voteType: VoteType;
  status: RoomStatus;
  timerSecs: number;
  dishes: Dish[];
  participants: Participant[];
  createdAt: string;
  expiresAt: string;
}

export interface Participant {
  name: string;
  avatarColor: string;
  hasVoted: boolean;
}

export interface VoteResult {
  dishId: string;
  dish: Dish;
  score: number;
  rank: number;
}

export interface TournamentRound {
  roundNumber: number;
  matchups: Matchup[];
}

export interface Matchup {
  dish1: Dish;
  dish2: Dish;
  winnerId?: string;
}

export type RoomEvent =
  | 'participant-joined'
  | 'vote-submitted'
  | 'round-complete'
  | 'voting-started'
  | 'voting-finished'
  | 'timer-tick';
