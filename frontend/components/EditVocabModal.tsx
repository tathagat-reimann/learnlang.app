"use client";

import { useEffect, useState } from "react";
import type { Vocab } from "@/lib/api";

type Props = {
  vocab: Vocab | null;
  onClose: () => void;
  onSave: (changes: { name: string; translation?: string; image?: File }) => Promise<void> | void;
};

export default function EditVocabModal({ vocab, onClose, onSave }: Props) {
  const [name, setName] = useState("");
  const [translation, setTranslation] = useState<string>("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [imageUrl, setImageUrl] = useState("");
  const [dragging, setDragging] = useState(false);

  useEffect(() => {
    setName(vocab?.name ?? "");
    setTranslation(vocab?.translation ?? "");
  setFile(null);
  setImageUrl("");
  }, [vocab]);

  // Close on ESC when modal is open
  useEffect(() => {
    if (!vocab) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [vocab, onClose]);

  if (!vocab) return null;

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    // Build optional image file from either file input or attached URL
    let uploadFile: File | null = file;
    if (!uploadFile && imageUrl.trim()) {
      try {
        uploadFile = await fetchUrlAsFile(imageUrl.trim());
      } catch (err: unknown) {
        setError(err instanceof Error ? err.message : "Could not fetch image from URL");
        return;
      }
    }
    setBusy(true);
    try {
      await onSave({ name: name.trim(), translation: translation.trim() || undefined, image: uploadFile ?? undefined });
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to save");
    } finally {
      setBusy(false);
    }
  };

  function onDragOver(e: React.DragEvent<HTMLDivElement>) {
    e.preventDefault();
    e.stopPropagation();
    setDragging(true);
  }
  function onDragLeave(e: React.DragEvent<HTMLDivElement>) {
    e.preventDefault();
    e.stopPropagation();
    setDragging(false);
  }
  function onDrop(e: React.DragEvent<HTMLDivElement>) {
    e.preventDefault();
    e.stopPropagation();
    setDragging(false);
    const f = e.dataTransfer.files?.[0];
    if (f) {
      if (!f.type.startsWith("image/")) {
        setError("Only image files are allowed");
        return;
      }
      setFile(f);
      setImageUrl("");
    }
  }

  async function fetchUrlAsFile(url: string): Promise<File> {
    let res: Response;
    try {
      res = await fetch(url, { mode: "cors" });
    } catch {
      throw new Error("Failed to fetch the image URL (network/CORS)");
    }
    if (!res.ok) throw new Error(`Image URL responded ${res.status}`);
    const ct = res.headers.get("content-type") || "application/octet-stream";
    if (!ct.startsWith("image/")) throw new Error("URL does not point to an image");
    const blob = await res.blob();
    const u = new URL(url);
    const base = u.pathname.split("/").pop() || "image";
    const extFromCT = ct.includes("png")
      ? ".png"
      : ct.includes("webp")
      ? ".webp"
      : ct.includes("gif")
      ? ".gif"
      : ct.includes("jpeg") || ct.includes("jpg")
      ? ".jpg"
      : "";
    const hasExt = /\.[a-zA-Z0-9]{2,4}$/.test(base);
    const filename = hasExt ? base : `${base}${extFromCT || ".jpg"}`;
    return new File([blob], filename, { type: ct });
  }

  return (
    <div
      className="fixed inset-0 z-50 bg-black/70 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
      onClick={onClose}
    >
      <form
        onSubmit={submit}
        onClick={(e) => e.stopPropagation()}
        className="w-full max-w-md bg-white text-gray-900 rounded shadow-lg p-5 space-y-4"
      >
        <h3 className="text-xl font-semibold">Edit Vocab</h3>
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-900">Name</span>
          <input
            className="border border-gray-400 rounded px-2 py-2 focus:outline-none focus:ring-2 focus:ring-blue-700 focus:border-blue-700"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-900">Translation</span>
          <input
            className="border border-gray-400 rounded px-2 py-2 focus:outline-none focus:ring-2 focus:ring-blue-700 focus:border-blue-700"
            value={translation}
            onChange={(e) => setTranslation(e.target.value)}
          />
        </label>
        <div className="space-y-2">
          <p className="text-sm text-gray-700">Image (leave empty to keep current)</p>
          <div
            onDragOver={onDragOver}
            onDragLeave={onDragLeave}
            onDrop={onDrop}
            className={`border-2 border-dashed rounded p-3 text-sm ${dragging ? "border-blue-700 bg-blue-50" : "border-gray-300"}`}
          >
            <p className="mb-2 text-gray-700">Drag & drop an image file here, or choose one:</p>
            <input
              type="file"
              accept="image/*"
              onChange={(e) => {
                const f = e.target.files?.[0] ?? null;
                setFile(f);
                if (f) setImageUrl("");
              }}
            />
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto] gap-2 items-end">
            <input
              className="border border-gray-400 rounded px-2 py-2 focus:outline-none focus:ring-2 focus:ring-blue-700 focus:border-blue-700"
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              placeholder="https://.../image.jpg"
            />
            <button
              type="button"
              onClick={async () => {
                setError(null);
                try {
                  if (!imageUrl.trim()) return;
                  const f = await fetchUrlAsFile(imageUrl.trim());
                  setFile(f);
                } catch (err: unknown) {
                  setError(err instanceof Error ? err.message : "Could not use URL");
                }
              }}
              className="px-3 py-2 rounded bg-gray-200 hover:bg-gray-300 text-gray-900"
            >
              Attach URL
            </button>
          </div>
        </div>
        {error && <p className="text-sm text-red-700">{error}</p>}
        <div className="flex justify-end gap-2">
          <button
            type="button"
            onClick={onClose}
            className="px-3 py-2 rounded border border-gray-400 text-gray-900 hover:bg-gray-100"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={busy}
            className="px-3 py-2 rounded bg-blue-700 hover:bg-blue-800 text-white disabled:opacity-60"
          >
            {busy ? "Saving…" : "Save"}
          </button>
        </div>
      </form>
    </div>
  );
}
