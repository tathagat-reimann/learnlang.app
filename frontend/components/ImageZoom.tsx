"use client";

import { useEffect, useState } from "react";

type Props = {
  src: string;
  alt?: string;
  thumbSize?: number; // px
  primaryLabel?: string; // shown immediately when zoomed (e.g., name)
  secondaryLabel?: string; // shown after click (e.g., translation)
  onEdit?: () => void; // optional edit handler; when present, show Edit control in overlay
};

export default function ImageZoom({ src, alt = "", thumbSize = 80, primaryLabel, secondaryLabel, onEdit }: Props) {
  const [zoomed, setZoomed] = useState(false);
  const [revealed, setRevealed] = useState(false);

  const toggleZoom = () => {
    if (!zoomed) {
      setZoomed(true);
      setRevealed(false);
    } else if (!revealed && secondaryLabel) {
      // second click reveals translation
      setRevealed(true);
    } else {
      // third click closes
      setZoomed(false);
      setRevealed(false);
    }
  };

  // Close on ESC when zoomed
  useEffect(() => {
    if (!zoomed) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setZoomed(false);
        setRevealed(false);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [zoomed]);

  return (
    <>
      {/* Thumbnail */}
      <button
        type="button"
        onClick={() => setZoomed(true)}
        className="inline-block rounded overflow-hidden border border-gray-200 hover:shadow focus:outline-none focus:ring"
        aria-label={alt || "View image"}
      >
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={src}
          alt={alt}
          width={thumbSize}
          height={thumbSize}
          style={{ objectFit: "cover", width: thumbSize, height: thumbSize }}
        />
      </button>

      {/* Modal overlay when zoomed */}
      {zoomed && (
        <div
          className="fixed inset-0 z-50 bg-black/70 flex items-center justify-center p-4"
          onClick={() => {
            // Close when clicking the backdrop
            setZoomed(false);
            setRevealed(false);
          }}
          role="dialog"
          aria-modal="true"
        >
          <div
            className="relative max-w-[90vw] max-h-[85vh]"
            onClick={(e) => {
              // Interact with the image/content toggling without closing
              e.stopPropagation();
              toggleZoom();
            }}
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={src}
              alt={alt}
              className="max-w-full max-h-[85vh] rounded shadow-lg"
            />
      {onEdit && (
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation();
                  onEdit();
                }}
        className="absolute top-2 right-2 px-2.5 py-1.5 text-sm rounded bg-white text-gray-900 shadow focus:outline-none focus:ring-2 focus:ring-blue-700"
                aria-label="Edit"
              >
                Edit
              </button>
            )}
            {(primaryLabel || (revealed && secondaryLabel)) && (
              <div className="absolute inset-0 flex items-center justify-center">
                <span className="px-4 py-2 bg-black/70 text-white text-xl font-semibold rounded">
                  {revealed && secondaryLabel ? secondaryLabel : primaryLabel}
                </span>
              </div>
            )}
          </div>
        </div>
      )}
    </>
  );
}
