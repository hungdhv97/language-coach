/**
 * Dictionary Lookup Page Component
 */

import { DictionarySearch } from '@/features/dictionary/components/DictionarySearch';

export default function DictionaryLookupPage() {

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Từ điển</h1>
        <p className="text-muted-foreground">
          Tìm kiếm và tra cứu thông tin chi tiết về từ vựng
        </p>
      </div>

      <DictionarySearch />
    </div>
  );
}

