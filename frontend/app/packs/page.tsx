import Link from "next/link";
import { getPacks, type Pack } from "@/lib/api";

export const dynamic = "force-dynamic";

export default async function PacksPage() {
  let packs: Pack[] = [];
  try {
    packs = await getPacks();
  } catch {
    return (
      <main className="p-6 max-w-3xl mx-auto">
        <h1 className="text-2xl font-semibold mb-4">Packs</h1>
        <p className="text-red-600">Failed to load packs.</p>
      </main>
    );
  }

  return (
    <main className="p-6 max-w-3xl mx-auto">
      <h1 className="text-2xl font-semibold mb-4">Packs</h1>
      {packs.length === 0 ? (
        <p className="text-gray-500">No packs yet.</p>
      ) : (
        <ul className="space-y-2">
          {packs.map((p) => (
            <li key={p.id}>
              <Link className="text-blue-600 hover:underline" href={`/packs/${p.id}`}>
                {p.name}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </main>
  );
}
