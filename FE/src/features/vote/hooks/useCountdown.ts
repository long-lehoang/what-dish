'use client';

import { useCallback, useEffect, useRef, useState } from 'react';

interface UseCountdownReturn {
  remaining: number;
  isActive: boolean;
  start: () => void;
  stop: () => void;
}

export function useCountdown(seconds: number, onComplete?: () => void): UseCountdownReturn {
  const [remaining, setRemaining] = useState(seconds);
  const [isActive, setIsActive] = useState(false);
  const rafRef = useRef<number | null>(null);
  const startTimeRef = useRef<number>(0);
  const onCompleteRef = useRef(onComplete);
  onCompleteRef.current = onComplete;

  const stop = useCallback(() => {
    setIsActive(false);
    if (rafRef.current !== null) {
      cancelAnimationFrame(rafRef.current);
      rafRef.current = null;
    }
  }, []);

  const tick = useCallback(() => {
    const elapsed = (performance.now() - startTimeRef.current) / 1000;
    const newRemaining = Math.max(0, seconds - elapsed);

    setRemaining(Math.ceil(newRemaining));

    if (newRemaining <= 0) {
      stop();
      onCompleteRef.current?.();
      return;
    }

    rafRef.current = requestAnimationFrame(tick);
  }, [seconds, stop]);

  const start = useCallback(() => {
    setIsActive(true);
    setRemaining(seconds);
    startTimeRef.current = performance.now();
    rafRef.current = requestAnimationFrame(tick);
  }, [seconds, tick]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (rafRef.current !== null) {
        cancelAnimationFrame(rafRef.current);
      }
    };
  }, []);

  return { remaining, isActive, start, stop };
}
