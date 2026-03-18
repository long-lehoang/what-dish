'use client';

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from 'react';
import { useWakeLock } from '../../hooks/useWakeLock';
import type { Step } from '@features/dish';

interface CookModeContextValue {
  currentStep: number;
  totalSteps: number;
  isActive: boolean;
  steps: Step[];
  next: () => void;
  prev: () => void;
  exit: () => void;
}

const CookModeContext = createContext<CookModeContextValue | null>(null);

export function useCookMode(): CookModeContextValue {
  const ctx = useContext(CookModeContext);
  if (!ctx) {
    throw new Error('useCookMode must be used within CookMode.Root');
  }
  return ctx;
}

interface CookModeRootProps {
  steps: Step[];
  onExit?: () => void;
  children: ReactNode;
}

export function CookModeRoot({ steps, onExit, children }: CookModeRootProps) {
  const [currentStep, setCurrentStep] = useState(0);
  const [isActive, setIsActive] = useState(true);
  const { request, release } = useWakeLock();
  const touchStartX = useRef<number | null>(null);

  const totalSteps = steps.length;

  const next = useCallback(() => {
    setCurrentStep((prev) => Math.min(prev + 1, totalSteps - 1));
  }, [totalSteps]);

  const prev = useCallback(() => {
    setCurrentStep((prev) => Math.max(prev - 1, 0));
  }, []);

  const exit = useCallback(() => {
    setIsActive(false);
    void release();
    onExit?.();
  }, [release, onExit]);

  // Wake lock on enter
  useEffect(() => {
    void request();
    return () => {
      void release();
    };
  }, [request, release]);

  // Keyboard navigation
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'ArrowRight') next();
      else if (e.key === 'ArrowLeft') prev();
      else if (e.key === 'Escape') exit();
    }
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [next, prev, exit]);

  // Swipe gesture handlers
  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    const touch = e.touches[0];
    if (touch) {
      touchStartX.current = touch.clientX;
    }
  }, []);

  const handleTouchEnd = useCallback(
    (e: React.TouchEvent) => {
      const touch = e.changedTouches[0];
      if (touchStartX.current === null || !touch) return;

      const diff = touch.clientX - touchStartX.current;
      const threshold = 50;

      if (diff > threshold) {
        prev();
      } else if (diff < -threshold) {
        next();
      }

      touchStartX.current = null;
    },
    [next, prev],
  );

  if (!isActive) return null;

  return (
    <CookModeContext.Provider
      value={{ currentStep, totalSteps, isActive, steps, next, prev, exit }}
    >
      <div
        className="fixed inset-0 z-50 flex flex-col bg-gray-950 text-white"
        onTouchStart={handleTouchStart}
        onTouchEnd={handleTouchEnd}
        role="dialog"
        aria-modal="true"
        aria-label="Chế độ nấu ăn"
      >
        {children}
      </div>
    </CookModeContext.Provider>
  );
}
