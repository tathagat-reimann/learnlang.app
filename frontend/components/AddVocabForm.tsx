"use client";

import { useCallback, useState } from "react";
import { useRouter } from "next/navigation";

type Props = {
  packId: string;
};

export default function AddVocabForm({ packId }: Props) {
  const router = useRouter();
  const [name, setName] = useState("");
  const [translation, setTranslation] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [imageUrl, setImageUrl] = useState("");
  const [dragging, setDragging] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (!translation.trim()) {
      setError("Please enter a translation");
      return;
    }
    let uploadFile: File | null = file;
    if (!uploadFile && imageUrl.trim()) {
      try {
        uploadFile = await fetchUrlAsFile(imageUrl.trim());
      } catch (err: unknown) {
        setError(
          err instanceof Error
            ? err.message
            : "Could not fetch image from URL (CORS?)"
        );
        return;
      }
    }
    if (!uploadFile) {
      setError("Please choose an image or provide a valid image URL");
      return;
    }
    setSubmitting(true);
    try {
      const base = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";
      const fd = new FormData();
      fd.append("name", name.trim());
  fd.append("translation", translation.trim());
      fd.append("pack_id", packId);
      fd.append("image", uploadFile);
      const res = await fetch(`${base.replace(/\/$/, "")}/api/vocabs`, {
        method: "POST",
        body: fd,
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Upload failed: ${res.status}`);
      }
      // Reset form, refresh data
      setName("");
      setTranslation("");
      setFile(null);
      setImageUrl("");
      router.refresh();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setSubmitting(false);
    }
  };

  const onDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
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
  }, []);

  const onDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setDragging(true);
  };
  const onDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setDragging(false);
  };

  async function fetchUrlAsFile(url: string): Promise<File> {
    let res: Response;
    try {
      res = await fetch(url, { mode: "cors" });
    } catch {
      // Last resort try without CORS (will be opaque and unusable), so fail explicitly
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
    <form onSubmit={onSubmit} className="mb-6 p-4 border rounded space-y-3">
      <h2 className="text-lg font-medium">Add Vocab</h2>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 items-end">
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-700">Name</span>
          <input
            className="border rounded px-2 py-1"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. knife"
            required
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-700">Translation</span>
          <input
            className="border rounded px-2 py-1"
            value={translation}
            onChange={(e) => setTranslation(e.target.value)}
            placeholder="e.g. चाकू"
            required
          />
        </label>
        <div className="flex flex-col gap-1">
          <span className="text-sm text-gray-700">Image</span>
          <div
            onDragOver={onDragOver}
            onDragLeave={onDragLeave}
            onDrop={onDrop}
            className={`border-2 border-dashed rounded p-3 text-sm text-gray-600 ${dragging ? "border-blue-500 bg-blue-50" : "border-gray-300"}`}
          >
            <p className="mb-2">Drag & drop an image file here, or choose one:</p>
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
        </div>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto] gap-3 items-end">
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-700">Or provide image URL</span>
          <input
            className="border rounded px-2 py-1"
            value={imageUrl}
            onChange={(e) => setImageUrl(e.target.value)}
            placeholder="https://.../image.jpg"
          />
        </label>
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
          className="inline-flex justify-center px-3 py-1.5 rounded bg-gray-200 hover:bg-gray-300 text-sm"
        >
          Attach URL
        </button>
      </div>
      {error && <p className="text-sm text-red-600">{error}</p>}
      <button
        type="submit"
        disabled={submitting}
        className="inline-flex items-center gap-2 px-3 py-1.5 rounded bg-blue-600 text-white disabled:opacity-60"
      >
        {submitting ? "Uploading…" : "Add"}
      </button>
    </form>
  );
}
