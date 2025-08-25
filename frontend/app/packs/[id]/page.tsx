import { getPackDetail, toImageUrl } from "@/lib/api";
import ImageZoom from "@/components/ImageZoom";
import AddVocabForm from "@/components/AddVocabForm";
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
      {detail.vocabs.length === 0 ? (
        <p className="text-gray-500">No vocabs in this pack.</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full border-separate border-spacing-0">
            <thead>
              <tr>
                <th className="text-left p-2 border-b">Image</th>
              </tr>
            </thead>
            <tbody>
        {detail.vocabs.map((v) => (
                <tr key={v.id} className="align-top">
                  <td className="p-2 border-b">
                    <ImageZoom
                      src={toImageUrl(v.image)}
          alt={v.name}
                      thumbSize={96}
          onRevealLabel={v.translation ? `${v.name} / ${v.translation}` : v.name}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </main>
  );
}
