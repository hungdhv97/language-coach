/**
 * User domain types
 */

export interface User {
  id: number;
  email?: string;
  username?: string;
  created_at: string;
  updated_at: string;
  is_active: boolean;
}

export interface UserProfile {
  user_id: number;
  display_name?: string;
  avatar_url?: string;
  birth_day?: string; // YYYY-MM-DD format
  bio?: string;
  created_at: string;
  updated_at: string;
}

export interface RegisterRequest {
  display_name?: string;
  email?: string;
  username?: string;
  password: string;
}

export interface RegisterResponse {
  user_id: number;
  email?: string;
  username?: string;
}

export interface LoginRequest {
  email?: string;
  username?: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user_id: number;
  email?: string;
  username?: string;
}

export interface UpdateProfileRequest {
  display_name?: string;
  avatar_url?: string;
  birth_day?: string; // YYYY-MM-DD format
  bio?: string;
}

export interface UpdateProfileResponse {
  user_id: number;
  display_name?: string;
  avatar_url?: string;
  birth_day?: string;
  bio?: string;
}
