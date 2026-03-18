import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useServingScale } from '../useServingScale';

describe('useServingScale', () => {
  it('initializes with original servings', () => {
    const { result } = renderHook(() => useServingScale(2));

    expect(result.current.servings).toBe(2);
  });

  it('scales amount proportionally', () => {
    const { result } = renderHook(() => useServingScale(2));

    // 2 servings, amount 100 → 100
    expect(result.current.scaleAmount(100)).toBe(100);

    act(() => {
      result.current.setServings(4);
    });

    // 4 servings, amount 100 → 200
    expect(result.current.scaleAmount(100)).toBe(200);
  });

  it('returns null for null amounts', () => {
    const { result } = renderHook(() => useServingScale(2));
    expect(result.current.scaleAmount(null)).toBeNull();
  });

  it('rounds to 1 decimal place', () => {
    const { result } = renderHook(() => useServingScale(3));

    act(() => {
      result.current.setServings(2);
    });

    // 100 * 2 / 3 = 66.666... → 66.7
    expect(result.current.scaleAmount(100)).toBe(66.7);
  });

  it('clamps servings to minimum 1', () => {
    const { result } = renderHook(() => useServingScale(2));

    act(() => {
      result.current.setServings(0);
    });

    expect(result.current.servings).toBe(1);

    act(() => {
      result.current.setServings(-5);
    });

    expect(result.current.servings).toBe(1);
  });

  it('clamps servings to maximum 20', () => {
    const { result } = renderHook(() => useServingScale(2));

    act(() => {
      result.current.setServings(25);
    });

    expect(result.current.servings).toBe(20);
  });

  it('handles scaling down to 1 serving', () => {
    const { result } = renderHook(() => useServingScale(4));

    act(() => {
      result.current.setServings(1);
    });

    // 200 * 1 / 4 = 50
    expect(result.current.scaleAmount(200)).toBe(50);
  });
});
