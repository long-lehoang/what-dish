'use client';

import { useCallback, useState } from 'react';

interface UseServingScaleReturn {
  servings: number;
  setServings: (value: number) => void;
  scaleAmount: (amount: number | null | undefined) => number | null;
}

const MIN_SERVINGS = 1;
const MAX_SERVINGS = 20;

export function useServingScale(originalServings: number): UseServingScaleReturn {
  const [servings, setServingsRaw] = useState(originalServings);

  const setServings = useCallback((value: number) => {
    setServingsRaw(Math.max(MIN_SERVINGS, Math.min(MAX_SERVINGS, value)));
  }, []);

  const scaleAmount = useCallback(
    (amount: number | null | undefined): number | null => {
      if (amount === null || amount === undefined) return null;
      const scaled = (amount * servings) / originalServings;
      return Math.round(scaled * 10) / 10;
    },
    [servings, originalServings],
  );

  return { servings, setServings, scaleAmount };
}
