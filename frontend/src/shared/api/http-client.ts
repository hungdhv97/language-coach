/**
 * HTTP Client wrapper using fetch API
 */

import { API_CONFIG } from './config';

export interface ApiError {
  code: string;
  message: string;
  details?: unknown;
}

export interface RequestConfig extends RequestInit {
  timeout?: number;
}

// Custom error class for token expiration
export class TokenExpiredError extends Error {
  code = 'TOKEN_EXPIRED';
  constructor(message: string = 'Token đã hết hạn') {
    super(message);
    this.name = 'TokenExpiredError';
  }
}

// Callback type for token expiration handling
type TokenExpiredCallback = () => void;

class HttpClient {
  private baseURL: string;
  private defaultTimeout: number;
  private authToken: string | null = null;
  private onTokenExpired: TokenExpiredCallback | null = null;

  constructor(baseURL: string, defaultTimeout: number = 30000) {
    this.baseURL = baseURL;
    this.defaultTimeout = defaultTimeout;
    // Load token from localStorage on initialization
    this.authToken = localStorage.getItem('auth_token');
  }

  setAuthToken(token: string | null) {
    this.authToken = token;
    if (token) {
      localStorage.setItem('auth_token', token);
    } else {
      localStorage.removeItem('auth_token');
    }
  }

  getAuthToken(): string | null {
    return this.authToken;
  }

  /**
   * Set callback to handle token expiration
   * This will be called automatically when TOKEN_EXPIRED error is detected
   */
  setTokenExpiredCallback(callback: TokenExpiredCallback) {
    this.onTokenExpired = callback;
  }

  private async request<T>(
    endpoint: string,
    config: RequestConfig = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const controller = new AbortController();
    const timeoutId = setTimeout(
      () => controller.abort(),
      config.timeout || this.defaultTimeout
    );

    const headers: Record<string, string> = {
      ...(API_CONFIG.headers as Record<string, string>),
      ...(config.headers as Record<string, string>),
    };

    // Add auth token if available
    if (this.authToken) {
      headers['Authorization'] = `Bearer ${this.authToken}`;
    }

    try {
      const response = await fetch(url, {
        ...config,
        signal: controller.signal,
        headers,
      });

      clearTimeout(timeoutId);

      const data = await response.json();

      if (!response.ok) {
        // Handle error response format
        let apiError: ApiError;
        if (data.error) {
          apiError = {
            code: data.error.code || 'UNKNOWN_ERROR',
            message: data.error.message || 'An error occurred',
            details: data.error.details,
          };
        } else {
          apiError = data as ApiError;
        }

        // Handle token expiration
        // Backend returns code "TOKEN_EXPIRED" when token is expired
        if (apiError.code === 'TOKEN_EXPIRED') {
          // Clear token immediately
          this.setAuthToken(null);
          
          // Call token expired callback if set
          if (this.onTokenExpired) {
            this.onTokenExpired();
          }
          
          // Throw TokenExpiredError
          throw new TokenExpiredError(apiError.message || 'Token đã hết hạn');
        }

        // Throw an Error instance with the message, and attach ApiError properties
        const error = new Error(apiError.message || 'An error occurred');
        Object.assign(error, { code: apiError.code, details: apiError.details });
        throw error;
      }

      // Backend returns { success: true, data: T } format
      // Return the whole response, let endpoints extract .data if needed
      return data as T;
    } catch (error) {
      clearTimeout(timeoutId);
      if (error instanceof Error) {
        throw error;
      }
      // If it's an ApiError object (plain object), convert it to Error
      if (error && typeof error === 'object' && 'message' in error) {
        const apiError = error as ApiError;
        const err = new Error(apiError.message);
        Object.assign(err, { code: apiError.code, details: apiError.details });
        throw err;
      }
      throw new Error('Unknown error occurred');
    }
  }

  async get<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, { ...config, method: 'GET' });
  }

  async post<T>(
    endpoint: string,
    data?: unknown,
    config?: RequestConfig
  ): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(
    endpoint: string,
    data?: unknown,
    config?: RequestConfig
  ): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async patch<T>(
    endpoint: string,
    data?: unknown,
    config?: RequestConfig
  ): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, { ...config, method: 'DELETE' });
  }
}

export const httpClient = new HttpClient(
  API_CONFIG.baseURL,
  API_CONFIG.timeout
);

