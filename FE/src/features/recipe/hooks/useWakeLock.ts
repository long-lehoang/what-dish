'use client';

import { useCallback, useEffect, useRef, useState } from 'react';

interface UseWakeLockReturn {
  isSupported: boolean;
  isActive: boolean;
  request: () => Promise<void>;
  release: () => Promise<void>;
}

export function useWakeLock(): UseWakeLockReturn {
  const [isSupported, setIsSupported] = useState(false);
  const [isActive, setIsActive] = useState(false);
  const wakeLockRef = useRef<WakeLockSentinel | null>(null);

  useEffect(() => {
    setIsSupported('wakeLock' in navigator);
  }, []);

  const request = useCallback(async () => {
    if (!('wakeLock' in navigator)) return;

    try {
      wakeLockRef.current = await navigator.wakeLock.request('screen');
      setIsActive(true);

      wakeLockRef.current.addEventListener('release', () => {
        setIsActive(false);
        wakeLockRef.current = null;
      });
    } catch {
      setIsActive(false);
    }
  }, []);

  const release = useCallback(async () => {
    if (wakeLockRef.current) {
      try {
        await wakeLockRef.current.release();
      } catch {
        // Already released
      }
      wakeLockRef.current = null;
      setIsActive(false);
    }
  }, []);

  // Re-acquire wake lock when tab becomes visible again
  useEffect(() => {
    function handleVisibilityChange() {
      if (document.visibilityState === 'visible' && isActive && wakeLockRef.current === null) {
        void request();
      }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [isActive, request]);

  // Release on unmount
  useEffect(() => {
    return () => {
      if (wakeLockRef.current) {
        void wakeLockRef.current.release().catch(() => {
          // Ignore errors on cleanup
        });
        wakeLockRef.current = null;
      }
    };
  }, []);

  return { isSupported, isActive, request, release };
}
