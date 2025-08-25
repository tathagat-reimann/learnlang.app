import { getPackDetail } from "@/lib/api";
import AddVocabForm from "@/components/AddVocabForm";
import PackVocabViewer from "@/components/PackVocabViewer";
import Link from "next/link";

export const dynamic = "force-dynamic";

type Props = { params: Promise<{ id: string }> };

export default async function PackDetailPage({ params }: Props) {
  const { id } = await params;
  const detail = await getPackDetail(id);

  return (
    <main className="p-6 max-w-5xl mx-auto">
      <div className="mb-4 flex items-center gap-3">
        <Link href="/packs" className="text-blue-600 hover:underline">← Back</Link>
        <h1 className="text-2xl font-semibold">{detail.pack.name}</h1>
      </div>
  <AddVocabForm packId={detail.pack.id} />
  <PackVocabViewer vocabs={detail.vocabs} />
    </main>
  );
}
