import type { DishFilters } from '@features/dish/types';

export type ShufflePhase = 'idle' | 'shuffle' | 'converge' | 'select' | 'reveal' | 'settle';

export interface ShuffleState {
  phase: ShufflePhase;
  selectedIndex: number | null;
  cardPositions: CardPosition[];
}

export interface CardPosition {
  x: number;
  y: number;
  rotation: number;
}

export interface HistoryEntry {
  dishId: string;
  dishName: string;
  timestamp: number;
}

export interface RandomFilters extends DishFilters {
  excludeIds?: string[];
}
