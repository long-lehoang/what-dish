import { create } from 'zustand';
import type { VoteRoom, RoomStatus, Participant, VoteResult, TournamentRound } from '../types';

const VALID_TRANSITIONS: Record<RoomStatus, RoomStatus[]> = {
  idle: ['waiting'],
  waiting: ['voting'],
  voting: ['finished'],
  finished: ['idle'],
};

function canTransition(from: RoomStatus, to: RoomStatus): boolean {
  return VALID_TRANSITIONS[from].includes(to);
}

interface VoteState {
  room: VoteRoom | null;
  status: RoomStatus;
  results: VoteResult[];
  currentRound: TournamentRound | null;
  countdown: number;
  error: string | null;

  // Actions
  setRoom: (room: VoteRoom) => void;
  addParticipant: (participant: Participant) => void;
  removeParticipant: (name: string) => void;
  startVoting: () => void;
  submitVote: (participantName: string) => void;
  receiveVoteUpdate: (participantName: string) => void;
  setResults: (results: VoteResult[]) => void;
  setCountdown: (seconds: number) => void;
  setCurrentRound: (round: TournamentRound | null) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

const initialState = {
  room: null,
  status: 'idle' as RoomStatus,
  results: [],
  currentRound: null,
  countdown: 0,
  error: null,
};

export const useVoteStore = create<VoteState>((set, get) => ({
  ...initialState,

  setRoom: (room) => {
    set({ room, status: room.status, error: null });
  },

  addParticipant: (participant) => {
    const { room } = get();
    if (!room) return;

    const exists = room.participants.some((p) => p.name === participant.name);
    if (exists) return;

    set({
      room: {
        ...room,
        participants: [...room.participants, participant],
      },
    });
  },

  removeParticipant: (name) => {
    const { room } = get();
    if (!room) return;

    set({
      room: {
        ...room,
        participants: room.participants.filter((p) => p.name !== name),
      },
    });
  },

  startVoting: () => {
    const { status, room } = get();
    if (!canTransition(status, 'voting')) {
      if (process.env.NODE_ENV === 'development') {
        // eslint-disable-next-line no-console
        console.warn(`Invalid vote store transition: ${status} -> voting`);
      }
      return;
    }

    set({
      status: 'voting',
      room: room ? { ...room, status: 'voting' } : null,
    });
  },

  submitVote: (participantName) => {
    const { room } = get();
    if (!room) return;

    set({
      room: {
        ...room,
        participants: room.participants.map((p) =>
          p.name === participantName ? { ...p, hasVoted: true } : p,
        ),
      },
    });
  },

  receiveVoteUpdate: (participantName) => {
    const { room } = get();
    if (!room) return;

    set({
      room: {
        ...room,
        participants: room.participants.map((p) =>
          p.name === participantName ? { ...p, hasVoted: true } : p,
        ),
      },
    });
  },

  setResults: (results) => {
    const { status, room } = get();
    if (!canTransition(status, 'finished')) {
      if (process.env.NODE_ENV === 'development') {
        // eslint-disable-next-line no-console
        console.warn(`Invalid vote store transition: ${status} -> finished`);
      }
      return;
    }

    set({
      results,
      status: 'finished',
      room: room ? { ...room, status: 'finished' } : null,
    });
  },

  setCountdown: (seconds) => {
    set({ countdown: seconds });
  },

  setCurrentRound: (round) => {
    set({ currentRound: round });
  },

  setError: (error) => {
    set({ error });
  },

  reset: () => {
    set(initialState);
  },
}));
