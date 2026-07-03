<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import {
		Package,
		Search,
		Plus,
		Edit,
		Trash2,
		AlertCircle,
		PackagePlus,
		PackageMinus
	} from 'lucide-svelte';

	interface InventarisData {
		id: string;
		nama_barang: string;
		kategori: string;
		stok: number;
	}

	let items = $state<InventarisData[]>([]);
	let isLoading = $state(true);
	let searchQuery = $state('');
	let errorMessage = $state('');

	// State untuk Modal
	let isModalOpen = $state(false);
	let isStokModalOpen = $state(false);
	let modalMode = $state<'tambah' | 'edit'>('tambah');
	let selectedId = $state('');

	let formData = $state({
		nama_barang: '',
		kategori: 'Obat',
		stok: 0
	});

	let stokData = $state({
		stok: 0
	});

	async function loadItems() {
		isLoading = true;
		errorMessage = '';
		try {
			const res = await fetchAPI('/inventaris');
			items = res.data || [];
		} catch (error: unknown) {
			errorMessage = 'Gagal memuat data inventaris.';
			console.error(error);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadItems();
	});

	let filteredItems = $derived(
		items.filter(
			(item) =>
				item.nama_barang.toLowerCase().includes(searchQuery.toLowerCase()) ||
				item.kategori.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	function openAddModal() {
		modalMode = 'tambah';
		formData = { nama_barang: '', kategori: 'Obat', stok: 0 };
		isModalOpen = true;
	}

	function openEditModal(item: InventarisData) {
		modalMode = 'edit';
		selectedId = item.id;
		formData = { nama_barang: item.nama_barang, kategori: item.kategori, stok: item.stok };
		isModalOpen = true;
	}

	function openStokModal(item: InventarisData) {
		selectedId = item.id;
		stokData = { stok: item.stok };
		isStokModalOpen = true;
	}

	function closeModal() {
		isModalOpen = false;
		isStokModalOpen = false;
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		try {
			if (modalMode === 'tambah') {
				await fetchAPI('/inventaris', {
					method: 'POST',
					body: JSON.stringify(formData)
				});
			} else {
				// Saat edit profil barang, stok dikirim tapi di backend diabaikan sesuai logic kita
				await fetchAPI(`/inventaris/${selectedId}`, {
					method: 'PUT',
					body: JSON.stringify(formData)
				});
			}
			closeModal();
			loadItems();
		} catch (error: unknown) {
			console.error(error);
			alert((error as { message: string }).message || 'Terjadi kesalahan saat menyimpan data');
		}
	}

	async function handleStokSubmit(event: Event) {
		event.preventDefault();
		try {
			await fetchAPI(`/inventaris/${selectedId}/stok`, {
				method: 'PATCH',
				body: JSON.stringify({ stok: stokData.stok })
			});
			closeModal();
			loadItems();
		} catch (error: unknown) {
			console.error(error);
			alert((error as { message: string }).message || 'Gagal memperbarui stok');
		}
	}

	async function deleteItem(id: string, name: string) {
		const confirmDelete = confirm(`Hapus barang "${name}" dari inventaris?`);
		if (!confirmDelete) return;

		try {
			await fetchAPI(`/inventaris/${id}`, { method: 'DELETE' });
			loadItems();
		} catch (error: unknown) {
			console.error(error);
			alert((error as { message: string }).message || 'Gagal menghapus barang.');
		}
	}

	// Fungsi helper penyesuaian stok instan di UI Modal
	function adjustStok(amount: number) {
		stokData.stok = Math.max(0, stokData.stok + amount);
	}
</script>

<svelte:head>
	<title>Inventaris Logistik - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<!-- HEADER -->
	<div
		class="flex flex-col items-start justify-between gap-4 rounded-2xl border border-gray-100 bg-white p-6 shadow-sm sm:flex-row sm:items-center"
	>
		<div class="flex items-center gap-4">
			<div
				class="flex h-12 w-12 items-center justify-center rounded-xl bg-[#117064]/10 text-[#117064]"
			>
				<Package size={24} />
			</div>
			<div>
				<h1 class="text-xl font-bold text-gray-800">Manajemen Inventaris Logistik</h1>
				<p class="text-sm text-gray-500">
					Pantau ketersediaan stok obat, vitamin, dan PMT Posyandu.
				</p>
			</div>
		</div>
		<button
			onclick={openAddModal}
			class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-xl bg-[#117064] px-4 py-2.5 text-sm font-bold text-white transition hover:bg-[#0e5c52]"
		>
			<Plus size={18} />
			Tambah Barang
		</button>
	</div>

	<!-- KONTEN TABEL -->
	<div class="rounded-2xl border border-gray-200 bg-white shadow-sm">
		<div class="border-b border-gray-100 p-5">
			<div class="relative max-w-md">
				<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
					<Search size={18} class="text-gray-400" />
				</div>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Cari nama barang atau kategori..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-2.5 pr-4 pl-11 text-sm focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064] focus:outline-none"
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
						<th class="px-6 py-4">Nama Barang</th>
						<th class="px-6 py-4">Kategori</th>
						<th class="px-6 py-4 text-center">Sisa Stok</th>
						<th class="px-6 py-4 text-right">Tindakan</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100">
					{#if isLoading}
						<tr
							><td colspan="4" class="py-10 text-center text-gray-400">Memuat data inventaris...</td
							></tr
						>
					{:else if filteredItems.length === 0}
						<tr
							><td colspan="4" class="py-10 text-center text-gray-400"
								>Belum ada barang di inventaris.</td
							></tr
						>
					{/if}

					{#each filteredItems as item (item.id)}
						<tr class="transition hover:bg-gray-50/50">
							<td class="px-6 py-4 font-bold text-gray-800">{item.nama_barang}</td>
							<td class="px-6 py-4">
								<span
									class="inline-flex rounded-full border border-blue-100 bg-blue-50 px-2.5 py-1 text-xs font-bold text-blue-700"
								>
									{item.kategori}
								</span>
							</td>
							<td class="px-6 py-4 text-center">
								<!-- Badge Stok: Merah jika <= 5, Hijau jika aman -->
								<span
									class="inline-flex min-w-[3rem] items-center justify-center rounded-lg px-2.5 py-1 text-sm font-bold {item.stok <=
									5
										? 'bg-red-100 text-red-700'
										: 'bg-green-100 text-green-700'}"
								>
									{item.stok}
								</span>
							</td>
							<td class="px-6 py-4 text-right">
								<div class="flex items-center justify-end gap-2">
									<button
										onclick={() => openStokModal(item)}
										class="inline-flex cursor-pointer items-center justify-center rounded-lg bg-teal-50 p-2 text-[#117064] transition hover:bg-teal-100"
										title="Ubah Stok Saja"
									>
										<PackagePlus size={16} />
									</button>
									<button
										onclick={() => openEditModal(item)}
										class="inline-flex cursor-pointer items-center justify-center rounded-lg bg-gray-50 p-2 text-gray-500 transition hover:bg-gray-100 hover:text-gray-700"
										title="Edit Detail Barang"
									>
										<Edit size={16} />
									</button>
									<button
										onclick={() => deleteItem(item.id, item.nama_barang)}
										class="inline-flex cursor-pointer items-center justify-center rounded-lg bg-red-50 p-2 text-red-500 transition hover:bg-red-100 hover:text-red-700"
										title="Hapus Barang"
									>
										<Trash2 size={16} />
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- MODAL TAMBAH/EDIT BARANG -->
{#if isModalOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 p-4 backdrop-blur-sm"
	>
		<div class="w-full max-w-md overflow-hidden rounded-3xl bg-white shadow-2xl">
			<form onsubmit={handleSubmit} class="p-6">
				<h2 class="mb-6 text-xl font-bold text-gray-800">
					{modalMode === 'tambah' ? 'Tambah Barang Baru' : 'Edit Detail Barang'}
				</h2>

				<div class="space-y-4">
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-600">Nama Barang</label>
						<input
							type="text"
							required
							bind:value={formData.nama_barang}
							class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064] focus:outline-none"
							placeholder="Contoh: Vitamin A Kapsul Biru"
						/>
					</div>
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-600">Kategori</label>
						<select
							bind:value={formData.kategori}
							class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064] focus:outline-none"
						>
							<option value="Vitamin">Vitamin</option>
							<option value="Obat">Obat-obatan</option>
							<option value="PMT">PMT (Makanan Tambahan)</option>
							<option value="Alat Kesehatan">Alat Kesehatan</option>
						</select>
					</div>
					{#if modalMode === 'tambah'}
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-600">Stok Awal</label>
							<input
								type="number"
								min="0"
								required
								bind:value={formData.stok}
								class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064] focus:outline-none"
							/>
						</div>
					{/if}
				</div>

				<div class="mt-8 flex gap-3">
					<button
						type="button"
						onclick={closeModal}
						class="flex-1 cursor-pointer rounded-xl bg-gray-100 py-2.5 text-sm font-bold text-gray-600 transition hover:bg-gray-200"
					>
						Batal
					</button>
					<button
						type="submit"
						class="flex-1 cursor-pointer rounded-xl bg-[#117064] py-2.5 text-sm font-bold text-white transition hover:bg-[#0e5c52]"
					>
						Simpan Data
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- MODAL KHUSUS UPDATE STOK -->
{#if isStokModalOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 p-4 backdrop-blur-sm"
	>
		<div class="w-full max-w-sm overflow-hidden rounded-3xl bg-white shadow-2xl">
			<form onsubmit={handleStokSubmit} class="p-6 text-center">
				<div
					class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-teal-50 text-[#117064]"
				>
					<Package size={24} />
				</div>
				<h2 class="mb-2 text-xl font-bold text-gray-800">Sesuaikan Stok</h2>
				<p class="mb-6 text-xs text-gray-500">Sesuaikan jumlah akhir stok fisik yang ada.</p>

				<div class="mb-6 flex items-center justify-center gap-4">
					<button
						type="button"
						onclick={() => adjustStok(-1)}
						class="cursor-pointer rounded-full p-2 text-red-500 hover:bg-red-50"
					>
						<PackageMinus size={28} />
					</button>

					<input
						type="number"
						min="0"
						required
						bind:value={stokData.stok}
						class="w-24 rounded-xl border-2 border-[#117064] bg-white px-2 py-3 text-center text-2xl font-bold focus:ring-4 focus:ring-teal-500/20 focus:outline-none"
					/>

					<button
						type="button"
						onclick={() => adjustStok(1)}
						class="cursor-pointer rounded-full p-2 text-green-600 hover:bg-green-50"
					>
						<PackagePlus size={28} />
					</button>
				</div>

				<div class="flex gap-3">
					<button
						type="button"
						onclick={closeModal}
						class="flex-1 cursor-pointer rounded-xl bg-gray-100 py-2.5 text-sm font-bold text-gray-600 transition hover:bg-gray-200"
					>
						Batal
					</button>
					<button
						type="submit"
						class="flex-1 cursor-pointer rounded-xl bg-[#117064] py-2.5 text-sm font-bold text-white transition hover:bg-[#0e5c52]"
					>
						Simpan Stok
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
