import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { createApiClient, ApiError } from '../api-client';

const mockFetch = vi.fn();
global.fetch = mockFetch;

function jsonResponse(data: unknown, status = 200, ok = true) {
  return {
    ok,
    status,
    statusText: status === 200 ? 'OK' : 'Error',
    headers: new Headers({ 'content-type': 'application/json' }),
    json: () => Promise.resolve(data),
    text: () => Promise.resolve(JSON.stringify(data)),
  } as unknown as Response;
}

describe('createApiClient', () => {
  const api = createApiClient('http://localhost:8080');

  beforeEach(() => {
    mockFetch.mockReset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('GET', () => {
    it('sends GET request with correct URL and headers', async () => {
      mockFetch.mockResolvedValueOnce(jsonResponse({ id: '1' }));

      const result = await api.get('/api/dishes');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dishes',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Accept: 'application/json',
          }),
        }),
      );
      expect(result).toEqual({ id: '1' });
    });

    it('passes custom headers', async () => {
      mockFetch.mockResolvedValueOnce(jsonResponse({}));

      await api.get('/test', { headers: { 'X-Custom': 'value' } });

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/test',
        expect.objectContaining({
          headers: expect.objectContaining({ 'X-Custom': 'value' }),
        }),
      );
    });
  });

  describe('POST', () => {
    it('sends POST request with JSON body', async () => {
      mockFetch.mockResolvedValueOnce(jsonResponse({ created: true }));

      const result = await api.post('/api/rooms', { name: 'test' });

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/rooms',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ name: 'test' }),
        }),
      );
      expect(result).toEqual({ created: true });
    });

    it('omits body when undefined', async () => {
      mockFetch.mockResolvedValueOnce(jsonResponse({}));

      await api.post('/api/action');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/action',
        expect.objectContaining({
          method: 'POST',
          body: undefined,
        }),
      );
    });
  });

  describe('PUT', () => {
    it('sends PUT request with JSON body', async () => {
      mockFetch.mockResolvedValueOnce(jsonResponse({ updated: true }));

      const result = await api.put('/api/dishes/1', { name: 'updated' });

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dishes/1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ name: 'updated' }),
        }),
      );
      expect(result).toEqual({ updated: true });
    });
  });

  describe('DELETE', () => {
    it('sends DELETE request', async () => {
      mockFetch.mockResolvedValueOnce(jsonResponse({ deleted: true }));

      const result = await api.delete('/api/dishes/1');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dishes/1',
        expect.objectContaining({ method: 'DELETE' }),
      );
      expect(result).toEqual({ deleted: true });
    });
  });

  describe('error handling', () => {
    it('throws ApiError for non-ok responses', async () => {
      mockFetch.mockResolvedValueOnce(
        jsonResponse({ message: 'Not found', code: 'DISH_NOT_FOUND' }, 404, false),
      );

      await expect(api.get('/api/dishes/missing')).rejects.toThrow(ApiError);

      try {
        mockFetch.mockResolvedValueOnce(
          jsonResponse({ message: 'Not found', code: 'DISH_NOT_FOUND' }, 404, false),
        );
        await api.get('/api/dishes/missing');
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        const apiError = error as ApiError;
        expect(apiError.status).toBe(404);
        expect(apiError.message).toBe('Not found');
        expect(apiError.code).toBe('DISH_NOT_FOUND');
      }
    });

    it('retries on 5xx server errors', async () => {
      mockFetch
        .mockResolvedValueOnce(jsonResponse({}, 500, false))
        .mockResolvedValueOnce(jsonResponse({ ok: true }));

      const result = await api.get('/api/test');

      expect(mockFetch).toHaveBeenCalledTimes(2);
      expect(result).toEqual({ ok: true });
    });

    it('throws after max retries on persistent server error', async () => {
      mockFetch
        .mockResolvedValueOnce(jsonResponse({ message: 'Server error' }, 500, false))
        .mockResolvedValueOnce(jsonResponse({ message: 'Server error' }, 500, false));

      await expect(api.get('/api/test')).rejects.toThrow(ApiError);
    });
  });

  describe('signal support', () => {
    it('passes abort signal to fetch', async () => {
      mockFetch.mockResolvedValueOnce(jsonResponse({}));
      const controller = new AbortController();

      await api.get('/test', { signal: controller.signal });

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/test',
        expect.objectContaining({ signal: controller.signal }),
      );
    });
  });
});
