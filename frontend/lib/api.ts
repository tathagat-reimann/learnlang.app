export type Pack = {
  id: string;
  name: string;
  lang_id: string;
  user_id: string;
  public?: boolean;
};

export type Vocab = {
  id: string;
  image: string; // path from backend like /files/images/xyz.jpg
  name: string;
  translation?: string;
  pack_id: string;
};

export type PackDetail = {
  pack: Pack;
  vocabs: Vocab[];
};

function getBaseUrl() {
  // Prefer explicit NEXT_PUBLIC_API_BASE for both server and client usage
  const base = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";
  return base.replace(/\/$/, "");
}

type ApiEnvelope<T> = { data: T; meta?: unknown };

export async function getPacks(): Promise<Pack[]> {
  const res = await fetch(`${getBaseUrl()}/api/packs`, {
    next: { revalidate: 0 },
  });
  if (!res.ok) throw new Error(`Failed to fetch packs: ${res.status}`);
  const raw: unknown = await res.json();
  // If server returned bare array
  if (Array.isArray(raw)) return raw as Pack[];
  // Else expect envelope shape { data: T }
  if (
    typeof raw === "object" &&
    raw !== null &&
    "data" in (raw as Record<string, unknown>) &&
    Array.isArray((raw as { data?: unknown }).data)
  ) {
    return (raw as ApiEnvelope<Pack[]>).data;
  }
  return [];
}

export async function getPackDetail(id: string): Promise<PackDetail> {
  const res = await fetch(`${getBaseUrl()}/api/packs/${id}`, {
    next: { revalidate: 0 },
  });
  if (!res.ok) throw new Error(`Failed to fetch pack ${id}: ${res.status}`);
  const raw: unknown = await res.json();
  if (
    typeof raw === "object" &&
    raw !== null &&
    "data" in (raw as Record<string, unknown>)
  ) {
    return (raw as ApiEnvelope<PackDetail>).data;
  }
  return raw as PackDetail;
}

export function toImageUrl(imagePath: string): string {
  // Backend returns "/files/..."; make it absolute for the browser
  const base = getBaseUrl();
  if (imagePath.startsWith("http://") || imagePath.startsWith("https://")) {
    return imagePath;
  }
  if (!imagePath.startsWith("/")) {
    return `${base}/${imagePath}`;
  }
  return `${base}${imagePath}`;
}
