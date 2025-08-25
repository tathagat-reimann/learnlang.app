"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { createPack, type Language } from "@/lib/api";

type Props = {
	languages: Language[];
};

export default function CreatePackForm({ languages }: Props) {
	const router = useRouter();
	const [name, setName] = useState("");
	const [lang, setLang] = useState(languages[0]?.id ?? "");
	const [userId, setUserId] = useState("u1");
	const [error, setError] = useState<string | null>(null);
	const [busy, setBusy] = useState(false);
	const noLang = languages.length === 0;

	const onSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError(null);
		if (!name.trim() || !lang) {
			setError("Name and language are required");
			return;
		}
		setBusy(true);
		try {
			await createPack({ name: name.trim(), lang_id: lang, user_id: userId.trim() || "u1" });
			setName("");
			router.refresh();
		} catch (err: unknown) {
			setError(err instanceof Error ? err.message : "Failed to create pack");
		} finally {
			setBusy(false);
		}
	};

	return (
		<form onSubmit={onSubmit} className="mb-6 p-4 border rounded space-y-3">
			<h2 className="text-lg font-medium">Create New Pack</h2>
			{noLang && (
				<p className="text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded px-2 py-1">
					No languages found. Please seed languages in the backend to enable pack creation.
				</p>
			)}
			<div className="grid grid-cols-1 sm:grid-cols-3 gap-3 items-end">
				<label className="flex flex-col gap-1">
					<span className="text-sm text-gray-700">Name</span>
					<input
						className="border rounded px-2 py-1"
						value={name}
						onChange={(e) => setName(e.target.value)}
						placeholder="e.g. Kitchen"
						required
						disabled={noLang}
					/>
				</label>
				<label className="flex flex-col gap-1">
					<span className="text-sm text-gray-700">Language</span>
					<select
						className="border rounded px-2 py-1"
						value={lang}
						onChange={(e) => setLang(e.target.value)}
						required
						disabled={noLang}
					>
						{languages.map((l) => (
							<option key={l.id} value={l.id}>
								{l.name}
							</option>
						))}
					</select>
				</label>
				<label className="flex flex-col gap-1">
					<span className="text-sm text-gray-700">User ID</span>
					<input
						className="border rounded px-2 py-1"
						value={userId}
						onChange={(e) => setUserId(e.target.value)}
						placeholder="u1"
						disabled={noLang}
					/>
				</label>
			</div>
			{error && <p className="text-sm text-red-600">{error}</p>}
			<button
				type="submit"
				disabled={busy || noLang}
				className="inline-flex items-center gap-2 px-3 py-1.5 rounded bg-green-600 text-white disabled:opacity-60"
			>
				{busy ? "Creating…" : "Create"}
			</button>
		</form>
	);
}


