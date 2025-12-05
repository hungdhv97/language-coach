/**
 * User API endpoints
 */

import { httpClient } from '@/shared/api/http-client';
import type {
  RegisterRequest,
  RegisterResponse,
  LoginRequest,
  LoginResponse,
  UserProfile,
  UpdateProfileRequest,
  UpdateProfileResponse,
} from '../model/user.types';

export interface ApiResponse<T> {
  success: boolean;
  data: T;
}

export const userEndpoints = {
  /**
   * Register a new user
   */
  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await httpClient.post<ApiResponse<RegisterResponse>>('/auth/register', data);
    return response.data;
  },

  /**
   * Login user
   */
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await httpClient.post<ApiResponse<LoginResponse>>('/auth/login', data);
    return response.data;
  },

  /**
   * Get user profile (requires authentication)
   */
  getProfile: async (): Promise<UserProfile> => {
    const response = await httpClient.get<ApiResponse<UserProfile>>('/users/profile');
    return response.data;
  },

  /**
   * Update user profile (requires authentication)
   */
  updateProfile: async (data: UpdateProfileRequest): Promise<UpdateProfileResponse> => {
    const response = await httpClient.put<ApiResponse<UpdateProfileResponse>>('/users/profile', data);
    return response.data;
  },

  /**
   * Check if email is available
   */
  checkEmailAvailability: async (email: string): Promise<{ available: boolean; exists: boolean }> => {
    const response = await httpClient.get<ApiResponse<{ available: boolean; exists: boolean }>>(
      `/auth/check-email?email=${encodeURIComponent(email)}`
    );
    return response.data;
  },

  /**
   * Check if username is available
   */
  checkUsernameAvailability: async (username: string): Promise<{ available: boolean; exists: boolean }> => {
    const response = await httpClient.get<ApiResponse<{ available: boolean; exists: boolean }>>(
      `/auth/check-username?username=${encodeURIComponent(username)}`
    );
    return response.data;
  },
};
