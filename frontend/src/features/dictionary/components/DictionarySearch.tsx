/**
 * Dictionary Search Component
 */

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDebounce } from '@/hooks/useDebounce';
import { dictionaryQueries } from '@/entities/dictionary/api/dictionary.queries';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';
import type { Word, Language } from '@/entities/dictionary/model/dictionary.types';

export function DictionarySearch() {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [languageId, setLanguageId] = useState<number | ''>('');
  const [currentPage, setCurrentPage] = useState(0);
  const pageSize = 20;
  const debouncedQuery = useDebounce(searchQuery, 500);

  // Fetch languages
  const { data: languages = [], isLoading: languagesLoading } = dictionaryQueries.useLanguages();

  const {
    data: searchResults,
    isLoading,
    isError,
  } = dictionaryQueries.useSearchWords(
    debouncedQuery,
    languageId ? Number(languageId) : 0,
    pageSize,
    currentPage * pageSize,
    debouncedQuery.length > 0 && !!languageId
  );

  const totalPages = searchResults ? Math.ceil(searchResults.total / pageSize) : 0;

  // Reset page when query or language changes
  useEffect(() => {
    setCurrentPage(0);
  }, [debouncedQuery, languageId]);

  const handleWordClick = (wordId: number) => {
    navigate(`/dictionary/words/${wordId}`);
  };

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="language-select">Ngôn ngữ</Label>
        <Select
          value={languageId ? String(languageId) : undefined}
          onValueChange={(value) => setLanguageId(value ? Number(value) : '')}
          disabled={languagesLoading}
        >
          <SelectTrigger id="language-select">
            <SelectValue placeholder="Chọn ngôn ngữ để tra cứu" />
          </SelectTrigger>
          <SelectContent>
            {languages.map((lang: Language) => (
              <SelectItem key={lang.id} value={String(lang.id)}>
                {lang.native_name || lang.name} ({lang.code})
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="flex gap-2">
        <input
          type="text"
          placeholder={languageId ? "Tìm kiếm từ..." : "Vui lòng chọn ngôn ngữ trước"}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          disabled={!languageId}
          className={cn(
            "flex-1 h-9 rounded-md border border-input bg-background px-3 py-1 text-sm",
            "ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium",
            "placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2",
            "focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          )}
        />
      </div>

      {isLoading && debouncedQuery.length > 0 && (
        <div className="text-center py-8 text-muted-foreground">
          Đang tìm kiếm...
        </div>
      )}

      {isError && (
        <div className="text-center py-8 text-destructive">
          Có lỗi xảy ra khi tìm kiếm. Vui lòng thử lại.
        </div>
      )}

      {!languageId && (
        <div className="text-center py-8 text-muted-foreground">
          Vui lòng chọn ngôn ngữ để bắt đầu tra cứu
        </div>
      )}

      {languageId && !isLoading && !isError && debouncedQuery.length === 0 && (
        <div className="text-center py-8 text-muted-foreground">
          Nhập từ cần tìm kiếm vào ô trên
        </div>
      )}

      {!isLoading && !isError && debouncedQuery.length > 0 && searchResults && (
        <>
          {searchResults.words.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              Không tìm thấy từ nào phù hợp với "{debouncedQuery}"
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="text-sm text-muted-foreground">
                  Tìm thấy {searchResults.total} kết quả
                  {searchResults.total > 0 && (
                    <span className="ml-2">
                      (Trang {currentPage + 1} / {totalPages})
                    </span>
                  )}
                </div>
              </div>
              <div className="space-y-2">
                {searchResults.words.map((word: Word) => (
                  <Card
                    key={word.id}
                    className="cursor-pointer hover:bg-accent transition-colors"
                    onClick={() => handleWordClick(word.id)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold text-lg">{word.lemma}</h3>
                          {word.romanization && (
                            <p className="text-sm text-muted-foreground">
                              {word.romanization}
                            </p>
                          )}
                        </div>
                        {word.part_of_speech_id && (
                          <span className="text-xs text-muted-foreground">
                            POS: {word.part_of_speech_id}
                          </span>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>

              {/* Pagination Controls */}
              {totalPages > 1 && (
                <div className="flex items-center justify-center gap-2 pt-4">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setCurrentPage((prev) => Math.max(0, prev - 1))}
                    disabled={currentPage === 0 || isLoading}
                  >
                    Trước
                  </Button>
                  <span className="text-sm text-muted-foreground px-4">
                    Trang {currentPage + 1} / {totalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setCurrentPage((prev) => Math.min(totalPages - 1, prev + 1))}
                    disabled={currentPage >= totalPages - 1 || isLoading}
                  >
                    Sau
                  </Button>
                </div>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}

