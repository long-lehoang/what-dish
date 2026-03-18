'use client';

import { useCallback, useEffect, useState } from 'react';

export function useMediaQuery(query: string): boolean {
  const [matches, setMatches] = useState(false);

  const handleChange = useCallback((event: MediaQueryListEvent) => {
    setMatches(event.matches);
  }, []);

  useEffect(() => {
    const mediaQuery = window.matchMedia(query);
    setMatches(mediaQuery.matches);

    mediaQuery.addEventListener('change', handleChange);
    return () => {
      mediaQuery.removeEventListener('change', handleChange);
    };
  }, [query, handleChange]);

  return matches;
}

export function useIsMobile(): boolean {
  return useMediaQuery('(max-width: 767px)');
}

export function useIsDesktop(): boolean {
  return useMediaQuery('(min-width: 1024px)');
}
