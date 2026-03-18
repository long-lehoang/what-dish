'use client';

import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from 'react';

type Theme = 'light' | 'dark';

interface ThemeContextValue {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

const STORAGE_KEY = 'toi-nay-an-gi-theme';

function getInitialTheme(): Theme {
  if (typeof window === 'undefined') {
    return 'light';
  }

  try {
    const stored = window.localStorage.getItem(STORAGE_KEY);
    if (stored === 'light' || stored === 'dark') {
      return stored;
    }
  } catch {
    // localStorage unavailable
  }

  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function applyThemeToDOM(theme: Theme): void {
  const root = document.documentElement;
  root.classList.remove('light', 'dark');
  root.classList.add(theme);
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<Theme>('light');

  // Initialize on mount (SSR-safe)
  useEffect(() => {
    const initial = getInitialTheme();
    setThemeState(initial);
    applyThemeToDOM(initial);
  }, []);

  const setTheme = useCallback((newTheme: Theme) => {
    setThemeState(newTheme);
    applyThemeToDOM(newTheme);
    try {
      window.localStorage.setItem(STORAGE_KEY, newTheme);
    } catch {
      // localStorage unavailable
    }
  }, []);

  const toggleTheme = useCallback(() => {
    setTheme(theme === 'light' ? 'dark' : 'light');
  }, [theme, setTheme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme(): ThemeContextValue {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}
