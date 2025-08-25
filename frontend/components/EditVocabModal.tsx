"use client";

import { useEffect, useState } from "react";
import type { Vocab } from "@/lib/api";

type Props = {
  vocab: Vocab | null;
  onClose: () => void;
  onSave: (changes: { name: string; translation?: string }) => Promise<void> | void;
};

export default function EditVocabModal({ vocab, onClose, onSave }: Props) {
  const [name, setName] = useState("");
  const [translation, setTranslation] = useState<string>("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setName(vocab?.name ?? "");
    setTranslation(vocab?.translation ?? "");
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
    setBusy(true);
    try {
      await onSave({ name: name.trim(), translation: translation.trim() || undefined });
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to save");
    } finally {
      setBusy(false);
    }
  };

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
