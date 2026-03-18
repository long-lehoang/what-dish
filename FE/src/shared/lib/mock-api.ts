/**
 * Mock API handler — returns mock data when the real backend is unavailable.
 * Parses request paths and returns appropriate responses matching the real API contract.
 */

import type { Dish, DishDetail, DishListResponse } from '@features/dish/types';
import type { VoteRoom } from '@features/vote/types';
import { MOCK_DISHES, MOCK_ROOM } from './mock-data';
import { generateRoomCode } from './utils';

// ---- Internal helpers ----

function delay(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function matchPath(path: string, pattern: string): Record<string, string> | null {
  const pathParts = path.split('/');
  const patternParts = pattern.split('/');

  if (pathParts.length !== patternParts.length) return null;

  const params: Record<string, string> = {};
  for (let i = 0; i < patternParts.length; i++) {
    const pp = patternParts[i]!;
    const pv = pathParts[i]!;
    if (pp.startsWith(':')) {
      params[pp.slice(1)] = pv;
    } else if (pp !== pv) {
      return null;
    }
  }
  return params;
}

function parseQuery(url: string): URLSearchParams {
  const idx = url.indexOf('?');
  return idx >= 0 ? new URLSearchParams(url.slice(idx + 1)) : new URLSearchParams();
}

function stripBase(url: string): string {
  // Remove base URL prefix if present, keep /api/...
  const idx = url.indexOf('/api/');
  return idx >= 0 ? url.slice(idx) : url;
}

// ---- Route handlers ----

function handleGetDishes(query: URLSearchParams): DishListResponse {
  let dishes: Dish[] = MOCK_DISHES;

  const search = query.get('search');
  if (search) {
    const lower = search.toLowerCase();
    dishes = dishes.filter(
      (d) =>
        d.name.toLowerCase().includes(lower) || d.tags.some((t) => t.toLowerCase().includes(lower)),
    );
  }

  const category = query.get('category');
  if (category) {
    dishes = dishes.filter((d) => d.category === category);
  }

  const difficulty = query.get('difficulty');
  if (difficulty) {
    dishes = dishes.filter((d) => d.difficulty <= Number(difficulty));
  }

  const maxTime = query.get('maxTime');
  if (maxTime) {
    const max = Number(maxTime);
    dishes = dishes.filter((d) => {
      const total = (d.prepTime ?? 0) + (d.cookTime ?? 0);
      return total <= max;
    });
  }

  const page = Number(query.get('page') ?? 1);
  const pageSize = Number(query.get('pageSize') ?? 20);
  const start = (page - 1) * pageSize;
  const paged = dishes.slice(start, start + pageSize);

  return {
    dishes: paged,
    total: dishes.length,
    page,
    pageSize,
  };
}

function handleGetDishBySlug(slug: string): DishDetail | null {
  return MOCK_DISHES.find((d) => d.slug === slug) ?? null;
}

function handleGetRandomDish(query: URLSearchParams): Dish | null {
  let dishes: Dish[] = MOCK_DISHES;

  const category = query.get('category');
  if (category) {
    dishes = dishes.filter((d) => d.category === category);
  }

  const difficulty = query.get('difficulty');
  if (difficulty) {
    dishes = dishes.filter((d) => d.difficulty <= Number(difficulty));
  }

  const maxTime = query.get('maxTime');
  if (maxTime) {
    const max = Number(maxTime);
    dishes = dishes.filter((d) => {
      const total = (d.prepTime ?? 0) + (d.cookTime ?? 0);
      return total <= max;
    });
  }

  const exclude = query.get('exclude');
  if (exclude) {
    const ids = exclude.split(',');
    dishes = dishes.filter((d) => !ids.includes(d.id));
  }

  if (dishes.length === 0) return MOCK_DISHES[0] ?? null;

  const idx = Math.floor(Math.random() * dishes.length);
  return dishes[idx] ?? null;
}

function handleCreateRoom(body: Record<string, unknown>): VoteRoom {
  return {
    ...MOCK_ROOM,
    id: `room-${Date.now()}`,
    code: generateRoomCode(),
    hostName: String(body.hostName ?? 'Host'),
    voteType: (body.voteType as VoteRoom['voteType']) ?? 'tournament',
    timerSecs: Number(body.timerSecs ?? 60),
    status: 'waiting',
    participants: [
      { name: String(body.hostName ?? 'Host'), avatarColor: '#FF6B35', hasVoted: false },
    ],
  };
}

function handleGetRoom(roomId: string): VoteRoom {
  return { ...MOCK_ROOM, id: roomId };
}

function handleJoinRoom(roomId: string, body: Record<string, unknown>): VoteRoom {
  const name = String(body.name ?? 'Player');
  const colors = ['#FF6B35', '#E63946', '#FFB703', '#2EC4B6', '#9B5DE5'];
  return {
    ...MOCK_ROOM,
    id: roomId,
    participants: [
      ...MOCK_ROOM.participants,
      { name, avatarColor: colors[Math.floor(Math.random() * colors.length)]!, hasVoted: false },
    ],
  };
}

// ---- Public interface ----

export async function mockApiHandler<T>(
  method: string,
  url: string,
  body?: unknown,
): Promise<T | null> {
  // Simulate network latency
  await delay(200 + Math.random() * 300);

  const path = stripBase(url);
  const pathOnly = path.split('?')[0]!;
  const query = parseQuery(path);

  // GET /api/dishes/random
  if (method === 'GET' && pathOnly === '/api/dishes/random') {
    return handleGetRandomDish(query) as T;
  }

  // GET /api/dishes/:slug
  const slugMatch = matchPath(pathOnly, '/api/dishes/:slug');
  if (method === 'GET' && slugMatch) {
    const result = handleGetDishBySlug(slugMatch.slug!);
    if (!result) return null;
    return result as T;
  }

  // GET /api/dishes
  if (method === 'GET' && pathOnly === '/api/dishes') {
    return handleGetDishes(query) as T;
  }

  // POST /api/rooms
  if (method === 'POST' && pathOnly === '/api/rooms') {
    return handleCreateRoom((body ?? {}) as Record<string, unknown>) as T;
  }

  // GET /api/rooms/:id
  const roomMatch = matchPath(pathOnly, '/api/rooms/:id');
  if (method === 'GET' && roomMatch) {
    return handleGetRoom(roomMatch.id!) as T;
  }

  // POST /api/rooms/:id/join
  const joinMatch = matchPath(pathOnly, '/api/rooms/:id/join');
  if (method === 'POST' && joinMatch) {
    return handleJoinRoom(joinMatch.id!, (body ?? {}) as Record<string, unknown>) as T;
  }

  return null;
}
