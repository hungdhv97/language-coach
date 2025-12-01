/**
 * Word Detail Component
 * Displays detailed information about a word including senses, translations, examples, and pronunciation
 */

import { dictionaryQueries } from '@/entities/dictionary/api/dictionary.queries';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Spinner } from '@/components/ui/spinner';

interface WordDetailProps {
  wordId: number;
}

export function WordDetail({ wordId }: WordDetailProps) {
  const { data: wordDetail, isLoading, isError } = dictionaryQueries.useWordDetail(wordId);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Spinner />
      </div>
    );
  }

  if (isError || !wordDetail) {
    return (
      <div className="text-center py-12 text-destructive">
        Không thể tải thông tin từ. Vui lòng thử lại.
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Word Header */}
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">{wordDetail.word.lemma}</CardTitle>
          {wordDetail.word.romanization && (
            <p className="text-muted-foreground">{wordDetail.word.romanization}</p>
          )}
        </CardHeader>
        <CardContent>
          {/* Pronunciations */}
          {wordDetail.pronunciations.length > 0 && (
            <div className="space-y-2">
              <h3 className="font-semibold">Phát âm:</h3>
              <div className="space-y-1">
                {wordDetail.pronunciations.map((pron) => (
                  <div key={pron.id} className="text-sm">
                    {pron.ipa && <span className="font-mono">{pron.ipa}</span>}
                    {pron.phonetic && (
                      <span className="ml-2 text-muted-foreground">
                        ({pron.phonetic})
                      </span>
                    )}
                    {pron.dialect && (
                      <span className="ml-2 text-xs text-muted-foreground">
                        [{pron.dialect}]
                      </span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Senses */}
      {wordDetail.senses.length > 0 && (
        <div className="space-y-4">
          <h2 className="text-xl font-semibold">Nghĩa:</h2>
          {wordDetail.senses.map((sense) => (
            <Card key={sense.id}>
              <CardHeader>
                <CardTitle className="text-lg">
                  {sense.sense_order}. {sense.definition}
                </CardTitle>
                {sense.usage_label && (
                  <p className="text-sm text-muted-foreground">
                    {sense.usage_label}
                  </p>
                )}
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Translations */}
                {sense.translations.length > 0 && (
                  <div>
                    <h3 className="font-semibold mb-2">Bản dịch:</h3>
                    <div className="flex flex-wrap gap-2">
                      {sense.translations.map((translation) => (
                        <span
                          key={translation.id}
                          className="px-2 py-1 bg-secondary rounded text-sm"
                        >
                          {translation.lemma}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Examples */}
                {sense.examples.length > 0 && (
                  <div>
                    <h3 className="font-semibold mb-2">Ví dụ:</h3>
                    <div className="space-y-2">
                      {sense.examples.map((example) => (
                        <div key={example.id} className="text-sm">
                          <p className="italic">{example.content}</p>
                          {example.source && (
                            <p className="text-xs text-muted-foreground mt-1">
                              Nguồn: {example.source}
                            </p>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {sense.note && (
                  <p className="text-sm text-muted-foreground">
                    <strong>Ghi chú:</strong> {sense.note}
                  </p>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

