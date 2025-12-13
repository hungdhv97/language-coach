/**
 * Dictionary Search Component
 */

import { useState, useEffect } from 'react';
import { useDebounce } from '@/hooks/useDebounce';
import { dictionaryQueries } from '@/entities/dictionary/api/dictionary.queries';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { cn } from '@/lib/utils';
import { WordDetail } from './WordDetail';
import type { Word, Language } from '@/entities/dictionary/model/dictionary.types';

export function DictionarySearch() {
  const [searchQuery, setSearchQuery] = useState('');
  const [languageId, setLanguageId] = useState<number | ''>('');
  const [selectedWordId, setSelectedWordId] = useState<number | null>(null);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const debouncedQuery = useDebounce(searchQuery, 500);
  const maxResults = 5;

  // Fetch languages
  const { data: languages = [], isLoading: languagesLoading } = dictionaryQueries.useLanguages();

  // Set default language to English when languages are loaded
  useEffect(() => {
    if (languages.length > 0 && !languageId) {
      const englishLang = languages.find((lang: Language) => lang.code === 'en');
      if (englishLang) {
        setLanguageId(englishLang.id);
      }
    }
  }, [languages, languageId]);

  const {
    data: searchResults,
    isLoading,
    isError,
  } = dictionaryQueries.useSearchWords(
    debouncedQuery,
    languageId ? Number(languageId) : 0,
    maxResults,
    0,
    debouncedQuery.length > 0 && !!languageId
  );

  // Reset selected word when search query changes
  useEffect(() => {
    setSelectedWordId(null);
  }, [debouncedQuery, languageId]);

  // Show dropdown when typing
  useEffect(() => {
    if (languageId && debouncedQuery.length > 0) {
      setIsDropdownOpen(true);
    } else {
      setIsDropdownOpen(false);
    }
  }, [debouncedQuery, languageId]);

  const handleWordClick = (wordId: number) => {
    setSelectedWordId(wordId);
    setIsDropdownOpen(false);
  };

  const showDropdown = 
    isDropdownOpen &&
    languageId && 
    debouncedQuery.length > 0 && 
    !isError && 
    (isLoading || searchResults !== undefined);

  return (
    <div className="space-y-4">
      {/* Language selector and search bar */}
      <div className="flex flex-col sm:flex-row gap-2">
        {/* Search input with dropdown */}
        <div className="flex-1 relative">
          <input
            type="text"
            placeholder={languageId ? "Tìm kiếm từ..." : "Vui lòng chọn ngôn ngữ trước"}
            value={searchQuery}
            onChange={(e) => {
              setSearchQuery(e.target.value);
              setIsDropdownOpen(true);
            }}
            onFocus={() => {
              if (languageId && debouncedQuery.length > 0) {
                setIsDropdownOpen(true);
              }
            }}
            disabled={!languageId}
            className={cn(
              "w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm",
              "ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium",
              "placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2",
              "focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            )}
          />
          
          {/* Dropdown results */}
          {showDropdown && (
            <div className="absolute z-50 w-full mt-1 bg-popover text-popover-foreground rounded-md border shadow-md max-h-[300px] overflow-y-auto">
              {isLoading ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  Đang tìm kiếm...
                </div>
              ) : searchResults && searchResults.words.length === 0 ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  Không tìm thấy từ nào phù hợp với "{debouncedQuery}"
                </div>
              ) : searchResults && searchResults.words.length > 0 ? (
                <div className="p-1">
                  {searchResults.words.map((word: Word) => (
                    <div
                      key={word.id}
                      onClick={() => handleWordClick(word.id)}
                      className={cn(
                        "cursor-pointer rounded-sm px-2 py-1.5 text-sm transition-colors",
                        "hover:bg-accent hover:text-accent-foreground",
                        selectedWordId === word.id && "bg-accent"
                      )}
                    >
                      <div className="font-medium">{word.lemma}</div>
                      {word.romanization && (
                        <div className="text-xs text-muted-foreground mt-0.5">
                          {word.romanization}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              ) : null}
            </div>
          )}
        </div>
        
        {/* Language selector - above on small screens, right on larger screens */}
        <div className="w-full sm:w-auto sm:min-w-[140px]">
          <Select
            value={languageId ? String(languageId) : undefined}
            onValueChange={(value) => setLanguageId(value ? Number(value) : '')}
            disabled={languagesLoading}
          >
            <SelectTrigger id="language-select" className="w-full">
              <SelectValue placeholder="Chọn ngôn ngữ" />
            </SelectTrigger>
            <SelectContent align="start" side="bottom">
              {languages.map((lang: Language) => (
                <SelectItem key={lang.id} value={String(lang.id)}>
                  {lang.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

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

      {isError && (
        <div className="text-center py-8 text-destructive">
          Có lỗi xảy ra khi tìm kiếm. Vui lòng thử lại.
        </div>
      )}

      {/* Word Detail - shown below when a word is selected */}
      {selectedWordId && (
        <div className="mt-4">
          <WordDetail wordId={selectedWordId} />
        </div>
      )}
    </div>
  );
}

