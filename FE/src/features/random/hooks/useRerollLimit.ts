'use client';

import { useCallback, useState } from 'react';
import { MAX_REROLLS } from '@shared/constants';

interface UseRerollLimitReturn {
  rerollCount: number;
  canReroll: boolean;
  incrementReroll: () => void;
  resetRerolls: () => void;
}

export function useRerollLimit(): UseRerollLimitReturn {
  const [rerollCount, setRerollCount] = useState(0);

  const canReroll = rerollCount < MAX_REROLLS;

  const incrementReroll = useCallback(() => {
    setRerollCount((prev) => prev + 1);
  }, []);

  const resetRerolls = useCallback(() => {
    setRerollCount(0);
  }, []);

  return { rerollCount, canReroll, incrementReroll, resetRerolls };
}
