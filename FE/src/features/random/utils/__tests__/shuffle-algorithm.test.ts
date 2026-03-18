import { describe, it, expect } from 'vitest';
import {
  generateCardPositions,
  generateConvergePositions,
  fisherYatesShuffle,
} from '../shuffle-algorithm';

describe('generateCardPositions', () => {
  it('generates the correct number of positions', () => {
    const positions = generateCardPositions(8, 400, 600);
    expect(positions).toHaveLength(8);
  });

  it('each position has x, y, and rotation', () => {
    const positions = generateCardPositions(4, 400, 600);
    for (const pos of positions) {
      expect(pos).toHaveProperty('x');
      expect(pos).toHaveProperty('y');
      expect(pos).toHaveProperty('rotation');
      expect(typeof pos.x).toBe('number');
      expect(typeof pos.y).toBe('number');
      expect(typeof pos.rotation).toBe('number');
    }
  });

  it('distributes positions roughly in a circle (not all identical)', () => {
    const positions = generateCardPositions(6, 400, 600);
    const xs = positions.map((p) => p.x);
    const ys = positions.map((p) => p.y);
    // Not all x or y values should be the same
    const uniqueXs = new Set(xs.map((x) => Math.round(x)));
    const uniqueYs = new Set(ys.map((y) => Math.round(y)));
    expect(uniqueXs.size).toBeGreaterThan(1);
    expect(uniqueYs.size).toBeGreaterThan(1);
  });

  it('handles zero count', () => {
    const positions = generateCardPositions(0, 400, 600);
    expect(positions).toHaveLength(0);
  });
});

describe('generateConvergePositions', () => {
  it('generates the correct number of positions', () => {
    const positions = generateConvergePositions(6);
    expect(positions).toHaveLength(6);
  });

  it('all positions have x=0 (converge to center)', () => {
    const positions = generateConvergePositions(4);
    for (const pos of positions) {
      expect(pos.x).toBe(0);
    }
  });

  it('y offsets create a stacked look (increasingly negative)', () => {
    const positions = generateConvergePositions(5);
    for (let i = 0; i < positions.length; i++) {
      expect(positions[i]!.y).toBe(i * -2);
    }
  });

  it('rotations fan out from center', () => {
    const positions = generateConvergePositions(5);
    const mid = Math.floor(5 / 2);
    // Center card should have rotation = 0
    expect(positions[mid]!.rotation).toBe(0);
    // Cards before center should have negative rotation
    expect(positions[0]!.rotation).toBeLessThan(0);
    // Cards after center should have positive rotation
    expect(positions[4]!.rotation).toBeGreaterThan(0);
  });
});

describe('fisherYatesShuffle', () => {
  it('returns a new array (does not mutate input)', () => {
    const input = [1, 2, 3, 4, 5];
    const original = [...input];
    fisherYatesShuffle(input);
    expect(input).toEqual(original);
  });

  it('returns array with same elements', () => {
    const input = [1, 2, 3, 4, 5];
    const result = fisherYatesShuffle(input);
    expect(result).toHaveLength(input.length);
    expect(result.sort()).toEqual([...input].sort());
  });

  it('handles empty array', () => {
    expect(fisherYatesShuffle([])).toEqual([]);
  });

  it('handles single element', () => {
    expect(fisherYatesShuffle([42])).toEqual([42]);
  });

  it('produces different orderings (probabilistic)', () => {
    const input = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    const results = new Set<string>();
    for (let i = 0; i < 20; i++) {
      results.add(JSON.stringify(fisherYatesShuffle(input)));
    }
    // With 10 elements and 20 shuffles, we should get multiple unique orderings
    expect(results.size).toBeGreaterThan(1);
  });
});
