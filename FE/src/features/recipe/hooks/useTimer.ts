'use client';

import { useCallback, useEffect, useRef, useState } from 'react';

interface UseTimerReturn {
  seconds: number;
  isRunning: boolean;
  isPaused: boolean;
  start: () => void;
  pause: () => void;
  resume: () => void;
  reset: () => void;
  formattedTime: string;
}

function formatTime(totalSeconds: number): string {
  const mins = Math.floor(totalSeconds / 60);
  const secs = totalSeconds % 60;
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
}

export function useTimer(initialSeconds: number): UseTimerReturn {
  const [seconds, setSeconds] = useState(initialSeconds);
  const [isRunning, setIsRunning] = useState(false);
  const [isPaused, setIsPaused] = useState(false);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const clearTimer = useCallback(() => {
    if (intervalRef.current !== null) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  }, []);

  const start = useCallback(() => {
    clearTimer();
    setSeconds(initialSeconds);
    setIsRunning(true);
    setIsPaused(false);

    intervalRef.current = setInterval(() => {
      setSeconds((prev) => {
        if (prev <= 1) {
          clearTimer();
          setIsRunning(false);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  }, [initialSeconds, clearTimer]);

  const pause = useCallback(() => {
    if (isRunning && !isPaused) {
      clearTimer();
      setIsPaused(true);
    }
  }, [isRunning, isPaused, clearTimer]);

  const resume = useCallback(() => {
    if (isRunning && isPaused) {
      setIsPaused(false);
      intervalRef.current = setInterval(() => {
        setSeconds((prev) => {
          if (prev <= 1) {
            clearTimer();
            setIsRunning(false);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }
  }, [isRunning, isPaused, clearTimer]);

  const reset = useCallback(() => {
    clearTimer();
    setSeconds(initialSeconds);
    setIsRunning(false);
    setIsPaused(false);
  }, [initialSeconds, clearTimer]);

  useEffect(() => {
    return () => {
      clearTimer();
    };
  }, [clearTimer]);

  return {
    seconds,
    isRunning,
    isPaused,
    start,
    pause,
    resume,
    reset,
    formattedTime: formatTime(seconds),
  };
}
