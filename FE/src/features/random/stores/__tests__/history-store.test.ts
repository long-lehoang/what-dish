import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { useHistoryStore } from '../history-store';

describe('useHistoryStore', () => {
  beforeEach(() => {
    // Reset store state before each test
    useHistoryStore.setState({ history: [] });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('starts with empty history', () => {
    expect(useHistoryStore.getState().history).toEqual([]);
  });

  it('adds an entry with timestamp', () => {
    const now = Date.now();
    vi.spyOn(Date, 'now').mockReturnValue(now);

    useHistoryStore.getState().addEntry({ dishId: 'dish-1', dishName: 'Phở Bò' });

    const { history } = useHistoryStore.getState();
    expect(history).toHaveLength(1);
    expect(history[0]).toEqual({
      dishId: 'dish-1',
      dishName: 'Phở Bò',
      timestamp: now,
    });
  });

  it('adds multiple entries', () => {
    const { addEntry } = useHistoryStore.getState();

    addEntry({ dishId: 'dish-1', dishName: 'Phở Bò' });
    addEntry({ dishId: 'dish-2', dishName: 'Bún Chả' });

    expect(useHistoryStore.getState().history).toHaveLength(2);
  });

  it('clears expired entries (older than 7 days)', () => {
    const now = Date.now();
    const eightDaysAgo = now - 8 * 24 * 60 * 60 * 1000;
    const twoDaysAgo = now - 2 * 24 * 60 * 60 * 1000;

    useHistoryStore.setState({
      history: [
        { dishId: 'old', dishName: 'Old Dish', timestamp: eightDaysAgo },
        { dishId: 'recent', dishName: 'Recent Dish', timestamp: twoDaysAgo },
      ],
    });

    useHistoryStore.getState().clearExpired();

    const { history } = useHistoryStore.getState();
    expect(history).toHaveLength(1);
    expect(history[0]!.dishId).toBe('recent');
  });

  it('getExcludeIds returns only non-expired dish IDs', () => {
    const now = Date.now();
    const eightDaysAgo = now - 8 * 24 * 60 * 60 * 1000;
    const twoDaysAgo = now - 2 * 24 * 60 * 60 * 1000;

    useHistoryStore.setState({
      history: [
        { dishId: 'old', dishName: 'Old Dish', timestamp: eightDaysAgo },
        { dishId: 'recent-1', dishName: 'Recent 1', timestamp: twoDaysAgo },
        { dishId: 'recent-2', dishName: 'Recent 2', timestamp: now },
      ],
    });

    const excludeIds = useHistoryStore.getState().getExcludeIds();
    expect(excludeIds).toEqual(['recent-1', 'recent-2']);
    expect(excludeIds).not.toContain('old');
  });

  it('getExcludeIds returns empty array for empty history', () => {
    expect(useHistoryStore.getState().getExcludeIds()).toEqual([]);
  });
});
