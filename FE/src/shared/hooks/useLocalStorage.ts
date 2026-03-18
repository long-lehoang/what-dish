'use client';

import { useCallback, useEffect, useState } from 'react';

export function useLocalStorage<T>(
  key: string,
  initialValue: T,
): [T, (value: T | ((prev: T) => T)) => void] {
  const [storedValue, setStoredValue] = useState<T>(initialValue);

  // Hydrate from localStorage after mount (SSR-safe)
  useEffect(() => {
    try {
      const item = window.localStorage.getItem(key);
      if (item !== null) {
        setStoredValue(JSON.parse(item) as T);
      }
    } catch (error) {
      // eslint-disable-next-line no-console
      console.warn(`Error reading localStorage key "${key}":`, error);
    }
  }, [key]);

  const setValue = useCallback(
    (value: T | ((prev: T) => T)) => {
      setStoredValue((prev) => {
        const nextValue = value instanceof Function ? value(prev) : value;

        try {
          window.localStorage.setItem(key, JSON.stringify(nextValue));
        } catch (error) {
          // eslint-disable-next-line no-console
          console.warn(`Error writing localStorage key "${key}":`, error);
        }

        return nextValue;
      });
    },
    [key],
  );

  return [storedValue, setValue];
}
