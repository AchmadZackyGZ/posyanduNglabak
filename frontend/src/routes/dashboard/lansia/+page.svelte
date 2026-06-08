<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import { Search, Plus, Edit2, Trash2, Activity, UserCog } from 'lucide-svelte';

	// --- 1. INTERFACE TYPESCRIPT ---
	interface Lansia {
		id: string;
		nik: string;
		nama_lengkap: string;
		tanggal_lahir: string;
		jenis_kelamin: string;
		alamat: string;
		nama_pendamping: string;
		kontak_darurat: string;
		status: string;
	}

	// --- 2. STATE MANAGEMENT ---
	let searchQuery = $state('');
	let lansiaList = $state<Lansia[]>([]);
	let isLoading = $state(true);

	// --- STATE MODAL & FORM ---
	let isModalOpen = $state(false);
	let isSubmitting = $state(false);

	let isEditMode = $state(false);
	let editId = $state('');

	let formLansia = $state({
		nik: '',
		nama_lengkap: '',
		tanggal_lahir: '',
		jenis_kelamin: 'L',
		alamat: '',
		nama_pendamping: '',
		kontak_darurat: '',
		status: 'AKTIF' // Default status
	});

	// --- 3. HELPER: Hitung Usia (Tahun) ---
	function hitungUsia(tglLahir: string) {
		if (!tglLahir) return 0;
		const birth = new Date(tglLahir);
		const now = new Date();
		let age = now.getFullYear() - birth.getFullYear();
		const m = now.getMonth() - birth.getMonth();
		if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) {
			age--;
		}
		return age < 0 ? 0 : age;
	}

	// Warna status keanggotaan (Bukan status klinis)
	function getStatusColor(status: string) {
		const s = status ? status.toUpperCase() : '';
		switch (s) {
			case 'MENINGGAL':
				return 'bg-gray-100 text-gray-600 border-gray-200';
			case 'PINDAH':
				return 'bg-amber-50 text-amber-600 border-amber-200';
			default:
				return 'bg-teal-50 text-teal-600 border-teal-200'; // AKTIF
		}
	}

	// --- 4. FUNGSI AMBIL DATA DARI GOLANG ---
	async function loadDataLansia() {
		try {
			const response = await fetchAPI('/lansia');
			if (response.data) {
				lansiaList = response.data;
			}
		} catch (error) {
			console.error('Gagal memuat data lansia:', error);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadDataLansia();
	});

	// --- 5. FILTER REAL-TIME ---
	let filteredLansia = $derived(
		lansiaList.filter((l) => {
			if (!l || !l.nama_lengkap) return false;
			return (
				l.nama_lengkap.toLowerCase().includes(searchQuery.toLowerCase()) ||
				(l.nik && l.nik.includes(searchQuery))
			);
		})
	);

	// --- 6. FUNGSI KONTROL MODAL ---
	function openAddModal() {
		isEditMode = false;
		editId = '';
		formLansia = {
			nik: '',
			nama_lengkap: '',
			tanggal_lahir: '',
			jenis_kelamin: 'L',
			alamat: '',
			nama_pendamping: '',
			kontak_darurat: '',
			status: 'AKTIF'
		};
		isModalOpen = true;
	}

	function handleEditLansia(lansia: Lansia) {
		isEditMode = true;
		editId = lansia.id;

		// Format tanggal dari backend agar sesuai dengan input type="date"
		const tglLahirSaja = lansia.tanggal_lahir ? lansia.tanggal_lahir.split('T')[0] : '';

		formLansia = {
			nik: lansia.nik,
			nama_lengkap: lansia.nama_lengkap,
			tanggal_lahir: tglLahirSaja,
			jenis_kelamin: lansia.jenis_kelamin,
			alamat: lansia.alamat,
			nama_pendamping: lansia.nama_pendamping || '',
			kontak_darurat: lansia.kontak_darurat || '',
			status: lansia.status || 'AKTIF'
		};
		isModalOpen = true;
	}

	// --- 7. FUNGSI HAPUS (DELETE) ---
	async function handleHapusLansia(id: string) {
		if (
			!confirm(
				'Hapus data lansia ini secara permanen? Data rekam medis yang terkait juga mungkin akan terpengaruh.'
			)
		)
			return;
		try {
			await fetchAPI(`/lansia/${id}`, { method: 'DELETE' });
			await loadDataLansia();
			alert('Data Lansia berhasil dihapus!');
		} catch (error) {
			console.error('Gagal menghapus data lansia:', error);
			alert('Gagal menghapus data.');
		}
	}

	// --- 8. FUNGSI SUBMIT (POST / PUT) ---
	async function handleSubmit(event: Event) {
		event.preventDefault();
		isSubmitting = true;

		try {
			if (isEditMode) {
				await fetchAPI(`/lansia/${editId}`, {
					method: 'PUT',
					body: JSON.stringify(formLansia)
				});
				alert('Data Lansia berhasil diperbarui!');
			} else {
				await fetchAPI('/lansia/register', {
					method: 'POST',
					body: JSON.stringify(formLansia)
				});
				alert('Pendaftaran Lansia berhasil!');
			}

			await loadDataLansia();
			isModalOpen = false;
		} catch (error) {
			console.error('Gagal menyimpan data lansia:', error);
			if (error instanceof Error) {
				alert(`Gagal menyimpan data: ` + error.message);
			} else {
				alert('Terjadi kesalahan tidak terduga.');
			}
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Data Lansia - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-end">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-[#1e293b]">Manajemen Data Lansia</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Kelola pendaftaran master, wali pendamping, dan informasi dasar lansia.
			</p>
		</div>
		<button
			onclick={openAddModal}
			class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#0f6456] px-5 py-3 text-sm font-bold text-white shadow-md transition-all hover:bg-[#0c4e43] active:scale-95"
		>
			<Plus size={18} strokeWidth={3} />
			Tambah Lansia
		</button>
	</div>

	<div class="overflow-hidden rounded-[24px] border border-gray-100 bg-white shadow-sm">
		<div
			class="flex flex-col items-start justify-between gap-4 border-b border-gray-100 bg-white p-6 md:flex-row md:items-center"
		>
			<div class="relative w-full md:w-96">
				<div
					class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-gray-400"
				>
					<Search size={18} />
				</div>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Cari nama lansia atau NIK..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-3 pr-4 pl-11 text-sm font-medium text-gray-700 shadow-xs transition-all focus:border-[#0f6456] focus:bg-white focus:ring-2 focus:ring-[#0f6456]/20 focus:outline-none"
				/>
			</div>

			<div
				class="flex flex-shrink-0 items-center gap-3 rounded-xl border border-gray-100 bg-gray-50 px-4 py-2.5 text-sm font-medium text-gray-500"
			>
				<Activity size={20} class="text-[#14a38b]" />
				<span>Total Terdaftar:</span>
				<span class="rounded-lg bg-[#e6f6f4] px-2.5 py-0.5 font-black text-[#0f6456]">
					{filteredLansia.length}
				</span>
			</div>
		</div>

		<div class="w-full overflow-x-auto">
			<table class="w-full min-w-[1000px] border-collapse text-left text-sm text-gray-600">
				<thead
					class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-500 uppercase"
				>
					<tr>
						<th class="w-16 px-6 py-5 text-center">No</th>
						<th class="px-6 py-5">NIK / Nama Lansia</th>
						<th class="w-24 px-6 py-5 text-center">L/P</th>
						<th class="w-32 px-6 py-5 text-center">Usia</th>
						<th class="px-6 py-5">Pendamping & Kontak</th>
						<th class="px-6 py-5 text-center">Status</th>
						<th class="w-32 px-6 py-5 text-center">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50 bg-white">
					{#if isLoading}
						<tr>
							<td colspan="7" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<div
										class="mb-3 h-8 w-8 animate-spin rounded-full border-4 border-[#14a38b] border-t-transparent"
									></div>
									<p class="text-base font-medium">Memuat data dari server...</p>
								</div>
							</td>
						</tr>
					{:else if filteredLansia.length === 0}
						<tr>
							<td colspan="7" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<Search size={32} class="mb-3 opacity-50" />
									<p class="text-base font-medium">Tidak ada data Lansia ditemukan.</p>
								</div>
							</td>
						</tr>
					{/if}

					{#each filteredLansia as lansia, index (lansia.id)}
						<tr class="transition-colors hover:bg-[#f8fdfb]">
							<td class="px-6 py-4 text-center font-bold text-gray-400">{index + 1}</td>
							<td class="px-6 py-4">
								<p class="text-base font-bold text-[#1e293b]">{lansia.nama_lengkap}</p>
								<p class="mt-0.5 font-mono text-xs text-gray-400">{lansia.nik}</p>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-xs font-bold {lansia.jenis_kelamin ===
									'L'
										? 'bg-blue-50 text-blue-600'
										: 'bg-pink-50 text-pink-600'}"
								>
									{lansia.jenis_kelamin}
								</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="text-lg font-black text-[#0f6456]"
									>{hitungUsia(lansia.tanggal_lahir)}</span
								>
								<span class="text-xs font-medium text-gray-500"> Thn</span>
							</td>
							<td class="px-6 py-4">
								<div class="flex items-center gap-2">
									<UserCog size={16} class="text-gray-400" />
									<div>
										<p class="font-medium text-gray-700">{lansia.nama_pendamping || '-'}</p>
										<p class="text-xs text-gray-500">{lansia.kontak_darurat || '-'}</p>
									</div>
								</div>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`inline-flex items-center justify-center rounded-lg border px-3 py-1 text-xs font-bold ${getStatusColor(lansia.status)}`}
								>
									{lansia.status ? lansia.status.toUpperCase() : 'AKTIF'}
								</span>
							</td>
							<td class="px-6 py-4">
								<div class="flex items-center justify-center gap-2">
									<button
										onclick={() => handleEditLansia(lansia)}
										class="cursor-pointer rounded-lg bg-blue-50 p-2 text-blue-600 transition-colors hover:bg-blue-100"
										title="Edit Data"
									>
										<Edit2 size={18} />
									</button>
									<button
										onclick={() => handleHapusLansia(lansia.id)}
										class="cursor-pointer rounded-lg bg-red-50 p-2 text-red-600 transition-colors hover:bg-red-100"
										title="Hapus Data"
									>
										<Trash2 size={18} />
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

{#if isModalOpen}
	<div
		class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
	>
		<div class="absolute inset-0 cursor-pointer" onclick={() => (isModalOpen = false)}></div>

		<div class="relative z-10 w-full max-w-[600px] overflow-hidden rounded-2xl bg-white shadow-2xl">
			<div class="border-b border-gray-100 bg-gray-50/50 px-6 py-4">
				<h2 class="text-lg font-black text-gray-800">
					{isEditMode ? 'Edit Data Lansia' : 'Registrasi Lansia Baru'}
				</h2>
				<p class="text-xs text-gray-500">
					{isEditMode
						? 'Perbarui informasi profil dan pendamping lansia.'
						: 'Masukkan identitas lengkap untuk pendaftaran program lansia.'}
				</p>
			</div>

			<form onsubmit={handleSubmit} class="max-h-[75vh] space-y-4 overflow-y-auto p-6">
				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div class="md:col-span-2">
						<label class="mb-1.5 block text-xs font-bold text-gray-700">NIK Lansia</label>
						<input
							type="text"
							required
							bind:value={formLansia.nik}
							placeholder="Contoh: 3509123456782001"
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						/>
					</div>

					<div class="md:col-span-2">
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Nama Lengkap</label>
						<input
							type="text"
							required
							bind:value={formLansia.nama_lengkap}
							placeholder="Contoh: Slamet Riyadi"
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						/>
					</div>

					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Tanggal Lahir</label>
						<input
							type="date"
							required
							bind:value={formLansia.tanggal_lahir}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						/>
					</div>
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Jenis Kelamin</label>
						<select
							bind:value={formLansia.jenis_kelamin}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						>
							<option value="L">Laki-Laki (L)</option>
							<option value="P">Perempuan (P)</option>
						</select>
					</div>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Alamat Lengkap</label>
					<textarea
						required
						bind:value={formLansia.alamat}
						placeholder="Contoh: Jl. Diponegoro No. 45"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						rows="2"
					></textarea>
				</div>

				<hr class="my-4 border-gray-100" />

				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700"
							>Nama Wali / Pendamping</label
						>
						<input
							type="text"
							bind:value={formLansia.nama_pendamping}
							placeholder="Nama anak / keluarga"
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						/>
					</div>
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Kontak Darurat (Wali)</label
						>
						<input
							type="tel"
							bind:value={formLansia.kontak_darurat}
							placeholder="Contoh: 08123456789"
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						/>
					</div>
				</div>

				{#if isEditMode}
					<div class="pt-2">
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Status Keanggotaan</label>
						<select
							bind:value={formLansia.status}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						>
							<option value="AKTIF">AKTIF</option>
							<option value="PINDAH">PINDAH LOKASI</option>
							<option value="MENINGGAL">MENINGGAL DUNIA</option>
						</select>
					</div>
				{/if}

				<div class="mt-6 flex items-center justify-end gap-3 border-t border-gray-50 pt-4">
					<button
						type="button"
						onclick={() => (isModalOpen = false)}
						class="cursor-pointer rounded-xl px-5 py-2.5 text-sm font-bold text-gray-500 transition hover:bg-gray-100"
					>
						Batal
					</button>
					<button
						type="submit"
						disabled={isSubmitting}
						class="cursor-pointer rounded-xl bg-[#0f6456] px-5 py-2.5 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43] disabled:cursor-not-allowed disabled:opacity-70"
					>
						{isSubmitting ? 'Menyimpan...' : isEditMode ? 'Simpan Perubahan' : 'Daftarkan Lansia'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
