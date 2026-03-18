'use client';

import { useCallback, useMemo, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { DishCard } from '@features/dish';
import type { Dish } from '@features/dish';
import type { Matchup } from '../types';

interface TournamentProps {
  dishes: Dish[];
  onComplete: (data: unknown) => void;
  timeRemaining?: number;
}

function buildRounds(dishes: Dish[]): Matchup[] {
  const matchups: Matchup[] = [];
  for (let i = 0; i < dishes.length - 1; i += 2) {
    const dish1 = dishes[i];
    const dish2 = dishes[i + 1];
    if (dish1 && dish2) {
      matchups.push({ dish1, dish2 });
    }
  }
  return matchups;
}

export function Tournament({ dishes, onComplete, timeRemaining }: TournamentProps) {
  const initialMatchups = useMemo(() => buildRounds(dishes), [dishes]);
  const [matchups, setMatchups] = useState<Matchup[]>(initialMatchups);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [roundNumber, setRoundNumber] = useState(1);
  const [winners, setWinners] = useState<Dish[]>([]);

  const currentMatchup = matchups[currentIndex];

  const handleChoice = useCallback(
    (winnerId: string) => {
      if (!currentMatchup) return;

      const winner =
        currentMatchup.dish1.id === winnerId ? currentMatchup.dish1 : currentMatchup.dish2;
      const newWinners = [...winners, winner];

      if (currentIndex + 1 < matchups.length) {
        setCurrentIndex((prev) => prev + 1);
        setWinners(newWinners);
      } else {
        // Round complete
        if (newWinners.length === 1) {
          onComplete({ winnerId: newWinners[0]?.id, rounds: roundNumber });
        } else {
          // Start next round
          const nextMatchups = buildRounds(newWinners);
          setMatchups(nextMatchups);
          setCurrentIndex(0);
          setWinners([]);
          setRoundNumber((prev) => prev + 1);
        }
      }
    },
    [currentMatchup, currentIndex, matchups.length, winners, roundNumber, onComplete],
  );

  if (!currentMatchup) return null;

  return (
    <div className="flex flex-col items-center gap-6 px-4 py-6">
      <div className="text-center">
        <p className="text-sm font-medium text-gray-500 dark:text-gray-400">
          Vòng {roundNumber} &middot; Trận {currentIndex + 1}/{matchups.length}
        </p>
        {timeRemaining !== undefined && (
          <p className="mt-1 text-xs text-gray-400">Còn {timeRemaining}s</p>
        )}
      </div>

      <div className="flex w-full max-w-lg items-center gap-4">
        <AnimatePresence mode="wait">
          <motion.div
            key={`${currentMatchup.dish1.id}-${currentIndex}`}
            initial={{ x: -100, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: -100, opacity: 0 }}
            transition={{ duration: 0.3 }}
            className="flex-1"
          >
            <DishCard
              dish={currentMatchup.dish1}
              onClick={() => handleChoice(currentMatchup.dish1.id)}
            />
          </motion.div>
        </AnimatePresence>

        <span className="shrink-0 text-xl font-bold text-gray-400">VS</span>

        <AnimatePresence mode="wait">
          <motion.div
            key={`${currentMatchup.dish2.id}-${currentIndex}`}
            initial={{ x: 100, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: 100, opacity: 0 }}
            transition={{ duration: 0.3 }}
            className="flex-1"
          >
            <DishCard
              dish={currentMatchup.dish2}
              onClick={() => handleChoice(currentMatchup.dish2.id)}
            />
          </motion.div>
        </AnimatePresence>
      </div>
    </div>
  );
}
