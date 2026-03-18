'use client';

import { useCallback, useRef, useState } from 'react';
import { apiClient } from '@shared/lib/api-client';
import type { Dish } from '@features/dish/types';
import type { ShufflePhase, RandomFilters } from '../types';
import { PHASE_DURATIONS } from '../utils/animation-config';
import { useHistoryStore } from '../stores/history-store';

interface UseShuffleReturn {
  phase: ShufflePhase;
  selectedDish: Dish | null;
  triggerShuffle: (filters?: RandomFilters) => void;
  reset: () => void;
  isAnimating: boolean;
}

export function useShuffle(): UseShuffleReturn {
  const [phase, setPhase] = useState<ShufflePhase>('idle');
  const [selectedDish, setSelectedDish] = useState<Dish | null>(null);
  const timeoutRefs = useRef<ReturnType<typeof setTimeout>[]>([]);
  const addEntry = useHistoryStore((s) => s.addEntry);
  const getExcludeIds = useHistoryStore((s) => s.getExcludeIds);

  const clearTimeouts = useCallback(() => {
    timeoutRefs.current.forEach(clearTimeout);
    timeoutRefs.current = [];
  }, []);

  const schedulePhase = useCallback((nextPhase: ShufflePhase, delayMs: number) => {
    const id = setTimeout(() => {
      setPhase(nextPhase);
    }, delayMs);
    timeoutRefs.current.push(id);
    return id;
  }, []);

  const triggerShuffle = useCallback(
    (filters?: RandomFilters) => {
      clearTimeouts();
      setSelectedDish(null);
      setPhase('shuffle');

      const excludeIds = getExcludeIds();
      const params = new URLSearchParams();

      if (filters?.dishType) params.set('dish_type', filters.dishType);
      if (filters?.difficulty) params.set('difficulty', filters.difficulty);
      if (filters?.maxCookTime) params.set('max_cook_time', String(filters.maxCookTime));
      if (filters?.tags) params.set('tags', filters.tags);
      if (excludeIds.length > 0) params.set('exclude_ids', excludeIds.join(','));

      const queryString = params.toString();
      const path = `/api/v1/recipes/random${queryString ? `?${queryString}` : ''}`;

      const fetchPromise = apiClient.get<Dish>(path);

      // Phase 1 -> 2: shuffle -> converge
      schedulePhase('converge', PHASE_DURATIONS.shuffle);

      // Phase 2 -> 3: converge -> select
      const selectDelay = PHASE_DURATIONS.shuffle + PHASE_DURATIONS.converge;
      const selectId = setTimeout(() => {
        setPhase('select');

        // Ensure fetch is complete before reveal
        fetchPromise
          .then((dish) => {
            setSelectedDish(dish);
            addEntry({ dishId: dish.id, dishName: dish.name });

            // Phase 3 -> 4: select -> reveal
            const revealId = setTimeout(() => {
              setPhase('reveal');

              // Phase 4 -> 5: reveal -> settle
              const settleId = setTimeout(() => {
                setPhase('settle');
              }, PHASE_DURATIONS.reveal);
              timeoutRefs.current.push(settleId);
            }, PHASE_DURATIONS.select);
            timeoutRefs.current.push(revealId);
          })
          .catch(() => {
            setPhase('idle');
            setSelectedDish(null);
          });
      }, selectDelay);
      timeoutRefs.current.push(selectId);
    },
    [clearTimeouts, schedulePhase, addEntry, getExcludeIds],
  );

  const reset = useCallback(() => {
    clearTimeouts();
    setPhase('idle');
    setSelectedDish(null);
  }, [clearTimeouts]);

  const isAnimating = phase !== 'idle' && phase !== 'settle';

  return { phase, selectedDish, triggerShuffle, reset, isAnimating };
}
