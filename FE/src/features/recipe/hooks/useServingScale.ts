'use client';

import { useCallback, useState } from 'react';

interface UseServingScaleReturn {
  servings: number;
  setServings: (value: number) => void;
  scaleAmount: (amount: number | null) => number | null;
}

const MIN_SERVINGS = 1;
const MAX_SERVINGS = 20;

export function useServingScale(originalServings: number): UseServingScaleReturn {
  const [servings, setServingsRaw] = useState(originalServings);

  const setServings = useCallback((value: number) => {
    setServingsRaw(Math.max(MIN_SERVINGS, Math.min(MAX_SERVINGS, value)));
  }, []);

  const scaleAmount = useCallback(
    (amount: number | null): number | null => {
      if (amount === null) return null;
      const scaled = (amount * servings) / originalServings;
      return Math.round(scaled * 10) / 10;
    },
    [servings, originalServings],
  );

  return { servings, setServings, scaleAmount };
}
