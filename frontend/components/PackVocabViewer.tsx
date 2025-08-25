"use client";

import { useState } from "react";
import ImageZoom from "@/components/ImageZoom";
import { toImageUrl, type Vocab } from "@/lib/api";

type Props = {
  vocabs: Vocab[];
};

export default function PackVocabViewer({ vocabs }: Props) {
  const [mode, setMode] = useState<"list" | "grid">("list");

  if (vocabs.length === 0) {
    return <p className="text-gray-500">No vocabs in this pack.</p>;
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-600">{vocabs.length} items</div>
        <div className="inline-flex rounded overflow-hidden border">
          <button
            type="button"
            onClick={() => setMode("list")}
            className={`px-3 py-1.5 text-sm focus:outline-none focus:ring ${
              mode === "list"
                ? "bg-blue-600 text-white"
                : "bg-white text-gray-800 hover:bg-gray-100"
            }`}
            aria-pressed={mode === "list"}
          >
            List
          </button>
          <button
            type="button"
            onClick={() => setMode("grid")}
            className={`px-3 py-1.5 text-sm border-l focus:outline-none focus:ring ${
              mode === "grid"
                ? "bg-blue-600 text-white"
                : "bg-white text-gray-800 hover:bg-gray-100"
            }`}
            aria-pressed={mode === "grid"}
          >
            Icons
          </button>
        </div>
      </div>

      {mode === "list" ? (
        <div className="overflow-x-auto">
          <table className="min-w-full border-separate border-spacing-0">
            <thead>
              <tr>
                <th className="text-left p-2 border-b">Image</th>
              </tr>
            </thead>
            <tbody>
              {vocabs.map((v) => (
                <tr key={v.id} className="align-top">
                  <td className="p-2 border-b">
                    <ImageZoom
                      src={toImageUrl(v.image)}
                      alt={v.name}
                      thumbSize={96}
                      primaryLabel={v.name}
                      secondaryLabel={v.translation}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
          {vocabs.map((v) => (
            <div key={v.id} className="flex items-center justify-center">
              <ImageZoom
                src={toImageUrl(v.image)}
                alt={v.name}
                thumbSize={112}
                primaryLabel={v.name}
                secondaryLabel={v.translation}
              />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
