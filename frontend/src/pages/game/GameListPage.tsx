/**
 * Game List Page Component
 * Displays available games for users to select
 */

import { useNavigate } from 'react-router-dom';
import { BookOpen, ChevronRight } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

interface Game {
  id: string;
  name: string;
  description: string;
  icon: React.ReactNode;
  route: string;
  color: string;
}

const availableGames: Game[] = [
  {
    id: 'vocab',
    name: 'Học Từ Vựng',
    description: 'Học từ vựng qua các câu hỏi trắc nghiệm theo chủ đề hoặc cấp độ',
    icon: <BookOpen className="h-8 w-8" />,
    route: '/games/vocab/config',
    color: 'from-blue-500 to-blue-600',
  },
  // Có thể thêm các game khác sau
  // {
  //   id: 'flashcard',
  //   name: 'Flashcard',
  //   description: 'Học từ vựng bằng thẻ ghi nhớ',
  //   icon: <FileText className="h-8 w-8" />,
  //   route: '/games/flashcard',
  //   color: 'from-purple-500 to-purple-600',
  // },
];

export default function GameListPage() {
  const navigate = useNavigate();

  const handleGameSelect = (game: Game) => {
    navigate(game.route);
  };

  return (
    <div className="min-h-screen p-4 md:p-8 bg-gradient-to-br from-background to-muted/20">
      <div className="max-w-4xl mx-auto space-y-8">
        <header className="text-center space-y-2">
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">Chọn Game</h1>
          <p className="text-muted-foreground text-lg">
            Chọn một game để bắt đầu học từ vựng
          </p>
        </header>

        <main>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {availableGames.map((game) => (
              <Card
                key={game.id}
                className="cursor-pointer transition-all hover:shadow-lg hover:scale-[1.02] group"
                onClick={() => handleGameSelect(game)}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    handleGameSelect(game);
                  }
                }}
                aria-label={`Chọn game ${game.name}`}
              >
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div className={`p-3 rounded-lg bg-gradient-to-br ${game.color} text-white`}>
                      {game.icon}
                    </div>
                    <ChevronRight className="h-5 w-5 text-muted-foreground group-hover:text-foreground group-hover:translate-x-1 transition-all" />
                  </div>
                  <CardTitle className="text-xl mt-4">{game.name}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-base mb-4">
                    {game.description}
                  </CardDescription>
                  <Button className="w-full" onClick={() => handleGameSelect(game)}>
                    Bắt Đầu
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </main>
      </div>
    </div>
  );
}

