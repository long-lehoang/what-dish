import { describe, it, expect, vi, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useTimer } from '../useTimer';

describe('useTimer', () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it('initializes with correct values', () => {
    const { result } = renderHook(() => useTimer(120));

    expect(result.current.seconds).toBe(120);
    expect(result.current.isRunning).toBe(false);
    expect(result.current.isPaused).toBe(false);
    expect(result.current.formattedTime).toBe('02:00');
  });

  it('formats time correctly', () => {
    const { result } = renderHook(() => useTimer(65));
    expect(result.current.formattedTime).toBe('01:05');
  });

  it('starts counting down', () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useTimer(10));

    act(() => {
      result.current.start();
    });

    expect(result.current.isRunning).toBe(true);
    expect(result.current.isPaused).toBe(false);

    act(() => {
      vi.advanceTimersByTime(3000);
    });

    expect(result.current.seconds).toBe(7);
    expect(result.current.formattedTime).toBe('00:07');
  });

  it('pauses and resumes', () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useTimer(10));

    act(() => {
      result.current.start();
    });

    act(() => {
      vi.advanceTimersByTime(2000);
    });

    expect(result.current.seconds).toBe(8);

    act(() => {
      result.current.pause();
    });

    expect(result.current.isPaused).toBe(true);

    // Time should not advance while paused
    act(() => {
      vi.advanceTimersByTime(3000);
    });

    expect(result.current.seconds).toBe(8);

    act(() => {
      result.current.resume();
    });

    expect(result.current.isPaused).toBe(false);

    act(() => {
      vi.advanceTimersByTime(2000);
    });

    expect(result.current.seconds).toBe(6);
  });

  it('stops at zero and sets isRunning to false', () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useTimer(3));

    act(() => {
      result.current.start();
    });

    act(() => {
      vi.advanceTimersByTime(5000);
    });

    expect(result.current.seconds).toBe(0);
    expect(result.current.isRunning).toBe(false);
  });

  it('resets to initial seconds', () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useTimer(30));

    act(() => {
      result.current.start();
    });

    act(() => {
      vi.advanceTimersByTime(5000);
    });

    expect(result.current.seconds).toBe(25);

    act(() => {
      result.current.reset();
    });

    expect(result.current.seconds).toBe(30);
    expect(result.current.isRunning).toBe(false);
    expect(result.current.isPaused).toBe(false);
  });
});
