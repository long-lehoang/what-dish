'use client';

import type { Dish } from '@features/dish';
import type { VoteType } from '../types';
import { Tournament } from './Tournament';
import { SwipeVote } from './SwipeVote';
import { RankVote } from './RankVote';

interface VoteStrategyProps {
  voteType: VoteType;
  dishes: Dish[];
  onVoteComplete: (data: unknown) => void;
  timeRemaining?: number;
}

export function VoteStrategy({
  voteType,
  dishes,
  onVoteComplete,
  timeRemaining,
}: VoteStrategyProps) {
  switch (voteType) {
    case 'tournament':
      return (
        <Tournament dishes={dishes} onComplete={onVoteComplete} timeRemaining={timeRemaining} />
      );
    case 'swipe':
      return (
        <SwipeVote dishes={dishes} onComplete={onVoteComplete} timeRemaining={timeRemaining} />
      );
    case 'ranking':
      return <RankVote dishes={dishes} onComplete={onVoteComplete} timeRemaining={timeRemaining} />;
    default:
      return null;
  }
}
