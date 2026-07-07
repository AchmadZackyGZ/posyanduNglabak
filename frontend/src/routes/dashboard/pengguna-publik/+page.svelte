<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import { Users, Search, AlertCircle, Trash2 } from 'lucide-svelte';

	interface UserData {
		id: string;
		username: string;
		nama_lengkap: string;
		is_active: boolean;
	}

	let citizens = $state<UserData[]>([]);
	let isLoading = $state(true);
	let searchQuery = $state('');
	let errorMessage = $state('');

	async function loadCitizens() {
		isLoading = true;
		errorMessage = '';
		try {
			// Memanggil endpoint baru yang kita buat di Golang
			const res = await fetchAPI('/pengguna/publik');
			citizens = res.data || [];
		} catch (error: unknown) {
			errorMessage = (error as Error).message || 'Gagal memuat data warga dari server.';
			console.error(error);
		} finally {
			isLoading = false;
		}
	}

	async function deleteCitizen(id: string, name: string) {
		const confirmDelete = confirm(
			`Apakah Anda yakin ingin menghapus akun "${name}" secara permanen?`
		);
		if (!confirmDelete) return;

		try {
			// Memanggil endpoint DELETE yang baru dibuat di Golang
			await fetchAPI(`/pengguna/${id}`, { method: 'DELETE' });
			// Jika berhasil, muat ulang tabel
			loadCitizens();
		} catch (error: unknown) {
			errorMessage = (error as Error).message || 'Gagal menghapus data warga.';
		}
	}

	onMount(() => {
		loadCitizens();
	});

	// Filter pencarian
	let filteredCitizens = $derived(
		citizens.filter(
			(c) =>
				c.nama_lengkap.toLowerCase().includes(searchQuery.toLowerCase()) ||
				c.username.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<svelte:head>
	<title>Daftar Warga - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<!-- HEADER -->
	<div class="flex items-center gap-4 rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
		<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-teal-50 text-teal-600">
			<Users size={24} />
		</div>
		<div>
			<h1 class="text-xl font-bold text-gray-800">Daftar Akun Warga</h1>
			<p class="text-sm text-gray-500">
				Kumpulan akun publik yang terdaftar mandiri di sistem Posyandu.
			</p>
		</div>
	</div>

	<!-- TABEL & PENCARIAN -->
	<div class="rounded-2xl border border-gray-200 bg-white shadow-sm">
		<div class="border-b border-gray-100 p-5">
			<div class="relative max-w-md">
				<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
					<Search size={18} class="text-gray-400" />
				</div>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Cari nama warga atau ID..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-2.5 pr-4 pl-11 text-sm focus:border-teal-500 focus:bg-white focus:ring-1 focus:ring-teal-500 focus:outline-none"
				/>
			</div>
		</div>

		{#if errorMessage}
			<div
				class="m-5 flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 p-4 text-red-600"
			>
				<AlertCircle size={18} class="mt-0.5 shrink-0" />
				<p class="text-sm font-medium">{errorMessage}</p>
			</div>
		{/if}

		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead class="bg-gray-50/80 text-xs font-bold tracking-wider text-gray-500 uppercase">
					<tr>
						<th class="px-6 py-4">Nama Lengkap & ID</th>
						<th class="px-6 py-4 text-center">Status Akun</th>
						<th class="px-6 py-4 text-right">Tindakan</th>
						<!-- TAMBAHAN HEADER -->
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100">
					{#if isLoading}
						<tr
							><td colspan="3" class="py-10 text-center text-gray-400">Memuat data warga...</td></tr
						>
					{:else if filteredCitizens.length === 0}
						<tr
							><td colspan="3" class="py-10 text-center text-gray-400"
								>Tidak ada data warga yang terdaftar.</td
							></tr
						>
					{/if}

					{#each filteredCitizens as c (c.id)}
						<tr class="transition hover:bg-gray-50/50">
							<td class="px-6 py-4">
								<div class="font-bold text-gray-800">{c.nama_lengkap}</div>
								<div class="text-xs font-medium text-gray-500">@{c.username}</div>
							</td>
							<td class="px-6 py-4 text-center">
								{#if c.is_active}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-green-50 px-3 py-1 text-xs font-bold text-green-700"
									>
										<span class="h-1.5 w-1.5 rounded-full bg-green-500"></span> Terverifikasi
									</span>
								{:else}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-gray-100 px-3 py-1 text-xs font-bold text-gray-500"
									>
										<span class="h-1.5 w-1.5 rounded-full bg-gray-400"></span> Nonaktif
									</span>
								{/if}
							</td>
							<!-- TAMBAHAN TOMBOL DELETE -->
							<td class="px-6 py-4 text-right">
								<button
									onclick={() => deleteCitizen(c.id, c.nama_lengkap)}
									class="inline-flex cursor-pointer items-center justify-center rounded-lg bg-red-50 p-2 text-red-500 transition hover:bg-red-100 hover:text-red-700"
									title="Hapus Akun Warga"
								>
									<Trash2 size={16} />
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
