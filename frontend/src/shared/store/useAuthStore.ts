/**
 * Authentication store using Zustand
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { httpClient } from '../api/http-client';
import type { User } from '@/entities/user/model/user.types';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  setUser: (user: User | null) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (token: string, user: User) => {
        httpClient.setAuthToken(token);
        set({ token, user, isAuthenticated: true });
      },
      logout: () => {
        httpClient.setAuthToken(null);
        set({ token: null, user: null, isAuthenticated: false });
      },
      setUser: (user: User | null) => {
        set({ user });
      },
    }),
    {
      name: 'auth-storage',
      onRehydrateStorage: () => (state) => {
        // Restore auth token to httpClient on rehydration
        if (state?.token) {
          httpClient.setAuthToken(state.token);
        }
      },
    }
  )
);
