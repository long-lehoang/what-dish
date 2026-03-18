import { describe, it, expect, vi, afterEach } from 'vitest';
import {
  cn,
  formatDuration,
  formatCurrency,
  slugify,
  getRandomExcluding,
  isExpired,
  generateRoomCode,
} from '../utils';

describe('cn', () => {
  it('merges class names', () => {
    expect(cn('px-4', 'py-2')).toBe('px-4 py-2');
  });

  it('handles conditional classes', () => {
    expect(cn('base', false && 'hidden', 'extra')).toBe('base extra');
  });

  it('resolves Tailwind conflicts (last wins)', () => {
    const result = cn('px-4', 'px-6');
    expect(result).toBe('px-6');
  });

  it('handles empty input', () => {
    expect(cn()).toBe('');
  });
});

describe('formatDuration', () => {
  it('formats minutes under 60', () => {
    expect(formatDuration(30)).toBe('30 phút');
  });

  it('formats exact hours', () => {
    expect(formatDuration(120)).toBe('2 giờ');
  });

  it('formats hours and minutes', () => {
    expect(formatDuration(90)).toBe('1 giờ 30 phút');
  });

  it('handles zero minutes', () => {
    expect(formatDuration(0)).toBe('0 phút');
  });
});

describe('formatCurrency', () => {
  it('formats VND to thousands', () => {
    expect(formatCurrency(50000)).toBe('~50k');
  });

  it('rounds to nearest thousand', () => {
    expect(formatCurrency(45500)).toBe('~46k');
  });

  it('handles zero', () => {
    expect(formatCurrency(0)).toBe('~0k');
  });
});

describe('slugify', () => {
  it('converts Vietnamese text to slug', () => {
    expect(slugify('Phở Bò Hà Nội')).toBe('pho-bo-ha-noi');
  });

  it('handles đ character', () => {
    expect(slugify('Bánh Đúc')).toBe('banh-duc');
  });

  it('removes special characters', () => {
    expect(slugify('Cơm & Thịt!')).toBe('com-thit');
  });

  it('collapses multiple dashes', () => {
    expect(slugify('a   b---c')).toBe('a-b-c');
  });

  it('trims leading/trailing dashes', () => {
    expect(slugify(' -hello- ')).toBe('hello');
  });
});

describe('getRandomExcluding', () => {
  const items = [
    { id: '1', name: 'A' },
    { id: '2', name: 'B' },
    { id: '3', name: 'C' },
  ];

  it('returns an item not in the exclude list', () => {
    const result = getRandomExcluding(items, ['1', '2'], 'id');
    expect(result).toEqual({ id: '3', name: 'C' });
  });

  it('returns null when all items are excluded', () => {
    const result = getRandomExcluding(items, ['1', '2', '3'], 'id');
    expect(result).toBeNull();
  });

  it('returns any item when no exclusions', () => {
    const result = getRandomExcluding(items, [], 'id');
    expect(result).not.toBeNull();
    expect(items).toContainEqual(result);
  });

  it('handles empty items array', () => {
    const result = getRandomExcluding([], ['1'], 'id' as never);
    expect(result).toBeNull();
  });
});

describe('isExpired', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns true for timestamps older than specified days', () => {
    const eightDaysAgo = Date.now() - 8 * 24 * 60 * 60 * 1000;
    expect(isExpired(eightDaysAgo, 7)).toBe(true);
  });

  it('returns false for timestamps within the specified days', () => {
    const oneDayAgo = Date.now() - 1 * 24 * 60 * 60 * 1000;
    expect(isExpired(oneDayAgo, 7)).toBe(false);
  });

  it('returns false for current timestamp', () => {
    expect(isExpired(Date.now(), 7)).toBe(false);
  });
});

describe('generateRoomCode', () => {
  it('generates a 6-character code', () => {
    const code = generateRoomCode();
    expect(code).toHaveLength(6);
  });

  it('uses only uppercase letters and digits (no ambiguous chars)', () => {
    const code = generateRoomCode();
    expect(code).toMatch(/^[ABCDEFGHJKLMNPQRSTUVWXYZ23456789]+$/);
  });

  it('does not contain ambiguous characters (0, O, 1, I)', () => {
    // Run multiple times to increase confidence
    for (let i = 0; i < 50; i++) {
      const code = generateRoomCode();
      expect(code).not.toMatch(/[0O1I]/);
    }
  });
});
