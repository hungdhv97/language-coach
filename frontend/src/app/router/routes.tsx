/**
 * Application Routes Configuration
 */

import { createBrowserRouter } from 'react-router-dom';
import LandingPage from '../../pages/LandingPage';
import GameListPage from '../../pages/game/GameListPage';
import GameConfigPage from '../../pages/game/GameConfigPage';
import GamePlayPage from '../../pages/game/GamePlayPage';
import GameStatisticsPage from '../../pages/game/GameStatisticsPage';
import DictionaryLookupPage from '../../pages/dictionary/DictionaryLookupPage';
import WordDetailPage from '../../pages/dictionary/WordDetailPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <LandingPage />,
  },
  {
    path: '/games',
    element: <GameListPage />,
  },
  {
    path: '/games/vocab/config',
    element: <GameConfigPage />,
  },
  {
    path: '/games/vocab/play/:sessionId',
    element: <GamePlayPage />,
  },
  {
    path: '/games/vocab/statistics/:sessionId',
    element: <GameStatisticsPage />,
  },
  {
    path: '/dictionary',
    element: <DictionaryLookupPage />,
  },
  {
    path: '/dictionary/words/:wordId',
    element: <WordDetailPage />,
  },
]);
