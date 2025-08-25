"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import type { Vocab } from "@/lib/api";
import PackVocabViewer from "@/components/PackVocabViewer";
import EditVocabModal from "@/components/EditVocabModal";

type Props = { vocabs: Vocab[] | null | undefined };

export default function PackVocabViewerClient({ vocabs }: Props) {
  const router = useRouter();
  const [editing, setEditing] = useState<Vocab | null>(null);

  async function handleSave(_: { name: string; translation?: string }) {
    // TODO: Implement update vocab API; for now just close and refresh.
    setEditing(null);
    router.refresh();
  }

  return (
    <>
      <PackVocabViewer vocabs={vocabs} onEdit={(v) => setEditing(v)} />
      <EditVocabModal vocab={editing} onClose={() => setEditing(null)} onSave={handleSave} />
    </>
  );
}
