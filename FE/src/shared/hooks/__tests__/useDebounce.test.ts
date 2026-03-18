import { describe, it, expect, vi, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useDebounce } from '../useDebounce';

describe('useDebounce', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns the initial value immediately', () => {
    const { result } = renderHook(() => useDebounce('hello', 500));
    expect(result.current).toBe('hello');
  });

  it('does not update the value before the delay', () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: 'hello', delay: 500 },
    });

    rerender({ value: 'world', delay: 500 });

    // Before delay, still the old value
    expect(result.current).toBe('hello');

    vi.useRealTimers();
  });

  it('updates the value after the delay', () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: 'hello', delay: 500 },
    });

    rerender({ value: 'world', delay: 500 });

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(result.current).toBe('world');

    vi.useRealTimers();
  });

  it('resets the timer on rapid changes', () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: 'a', delay: 300 },
    });

    rerender({ value: 'b', delay: 300 });
    act(() => {
      vi.advanceTimersByTime(200);
    });

    rerender({ value: 'c', delay: 300 });
    act(() => {
      vi.advanceTimersByTime(200);
    });

    // 'b' should NOT have appeared — timer was reset
    expect(result.current).toBe('a');

    act(() => {
      vi.advanceTimersByTime(100);
    });

    // Now 300ms have passed since last change to 'c'
    expect(result.current).toBe('c');

    vi.useRealTimers();
  });
});
