import type { CardPosition } from '../types';

/**
 * Generate scattered card positions within a container.
 * Positions are distributed in a rough circle with random offsets.
 */
export function generateCardPositions(
  count: number,
  containerWidth: number,
  containerHeight: number,
): CardPosition[] {
  const positions: CardPosition[] = [];
  const radiusX = containerWidth * 0.35;
  const radiusY = containerHeight * 0.3;

  for (let i = 0; i < count; i++) {
    const angle = (i / count) * Math.PI * 2;
    const jitterX = (Math.random() - 0.5) * radiusX * 0.4;
    const jitterY = (Math.random() - 0.5) * radiusY * 0.4;

    positions.push({
      x: Math.cos(angle) * radiusX + jitterX,
      y: Math.sin(angle) * radiusY + jitterY,
      rotation: (Math.random() - 0.5) * 30,
    });
  }

  return positions;
}

/**
 * Generate positions for cards converging to center as a neat stack.
 * Each card is offset slightly to create a fanned-stack look.
 */
export function generateConvergePositions(count: number): CardPosition[] {
  const positions: CardPosition[] = [];

  for (let i = 0; i < count; i++) {
    positions.push({
      x: 0,
      y: i * -2,
      rotation: (i - Math.floor(count / 2)) * 1.5,
    });
  }

  return positions;
}

/**
 * Fisher-Yates shuffle algorithm. Returns a new shuffled array (pure function).
 */
export function fisherYatesShuffle<T>(array: readonly T[]): T[] {
  const result = [...array];
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    const temp = result[i]!;
    result[i] = result[j]!;
    result[j] = temp;
  }
  return result;
}
