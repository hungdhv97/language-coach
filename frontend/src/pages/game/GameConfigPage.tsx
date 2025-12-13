/**
 * Game Configuration Page Component
 * Allows users to configure game session settings
 */

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { dictionaryQueries } from '@/entities/dictionary/api/dictionary.queries';
import { gameMutations } from '@/features/game/api/game.mutations';
import type { Language, Topic, Level } from '@/entities/dictionary/model/dictionary.types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertCircle } from 'lucide-react';

export default function GameConfigPage() {
  const navigate = useNavigate();
  const [sourceLanguageId, setSourceLanguageId] = useState<number | ''>('');
  const [targetLanguageId, setTargetLanguageId] = useState<number | ''>('');
  const [mode, setMode] = useState<'topic' | 'level' | ''>('');
  const [topicId, setTopicId] = useState<number | ''>('');
  const [levelId, setLevelId] = useState<number | ''>('');
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Fetch reference data
  const { data: languages = [], isLoading: languagesLoading } = dictionaryQueries.useLanguages();
  const { data: topics = [], isLoading: topicsLoading } = dictionaryQueries.useTopics();
  const { data: levels = [], isLoading: levelsLoading } = dictionaryQueries.useLevels(
    sourceLanguageId ? Number(sourceLanguageId) : undefined
  );

  // Create session mutation
  const createSessionMutation = gameMutations.useCreateSession();

  // Validation
  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    // Source and target languages must be different (FR-010)
    if (sourceLanguageId && targetLanguageId && sourceLanguageId === targetLanguageId) {
      newErrors.languages = 'Ngôn ngữ nguồn và ngôn ngữ đích phải khác nhau';
    }

    // Mode is required
    if (!mode) {
      newErrors.mode = 'Vui lòng chọn chế độ chơi';
    }

    // Topic XOR Level required (FR-011)
    if (mode === 'topic') {
      if (!topicId) {
        newErrors.topic = 'Vui lòng chọn chủ đề';
      }
      if (levelId) {
        newErrors.level = 'Không thể chọn cả chủ đề và cấp độ cùng lúc';
      }
    } else if (mode === 'level') {
      if (!levelId) {
        newErrors.level = 'Vui lòng chọn cấp độ';
      }
      if (topicId) {
        newErrors.topic = 'Không thể chọn cả chủ đề và cấp độ cùng lúc';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    try {
      const session = await createSessionMutation.mutateAsync({
        source_language_id: Number(sourceLanguageId),
        target_language_id: Number(targetLanguageId),
        mode: mode as 'topic' | 'level',
        topic_id: mode === 'topic' ? Number(topicId) : undefined,
        level_id: mode === 'level' ? Number(levelId) : undefined,
      });

      // Navigate to game play page (will be implemented in Phase 6)
      navigate(`/games/vocab/play/${session.id}`);
    } catch (error: unknown) {
      const apiError = error as { code?: string; message?: string };
      if (apiError.code === 'INSUFFICIENT_WORDS') {
        setErrors({ submit: apiError.message || 'Không đủ từ vựng để tạo game session' });
      } else if (apiError.code === 'VALIDATION_ERROR') {
        setErrors({ submit: apiError.message || 'Dữ liệu không hợp lệ' });
      } else {
        setErrors({ submit: 'Không thể tạo game session. Vui lòng thử lại' });
      }
    }
  };

  // Reset level when source language changes
  useEffect(() => {
    if (mode === 'level') {
      setLevelId('');
    }
  }, [sourceLanguageId, mode]);

  // Reset topic/level when mode changes
  useEffect(() => {
    setTopicId('');
    setLevelId('');
    setErrors({});
  }, [mode]);

  return (
    <div className="min-h-screen p-4 md:p-8 bg-gradient-to-br from-background to-muted/20">
      <div className="max-w-2xl mx-auto space-y-6">
        <header className="text-center space-y-2">
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">Cấu Hình Game</h1>
          <p className="text-muted-foreground text-lg">
            Chọn ngôn ngữ và chế độ chơi để bắt đầu
          </p>
        </header>

        <main>
          <form onSubmit={handleSubmit}>
            <Card>
              <CardHeader>
                <CardTitle>Ngôn Ngữ</CardTitle>
                <CardDescription>Chọn ngôn ngữ nguồn và ngôn ngữ đích</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="source-language">Ngôn Ngữ Nguồn</Label>
                    <Select
                      value={sourceLanguageId ? String(sourceLanguageId) : undefined}
                      onValueChange={(value) => setSourceLanguageId(value ? Number(value) : '')}
                      disabled={languagesLoading}
                      required
                    >
                      <SelectTrigger id="source-language" className={errors.languages ? 'border-destructive' : ''}>
                        <SelectValue placeholder="Chọn ngôn ngữ nguồn" />
                      </SelectTrigger>
                      <SelectContent>
                        {languages.map((lang: Language) => (
                          <SelectItem key={lang.id} value={String(lang.id)}>
                            {lang.name} ({lang.code})
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="target-language">Ngôn Ngữ Đích</Label>
                    <Select
                      value={targetLanguageId ? String(targetLanguageId) : undefined}
                      onValueChange={(value) => setTargetLanguageId(value ? Number(value) : '')}
                      disabled={languagesLoading}
                      required
                    >
                      <SelectTrigger id="target-language" className={errors.languages ? 'border-destructive' : ''}>
                        <SelectValue placeholder="Chọn ngôn ngữ đích" />
                      </SelectTrigger>
                      <SelectContent>
                        {languages.map((lang: Language) => (
                          <SelectItem key={lang.id} value={String(lang.id)}>
                            {lang.name} ({lang.code})
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                {errors.languages && (
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>{errors.languages}</AlertDescription>
                  </Alert>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Chế Độ Chơi</CardTitle>
                <CardDescription>Chọn cách bạn muốn học từ vựng</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <Button
                    type="button"
                    variant={mode === 'topic' ? 'default' : 'outline'}
                    className="h-auto py-6 flex flex-col gap-2"
                    onClick={() => setMode('topic')}
                  >
                    <span className="text-2xl">📚</span>
                    <span>Theo Chủ Đề</span>
                  </Button>
                  <Button
                    type="button"
                    variant={mode === 'level' ? 'default' : 'outline'}
                    className="h-auto py-6 flex flex-col gap-2"
                    onClick={() => setMode('level')}
                  >
                    <span className="text-2xl">📊</span>
                    <span>Theo Cấp Độ</span>
                  </Button>
                </div>
                {errors.mode && (
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>{errors.mode}</AlertDescription>
                  </Alert>
                )}
              </CardContent>
            </Card>

            {/* Topic Selection (when mode is topic) */}
            {mode === 'topic' && (
              <Card>
                <CardHeader>
                  <CardTitle>Chọn Chủ Đề</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="topic">Chủ Đề</Label>
                    <Select
                      value={topicId ? String(topicId) : undefined}
                      onValueChange={(value) => setTopicId(value ? Number(value) : '')}
                      disabled={topicsLoading}
                      required
                    >
                      <SelectTrigger id="topic" className={errors.topic ? 'border-destructive' : ''}>
                        <SelectValue placeholder="Chọn chủ đề" />
                      </SelectTrigger>
                      <SelectContent>
                        {topics.map((topic: Topic) => (
                          <SelectItem key={topic.id} value={String(topic.id)}>
                            {topic.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  {errors.topic && (
                    <Alert variant="destructive">
                      <AlertCircle className="h-4 w-4" />
                      <AlertDescription>{errors.topic}</AlertDescription>
                    </Alert>
                  )}
                </CardContent>
              </Card>
            )}

            {/* Level Selection (when mode is level) */}
            {mode === 'level' && (
              <Card>
                <CardHeader>
                  <CardTitle>Chọn Cấp Độ</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="level">Cấp Độ</Label>
                    <Select
                      value={levelId ? String(levelId) : undefined}
                      onValueChange={(value) => setLevelId(value ? Number(value) : '')}
                      disabled={levelsLoading || !sourceLanguageId}
                      required
                    >
                      <SelectTrigger id="level" className={errors.level ? 'border-destructive' : ''}>
                        <SelectValue placeholder={!sourceLanguageId ? 'Vui lòng chọn ngôn ngữ nguồn trước' : 'Chọn cấp độ'} />
                      </SelectTrigger>
                      <SelectContent>
                        {levels.map((level: Level) => (
                          <SelectItem key={level.id} value={String(level.id)}>
                            {level.name} {level.description && `- ${level.description}`}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  {errors.level && (
                    <Alert variant="destructive">
                      <AlertCircle className="h-4 w-4" />
                      <AlertDescription>{errors.level}</AlertDescription>
                    </Alert>
                  )}
                </CardContent>
              </Card>
            )}

            {/* Submit Error */}
            {errors.submit && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Lỗi</AlertTitle>
                <AlertDescription>{errors.submit}</AlertDescription>
              </Alert>
            )}

            {/* Submit Button */}
            <div className="flex gap-4 justify-end">
              <Button
                type="button"
                variant="outline"
                onClick={() => navigate('/games')}
              >
                Quay Lại
              </Button>
              <Button
                type="submit"
                disabled={createSessionMutation.isPending}
              >
                {createSessionMutation.isPending ? 'Đang tạo...' : 'Bắt Đầu Chơi'}
              </Button>
            </div>
          </form>
        </main>
      </div>
    </div>
  );
}

