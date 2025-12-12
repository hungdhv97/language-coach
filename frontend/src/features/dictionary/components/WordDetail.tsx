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
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <CardTitle className="text-2xl">{wordDetail.word.lemma}</CardTitle>
              {wordDetail.word.romanization && (
                <p className="text-muted-foreground mt-1">{wordDetail.word.romanization}</p>
              )}
            </div>
            {wordDetail.word.frequency_rank && (
              <div className="text-sm text-muted-foreground">
                Tần suất: #{wordDetail.word.frequency_rank}
              </div>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Word Notes */}
          {wordDetail.word.notes && (
            <div className="p-3 bg-muted rounded-md">
              <p className="text-sm">
                <strong>Ghi chú:</strong> {wordDetail.word.notes}
              </p>
            </div>
          )}

          {/* Pronunciations */}
          {wordDetail.pronunciations.length > 0 && (
            <div className="space-y-2">
              <h3 className="font-semibold">Phát âm:</h3>
              <div className="space-y-2">
                {wordDetail.pronunciations.map((pron) => (
                  <div key={pron.id} className="flex items-center gap-2 text-sm">
                    <div className="flex-1">
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
                    {pron.audio_url && (
                      <audio controls className="h-8">
                        <source src={pron.audio_url} type="audio/mpeg" />
                        <source src={pron.audio_url} type="audio/wav" />
                        Trình duyệt của bạn không hỗ trợ phát audio.
                      </audio>
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
                <div className="space-y-2">
                  <div className="flex items-start justify-between gap-4">
                    <CardTitle className="text-lg flex-1">
                      {sense.sense_order}. {sense.definition}
                    </CardTitle>
                    <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
                      {sense.part_of_speech_name && (
                        <span className="px-2 py-1 bg-primary/10 rounded">
                          {sense.part_of_speech_name}
                        </span>
                      )}
                      {sense.level_name && (
                        <span className="px-2 py-1 bg-secondary rounded">
                          {sense.level_name}
                        </span>
                      )}
                    </div>
                  </div>
                  <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
                    {sense.definition_language_name && (
                      <span>Ngôn ngữ định nghĩa: {sense.definition_language_name}</span>
                    )}
                    {sense.usage_label && (
                      <span className="px-2 py-1 bg-accent rounded">
                        {sense.usage_label}
                      </span>
                    )}
                  </div>
                </div>
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
                          {translation.romanization && (
                            <span className="ml-1 text-xs text-muted-foreground">
                              ({translation.romanization})
                            </span>
                          )}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Examples */}
                {sense.examples.length > 0 && (
                  <div>
                    <h3 className="font-semibold mb-2">Ví dụ:</h3>
                    <div className="space-y-3">
                      {sense.examples.map((example) => (
                        <div key={example.id} className="text-sm space-y-1">
                          <div className="flex items-start gap-2">
                            <p className="italic flex-1">{example.content}</p>
                            {example.audio_url && (
                              <audio controls className="h-8 flex-shrink-0">
                                <source src={example.audio_url} type="audio/mpeg" />
                                <source src={example.audio_url} type="audio/wav" />
                                Trình duyệt của bạn không hỗ trợ phát audio.
                              </audio>
                            )}
                          </div>
                          {example.source && (
                            <p className="text-xs text-muted-foreground">
                              Nguồn: {example.source}
                            </p>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {sense.note && (
                  <div className="p-3 bg-muted rounded-md">
                    <p className="text-sm">
                      <strong>Ghi chú:</strong> {sense.note}
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

