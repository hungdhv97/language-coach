/**
 * Landing Page Component
 * Displays two prominent action buttons for navigation
 */

import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Gamepad2, BookOpen } from 'lucide-react';
import { useAuthStore } from '@/shared/store/useAuthStore';

export default function LandingPage() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();

  const handlePlayGame = () => {
    if (!isAuthenticated) {
      navigate('/auth/login');
    } else {
      navigate('/games');
    }
  };

  const handleDictionaryLookup = () => {
    navigate('/dictionary');
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-background to-muted/20">
      <div className="w-full max-w-4xl space-y-8">
        <header className="text-center space-y-4">
          <h1 className="text-4xl md:text-5xl font-bold tracking-tight">English Coach</h1>
          <p className="text-lg text-muted-foreground">
            Học từ vựng đa ngôn ngữ một cách hiệu quả
          </p>
        </header>

        <main className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card className="cursor-pointer transition-all hover:shadow-lg hover:scale-105" onClick={handlePlayGame}>
            <CardHeader>
              <div className="flex items-center gap-3 mb-2">
                <div className="p-2 rounded-lg bg-primary/10">
                  <Gamepad2 className="h-6 w-6 text-primary" />
                </div>
                <CardTitle className="text-2xl">Chơi Game</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-base">
                Học từ vựng qua các trò chơi thú vị
              </CardDescription>
              <Button className="w-full mt-4" size="lg">
                Bắt Đầu Chơi
              </Button>
            </CardContent>
          </Card>

          <Card className="cursor-pointer transition-all hover:shadow-lg hover:scale-105" onClick={handleDictionaryLookup}>
            <CardHeader>
              <div className="flex items-center gap-3 mb-2">
                <div className="p-2 rounded-lg bg-primary/10">
                  <BookOpen className="h-6 w-6 text-primary" />
                </div>
                <CardTitle className="text-2xl">Tra Cứu Từ Điển</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-base">
                Tìm kiếm và học từ vựng đa ngôn ngữ
              </CardDescription>
              <Button className="w-full mt-4" size="lg" variant="secondary">
                Mở Từ Điển
              </Button>
            </CardContent>
          </Card>
        </main>
      </div>
    </div>
  );
}

