import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { HISTORY_DAYS } from '@shared/constants';
import type { HistoryEntry } from '../types';

interface HistoryState {
  history: HistoryEntry[];
  addEntry: (entry: Omit<HistoryEntry, 'timestamp'>) => void;
  clearExpired: () => void;
  getExcludeIds: () => string[];
}

function filterExpired(entries: HistoryEntry[]): HistoryEntry[] {
  const expiryMs = HISTORY_DAYS * 24 * 60 * 60 * 1000;
  const now = Date.now();
  return entries.filter((entry) => now - entry.timestamp < expiryMs);
}

export const useHistoryStore = create<HistoryState>()(
  persist(
    (set, get) => ({
      history: [],

      addEntry: (entry) => {
        set((state) => ({
          history: [...state.history, { ...entry, timestamp: Date.now() }],
        }));
      },

      clearExpired: () => {
        set((state) => ({
          history: filterExpired(state.history),
        }));
      },

      getExcludeIds: () => {
        const { history } = get();
        const valid = filterExpired(history);
        return valid.map((entry) => entry.dishId);
      },
    }),
    {
      name: 'random-history',
      onRehydrateStorage: () => {
        return (state) => {
          if (state) {
            state.clearExpired();
          }
        };
      },
    },
  ),
);
