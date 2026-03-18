/**
 * Mock API handler — returns mock data when the real backend is unavailable.
 * Parses request paths and returns appropriate responses matching the real API contract.
 *
 * All responses are wrapped in { data } or { data, pagination } to match the BE format.
 * The api-client's unwrapData/getList will handle unwrapping.
 */

import type { Dish, DishDetail } from '@features/dish/types';
import type { Pagination } from '@shared/lib/api-client';
import type { VoteRoom } from '@features/vote/types';
import { MOCK_CATEGORIES, MOCK_DISHES, MOCK_ROOM, MOCK_TAGS } from './mock-data';
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
  const idx = url.indexOf('/api/');
  return idx >= 0 ? url.slice(idx) : url;
}

function makePagination(total: number, page: number, pageSize: number): Pagination {
  return { page, pageSize, total, totalPages: Math.ceil(total / pageSize) };
}

// ---- Shared filter logic ----

function filterDishes(dishes: Dish[], query: URLSearchParams): Dish[] {
  let result = dishes;

  const dishType = query.get('dish_type');
  if (dishType) {
    result = result.filter((d) => d.dishTypeId === dishType);
  }

  const region = query.get('region');
  if (region) {
    result = result.filter((d) => d.regionId === region);
  }

  const difficulty = query.get('difficulty');
  if (difficulty) {
    result = result.filter((d) => d.difficulty === difficulty);
  }

  const maxCookTime = query.get('max_cook_time');
  if (maxCookTime) {
    const max = Number(maxCookTime);
    result = result.filter((d) => (d.cookTime ?? 0) <= max);
  }

  return result;
}

// ---- Route handlers ----

function handleGetRecipes(query: URLSearchParams): { data: Dish[]; pagination: Pagination } {
  const dishes = filterDishes(MOCK_DISHES, query);

  const page = Number(query.get('page') ?? 1);
  const pageSize = Number(query.get('pageSize') ?? 20);
  const start = (page - 1) * pageSize;
  const paged = dishes.slice(start, start + pageSize);

  return { data: paged, pagination: makePagination(dishes.length, page, pageSize) };
}

function handleSearchRecipes(query: URLSearchParams): { data: Dish[]; pagination: Pagination } {
  const q = (query.get('q') ?? '').toLowerCase();
  let matched: Dish[] = MOCK_DISHES;

  if (q) {
    matched = matched.filter(
      (d) => d.name.toLowerCase().includes(q) || (d.description ?? '').toLowerCase().includes(q),
    );
  }

  const dishes = filterDishes(matched, query);

  const page = Number(query.get('page') ?? 1);
  const pageSize = Number(query.get('pageSize') ?? 20);
  const start = (page - 1) * pageSize;
  const paged = dishes.slice(start, start + pageSize);

  return { data: paged, pagination: makePagination(dishes.length, page, pageSize) };
}

function handleGetRecipeBySlug(slug: string): { data: DishDetail } | null {
  const dish = MOCK_DISHES.find((d) => d.slug === slug);
  return dish ? { data: dish } : null;
}

function handleGetRandomRecipe(query: URLSearchParams): { data: Dish } | null {
  let dishes: Dish[] = filterDishes(MOCK_DISHES, query);

  const excludeIds = query.get('exclude_ids');
  if (excludeIds) {
    const ids = excludeIds.split(',');
    dishes = dishes.filter((d) => !ids.includes(d.id));
  }

  if (dishes.length === 0) {
    const fallback = MOCK_DISHES[0];
    return fallback ? { data: fallback } : null;
  }

  const idx = Math.floor(Math.random() * dishes.length);
  const dish = dishes[idx];
  return dish ? { data: dish } : null;
}

function handleGetCategories(query: URLSearchParams): { data: typeof MOCK_CATEGORIES } {
  const type = query.get('type');
  if (type) {
    return { data: MOCK_CATEGORIES.filter((c) => c.type === type) };
  }
  return { data: MOCK_CATEGORIES };
}

function handleGetTags(): { data: typeof MOCK_TAGS } {
  return { data: MOCK_TAGS };
}

// ---- Vote room handlers ----

function handleCreateRoom(body: Record<string, unknown>): { data: VoteRoom } {
  return {
    data: {
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
    },
  };
}

function handleGetRoom(roomId: string): { data: VoteRoom } {
  return { data: { ...MOCK_ROOM, id: roomId } };
}

function handleJoinRoom(roomId: string, body: Record<string, unknown>): { data: VoteRoom } {
  const name = String(body.name ?? 'Player');
  const colors = ['#FF6B35', '#E63946', '#FFB703', '#2EC4B6', '#9B5DE5'];
  return {
    data: {
      ...MOCK_ROOM,
      id: roomId,
      participants: [
        ...MOCK_ROOM.participants,
        {
          name,
          avatarColor: colors[Math.floor(Math.random() * colors.length)]!,
          hasVoted: false,
        },
      ],
    },
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

  // GET /api/v1/recipes/random
  if (method === 'GET' && pathOnly === '/api/v1/recipes/random') {
    return handleGetRandomRecipe(query) as T;
  }

  // GET /api/v1/recipes/search
  if (method === 'GET' && pathOnly === '/api/v1/recipes/search') {
    return handleSearchRecipes(query) as T;
  }

  // GET /api/v1/recipes/:slug
  const slugMatch = matchPath(pathOnly, '/api/v1/recipes/:slug');
  if (method === 'GET' && slugMatch) {
    const result = handleGetRecipeBySlug(slugMatch.slug!);
    if (!result) return null;
    return result as T;
  }

  // GET /api/v1/recipes
  if (method === 'GET' && pathOnly === '/api/v1/recipes') {
    return handleGetRecipes(query) as T;
  }

  // GET /api/v1/categories
  if (method === 'GET' && pathOnly === '/api/v1/categories') {
    return handleGetCategories(query) as T;
  }

  // GET /api/v1/tags
  if (method === 'GET' && pathOnly === '/api/v1/tags') {
    return handleGetTags() as T;
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
