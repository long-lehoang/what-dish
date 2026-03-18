import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useRerollLimit } from '../useRerollLimit';

describe('useRerollLimit', () => {
  it('starts with 0 rerolls and canReroll true', () => {
    const { result } = renderHook(() => useRerollLimit());

    expect(result.current.rerollCount).toBe(0);
    expect(result.current.canReroll).toBe(true);
  });

  it('increments reroll count', () => {
    const { result } = renderHook(() => useRerollLimit());

    act(() => {
      result.current.incrementReroll();
    });

    expect(result.current.rerollCount).toBe(1);
    expect(result.current.canReroll).toBe(true);
  });

  it('disallows reroll after MAX_REROLLS (3)', () => {
    const { result } = renderHook(() => useRerollLimit());

    act(() => {
      result.current.incrementReroll();
      result.current.incrementReroll();
      result.current.incrementReroll();
    });

    expect(result.current.rerollCount).toBe(3);
    expect(result.current.canReroll).toBe(false);
  });

  it('resets rerolls to zero', () => {
    const { result } = renderHook(() => useRerollLimit());

    act(() => {
      result.current.incrementReroll();
      result.current.incrementReroll();
    });

    expect(result.current.rerollCount).toBe(2);

    act(() => {
      result.current.resetRerolls();
    });

    expect(result.current.rerollCount).toBe(0);
    expect(result.current.canReroll).toBe(true);
  });
});
