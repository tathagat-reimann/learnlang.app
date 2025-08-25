"use client";

import ImageZoom from "@/components/ImageZoom";
import { toImageUrl, type Vocab } from "@/lib/api";

type Props = { vocabs: Vocab[] | null | undefined; onEdit?: (v: Vocab) => void };

export default function PackVocabViewer({ vocabs, onEdit }: Props) {
  const list = Array.isArray(vocabs) ? vocabs : [];
  if (list.length === 0) {
    return <p className="text-gray-500">No vocabs in this pack.</p>;
  }

  return (
    <div className="space-y-3">
      <div className="text-sm text-gray-600">{list.length} items</div>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
        {list.map((v) => (
          <div key={v.id} className="relative flex items-center justify-center group">
            <ImageZoom
              src={toImageUrl(v.image)}
              alt={v.name}
              thumbSize={112}
              primaryLabel={v.name}
              secondaryLabel={v.translation}
            />
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                onEdit?.(v);
              }}
              disabled={!onEdit}
              className="absolute top-1 right-1 opacity-0 group-hover:opacity-100 transition-opacity px-1.5 py-0.5 text-xs rounded bg-white text-gray-900 shadow disabled:opacity-50 focus:opacity-100 focus:outline-none focus:ring-2 focus:ring-blue-700"
              aria-label={`Edit ${v.name}`}
            >
              ✏️
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
