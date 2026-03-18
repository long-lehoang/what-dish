import { describe, it, expect, beforeEach } from 'vitest';
import { useVoteStore } from '../vote-store';
import type { VoteRoom, Participant, VoteResult } from '../../types';

function createMockRoom(overrides?: Partial<VoteRoom>): VoteRoom {
  return {
    id: 'room-1',
    code: 'ABC123',
    hostName: 'Host',
    voteType: 'tournament',
    status: 'waiting',
    timerSecs: 60,
    dishes: [],
    participants: [],
    createdAt: '2024-01-01T00:00:00Z',
    expiresAt: '2024-01-01T01:00:00Z',
    ...overrides,
  };
}

function createMockParticipant(overrides?: Partial<Participant>): Participant {
  return {
    name: 'Player 1',
    avatarColor: '#FF6B35',
    hasVoted: false,
    ...overrides,
  };
}

describe('useVoteStore', () => {
  beforeEach(() => {
    useVoteStore.getState().reset();
  });

  describe('initial state', () => {
    it('starts with idle status and no room', () => {
      const state = useVoteStore.getState();
      expect(state.room).toBeNull();
      expect(state.status).toBe('idle');
      expect(state.results).toEqual([]);
      expect(state.currentRound).toBeNull();
      expect(state.countdown).toBe(0);
      expect(state.error).toBeNull();
    });
  });

  describe('setRoom', () => {
    it('sets room and syncs status', () => {
      const room = createMockRoom({ status: 'waiting' });
      useVoteStore.getState().setRoom(room);

      const state = useVoteStore.getState();
      expect(state.room).toEqual(room);
      expect(state.status).toBe('waiting');
      expect(state.error).toBeNull();
    });
  });

  describe('addParticipant', () => {
    it('adds a participant to the room', () => {
      useVoteStore.getState().setRoom(createMockRoom());
      const participant = createMockParticipant();

      useVoteStore.getState().addParticipant(participant);

      expect(useVoteStore.getState().room!.participants).toHaveLength(1);
      expect(useVoteStore.getState().room!.participants[0]).toEqual(participant);
    });

    it('does not add duplicate participants', () => {
      useVoteStore.getState().setRoom(createMockRoom());
      const participant = createMockParticipant();

      useVoteStore.getState().addParticipant(participant);
      useVoteStore.getState().addParticipant(participant);

      expect(useVoteStore.getState().room!.participants).toHaveLength(1);
    });

    it('does nothing when no room is set', () => {
      useVoteStore.getState().addParticipant(createMockParticipant());
      expect(useVoteStore.getState().room).toBeNull();
    });
  });

  describe('removeParticipant', () => {
    it('removes a participant by name', () => {
      const room = createMockRoom({
        participants: [
          createMockParticipant({ name: 'Alice' }),
          createMockParticipant({ name: 'Bob' }),
        ],
      });
      useVoteStore.getState().setRoom(room);

      useVoteStore.getState().removeParticipant('Alice');

      const participants = useVoteStore.getState().room!.participants;
      expect(participants).toHaveLength(1);
      expect(participants[0]!.name).toBe('Bob');
    });
  });

  describe('state machine transitions', () => {
    it('transitions from waiting → voting (valid)', () => {
      useVoteStore.getState().setRoom(createMockRoom({ status: 'waiting' }));

      useVoteStore.getState().startVoting();

      expect(useVoteStore.getState().status).toBe('voting');
      expect(useVoteStore.getState().room!.status).toBe('voting');
    });

    it('rejects idle → voting (invalid)', () => {
      useVoteStore.getState().setRoom(createMockRoom({ status: 'idle' }));

      useVoteStore.getState().startVoting();

      // Status should remain idle
      expect(useVoteStore.getState().status).toBe('idle');
    });

    it('transitions from voting → finished via setResults', () => {
      useVoteStore.getState().setRoom(createMockRoom({ status: 'waiting' }));
      useVoteStore.getState().startVoting();

      const results: VoteResult[] = [
        {
          dishId: 'dish-1',
          dish: {
            id: 'dish-1',
            name: 'Phở',
            slug: 'pho',
            status: 'PUBLISHED',
            viewCount: 0,
            favoriteCount: 0,
            createdAt: '2024-01-01',
            updatedAt: '2024-01-01',
          },
          score: 10,
          rank: 1,
        },
      ];
      useVoteStore.getState().setResults(results);

      expect(useVoteStore.getState().status).toBe('finished');
      expect(useVoteStore.getState().results).toEqual(results);
    });

    it('rejects setting results from idle status', () => {
      // idle → finished is not a valid transition
      useVoteStore.getState().setResults([]);

      expect(useVoteStore.getState().status).toBe('idle');
      expect(useVoteStore.getState().results).toEqual([]);
    });
  });

  describe('submitVote', () => {
    it('marks participant as having voted', () => {
      const room = createMockRoom({
        participants: [createMockParticipant({ name: 'Alice', hasVoted: false })],
      });
      useVoteStore.getState().setRoom(room);

      useVoteStore.getState().submitVote('Alice');

      expect(useVoteStore.getState().room!.participants[0]!.hasVoted).toBe(true);
    });

    it('does not affect other participants', () => {
      const room = createMockRoom({
        participants: [
          createMockParticipant({ name: 'Alice' }),
          createMockParticipant({ name: 'Bob' }),
        ],
      });
      useVoteStore.getState().setRoom(room);

      useVoteStore.getState().submitVote('Alice');

      expect(useVoteStore.getState().room!.participants[1]!.hasVoted).toBe(false);
    });
  });

  describe('receiveVoteUpdate', () => {
    it('marks participant as voted on remote update', () => {
      const room = createMockRoom({
        participants: [createMockParticipant({ name: 'Bob', hasVoted: false })],
      });
      useVoteStore.getState().setRoom(room);

      useVoteStore.getState().receiveVoteUpdate('Bob');

      expect(useVoteStore.getState().room!.participants[0]!.hasVoted).toBe(true);
    });
  });

  describe('setCountdown', () => {
    it('updates countdown value', () => {
      useVoteStore.getState().setCountdown(30);
      expect(useVoteStore.getState().countdown).toBe(30);
    });
  });

  describe('setError', () => {
    it('sets and clears error', () => {
      useVoteStore.getState().setError('Connection lost');
      expect(useVoteStore.getState().error).toBe('Connection lost');

      useVoteStore.getState().setError(null);
      expect(useVoteStore.getState().error).toBeNull();
    });
  });

  describe('reset', () => {
    it('resets to initial state', () => {
      useVoteStore.getState().setRoom(createMockRoom());
      useVoteStore.getState().setCountdown(30);
      useVoteStore.getState().setError('test');

      useVoteStore.getState().reset();

      const state = useVoteStore.getState();
      expect(state.room).toBeNull();
      expect(state.status).toBe('idle');
      expect(state.results).toEqual([]);
      expect(state.countdown).toBe(0);
      expect(state.error).toBeNull();
    });
  });
});
