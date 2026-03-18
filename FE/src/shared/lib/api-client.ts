import { API_BASE_URL } from '@shared/constants';
import { mockApiHandler } from './mock-api';

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public override readonly message: string,
    public readonly code?: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

interface RequestOptions {
  headers?: Record<string, string>;
  signal?: AbortSignal;
}

interface ApiClient {
  get<T>(path: string, options?: RequestOptions): Promise<T>;
  post<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T>;
  put<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T>;
  delete<T>(path: string, options?: RequestOptions): Promise<T>;
}

const MAX_RETRIES = 1;
const RETRY_DELAY_MS = 1000;

function isServerError(status: number): boolean {
  return status >= 500 && status < 600;
}

async function delay(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function parseResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get('content-type');
  if (contentType?.includes('application/json')) {
    return (await response.json()) as T;
  }
  return (await response.text()) as unknown as T;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (response.ok) {
    return parseResponse<T>(response);
  }

  let errorMessage = response.statusText;
  let errorCode: string | undefined;

  try {
    const errorBody = (await response.json()) as {
      message?: string;
      code?: string;
    };
    errorMessage = errorBody.message ?? errorMessage;
    errorCode = errorBody.code;
  } catch {
    // Body is not JSON — use statusText
  }

  throw new ApiError(response.status, errorMessage, errorCode);
}

function shouldFallbackToMock(error: unknown): boolean {
  // Network errors (backend completely unreachable)
  if (
    error instanceof TypeError ||
    (error instanceof Error &&
      (error.message.includes('fetch') || error.message.includes('ECONNREFUSED')))
  ) {
    return true;
  }
  // API errors indicating the backend doesn't have this endpoint yet
  if (error instanceof ApiError && (error.status === 404 || error.status >= 500)) {
    return true;
  }
  return false;
}

async function fetchWithRetry<T>(url: string, init: RequestInit, body?: unknown): Promise<T> {
  let lastError: unknown;

  for (let attempt = 0; attempt <= MAX_RETRIES; attempt++) {
    try {
      const response = await fetch(url, init);

      if (isServerError(response.status) && attempt < MAX_RETRIES) {
        await delay(RETRY_DELAY_MS);
        continue;
      }

      return await handleResponse<T>(response);
    } catch (error) {
      lastError = error;

      // On ApiError, skip retries but let mock fallback handle it below
      if (error instanceof ApiError) {
        break;
      }

      if (attempt < MAX_RETRIES) {
        await delay(RETRY_DELAY_MS);
        continue;
      }
    }
  }

  // Backend unreachable or endpoint missing — fall back to mock data
  if (shouldFallbackToMock(lastError)) {
    const method = init.method ?? 'GET';
    const mockResult = await mockApiHandler<T>(method, url, body);
    if (mockResult !== null) {
      return mockResult;
    }
  }

  throw lastError;
}

function getBaseUrl(configuredUrl: string): string {
  // On the server, always use the configured URL directly
  if (typeof window === 'undefined') {
    return configuredUrl;
  }
  // On the client, use the public URL from env
  return configuredUrl;
}

export function createApiClient(baseUrl: string): ApiClient {
  const resolvedBaseUrl = getBaseUrl(baseUrl);

  function buildHeaders(customHeaders?: Record<string, string>): Record<string, string> {
    return {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      ...customHeaders,
    };
  }

  return {
    async get<T>(path: string, options?: RequestOptions): Promise<T> {
      return fetchWithRetry<T>(
        `${resolvedBaseUrl}${path}`,
        {
          method: 'GET',
          headers: buildHeaders(options?.headers),
          signal: options?.signal,
        },
        undefined,
      );
    },

    async post<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
      return fetchWithRetry<T>(
        `${resolvedBaseUrl}${path}`,
        {
          method: 'POST',
          headers: buildHeaders(options?.headers),
          body: body !== undefined ? JSON.stringify(body) : undefined,
          signal: options?.signal,
        },
        body,
      );
    },

    async put<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
      return fetchWithRetry<T>(
        `${resolvedBaseUrl}${path}`,
        {
          method: 'PUT',
          headers: buildHeaders(options?.headers),
          body: body !== undefined ? JSON.stringify(body) : undefined,
          signal: options?.signal,
        },
        body,
      );
    },

    async delete<T>(path: string, options?: RequestOptions): Promise<T> {
      return fetchWithRetry<T>(
        `${resolvedBaseUrl}${path}`,
        {
          method: 'DELETE',
          headers: buildHeaders(options?.headers),
          signal: options?.signal,
        },
        undefined,
      );
    },
  };
}

export const apiClient = createApiClient(API_BASE_URL);
