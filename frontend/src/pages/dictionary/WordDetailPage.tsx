/**
 * Word Detail Page Component
 */

import { useParams, useNavigate } from 'react-router-dom';
import { WordDetail } from '@/features/dictionary/components/WordDetail';
import { Button } from '@/components/ui/button';

export default function WordDetailPage() {
  const { wordId } = useParams<{ wordId: string }>();
  const navigate = useNavigate();

  if (!wordId) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-destructive">
          Không tìm thấy ID từ
        </div>
      </div>
    );
  }

  const wordIdNum = parseInt(wordId, 10);
  if (isNaN(wordIdNum)) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-destructive">
          ID từ không hợp lệ
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-6">
        <Button
          variant="ghost"
          onClick={() => navigate('/dictionary')}
          className="mb-4"
        >
          ← Quay lại tìm kiếm
        </Button>
      </div>

      <WordDetail wordId={wordIdNum} />
    </div>
  );
}

