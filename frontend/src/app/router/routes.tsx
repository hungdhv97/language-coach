/**
 * Application Routes Configuration
 */

import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from '../../components/layout/Layout';
import LandingPage from '../../pages/LandingPage';
import GameListPage from '../../pages/game/GameListPage';
import GameConfigPage from '../../pages/game/GameConfigPage';
import GamePlayPage from '../../pages/game/GamePlayPage';
import GameStatisticsPage from '../../pages/game/GameStatisticsPage';
import DictionaryLookupPage from '../../pages/dictionary/DictionaryLookupPage';
import WordDetailPage from '../../pages/dictionary/WordDetailPage';
import LoginPage from '../../pages/auth/LoginPage';
import RegisterPage from '../../pages/auth/RegisterPage';
import ProfilePage from '../../pages/auth/ProfilePage';
import { useAuthStore } from '../../shared/store/useAuthStore';

// Protected Route wrapper component
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  if (!isAuthenticated) {
    return <Navigate to="/auth/login" replace />;
  }
  return <>{children}</>;
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <LandingPage />,
      },
      {
        path: '/auth/login',
        element: <LoginPage />,
      },
      {
        path: '/auth/register',
        element: <RegisterPage />,
      },
      {
        path: '/auth/profile',
        element: (
          <ProtectedRoute>
            <ProfilePage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/games',
        element: (
          <ProtectedRoute>
            <GameListPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/games/vocab/config',
        element: (
          <ProtectedRoute>
            <GameConfigPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/games/vocab/play/:sessionId',
        element: (
          <ProtectedRoute>
            <GamePlayPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/games/vocab/statistics/:sessionId',
        element: (
          <ProtectedRoute>
            <GameStatisticsPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/dictionary',
        element: <DictionaryLookupPage />,
      },
      {
        path: '/dictionary/words/:wordId',
        element: <WordDetailPage />,
      },
    ],
  },
]);
