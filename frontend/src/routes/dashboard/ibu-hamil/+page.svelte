<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import { Search, Plus, Edit2, Trash2, HeartPulse, AlertCircle } from 'lucide-svelte';

	// --- 1. INTERFACE TYPESCRIPT ---
	interface IbuHamil {
		id: string;
		nik: string;
		nama_ibu: string;
		no_hp: string;
		hpl: string;
		status_kehamilan: string;
		alamat: string;
		tanggal_lahir: string;
	}

	// --- 2. STATE MANAGEMENT ---
	let searchQuery = $state('');
	let ibuHamilList = $state<IbuHamil[]>([]);
	let isLoading = $state(true);

	// --- STATE MODAL & FORM ---
	let isModalOpen = $state(false);
	let isSubmitting = $state(false);

	let isEditMode = $state(false);
	let editId = $state('');

	let formIbuHamil = $state({
		nik: '',
		nama_ibu: '',
		hpl: '',
		status_kehamilan: 'Normal', // Default
		alamat: '',
		no_hp: '',
		tanggal_lahir: ''
	});

	// --- 3. HELPER: Hitung Usia Kandungan Berdasarkan HPL ---
	function hitungUsiaKandungan(hpl: string) {
		if (!hpl) return 0;
		// Rumus: Manusia normal hamil 40 minggu (280 hari). HPL adalah titik 40 minggu.
		const hplDate = new Date(hpl);
		const today = new Date();
		const diffTime = hplDate.getTime() - today.getTime();
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
		const weeksLeft = Math.floor(diffDays / 7);

		let usia = 40 - weeksLeft;
		if (usia < 0) usia = 0;
		if (usia > 42) usia = 42; // Batas maksimal logis
		return usia;
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'Risiko Tinggi':
				return 'bg-red-50 text-red-600 border-red-200';
			case 'Pemantauan':
				return 'bg-amber-50 text-amber-600 border-amber-200';
			default:
				return 'bg-teal-50 text-teal-600 border-teal-200';
		}
	}

	// --- 4. FUNGSI AMBIL DATA DARI GOLANG ---
	async function loadDataIbuHamil() {
		try {
			const response = await fetchAPI('/ibu-hamil');
			if (response.data) {
				ibuHamilList = response.data;
			}
		} catch (error) {
			console.error('Gagal memuat data ibu hamil:', error);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadDataIbuHamil();
	});

	// --- 5. FILTER REAL-TIME DENGAN SAFEGUARD ---
	let filteredIbuHamil = $derived(
		ibuHamilList.filter((ibu) => {
			if (!ibu || !ibu.nama_ibu) return false;
			return (
				ibu.nama_ibu.toLowerCase().includes(searchQuery.toLowerCase()) ||
				(ibu.nik && ibu.nik.includes(searchQuery))
			);
		})
	);

	// --- 6. FUNGSI KONTROL MODAL ---
	function openAddModal() {
		isEditMode = false;
		editId = '';
		formIbuHamil = {
			nik: '',
			nama_ibu: '',
			hpl: '',
			status_kehamilan: 'Normal',
			alamat: '',
			no_hp: '',
			tanggal_lahir: ''
		};
		isModalOpen = true;
	}

	function handleEditIbuHamil(ibu: IbuHamil) {
		isEditMode = true;
		editId = ibu.id;

		// Potong format timestamp ISO dari Golang (misal 2026-06-10T00:00:00Z -> 2026-06-10)
		const hplSaja = ibu.hpl ? ibu.hpl.split('T')[0] : '';

		formIbuHamil = {
			nik: ibu.nik,
			nama_ibu: ibu.nama_ibu,
			hpl: hplSaja,
			status_kehamilan: ibu.status_kehamilan,
			alamat: ibu.alamat,
			no_hp: ibu.no_hp,
			tanggal_lahir: ibu.tanggal_lahir ? ibu.tanggal_lahir.split('T')[0] : ''
		};
		isModalOpen = true;
	}

	// --- 7. FUNGSI HAPUS (DELETE) ---
	async function handleHapusIbuHamil(id: string) {
		if (!confirm('Hapus data ibu hamil ini secara permanen?')) return;
		try {
			await fetchAPI(`/ibu-hamil/${id}`, { method: 'DELETE' });
			await loadDataIbuHamil();
			alert('Data Ibu Hamil berhasil dihapus!');
		} catch (error) {
			console.error('Gagal menghapus data ibu hamil:', error);
			alert('Gagal menghapus data.');
		}
	}

	// --- 8. FUNGSI SUBMIT (POST / PUT) ---
	async function handleSubmit(event: Event) {
		event.preventDefault();
		isSubmitting = true;

		try {
			if (isEditMode) {
				await fetchAPI(`/ibu-hamil/${editId}`, {
					method: 'PUT',
					body: JSON.stringify(formIbuHamil)
				});
				alert('Data Ibu Hamil berhasil diperbarui!');
			} else {
				await fetchAPI('/ibu-hamil/register', {
					method: 'POST',
					body: JSON.stringify(formIbuHamil)
				});
				alert('Pendaftaran Ibu Hamil berhasil!');
			}

			await loadDataIbuHamil();
			isModalOpen = false;
		} catch (error) {
			console.error('Gagal menyimpan data:', error);
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
	<title>Data Ibu Hamil - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-end">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-[#1e293b]">Manajemen Ibu Hamil</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Kelola pendaftaran kandungan, taksiran persalinan, dan rekam medis klinis.
			</p>
		</div>
		<button
			onclick={openAddModal}
			class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#0f6456] px-5 py-3 text-sm font-bold text-white shadow-md transition-all hover:bg-[#0c4e43] active:scale-95"
		>
			<Plus size={18} strokeWidth={3} />
			Tambah Ibu Hamil
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
					placeholder="Cari nama ibu atau NIK..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-3 pr-4 pl-11 text-sm font-medium text-gray-700 shadow-xs transition-all focus:border-[#0f6456] focus:bg-white focus:ring-2 focus:ring-[#0f6456]/20 focus:outline-none"
				/>
			</div>

			<div
				class="flex flex-shrink-0 items-center gap-3 rounded-xl border border-gray-100 bg-gray-50 px-4 py-2.5 text-sm font-medium text-gray-500"
			>
				<HeartPulse size={20} class="text-[#14a38b]" />
				<span>Total Kandungan:</span>
				<span class="rounded-lg bg-[#e6f6f4] px-2.5 py-0.5 font-black text-[#0f6456]">
					{filteredIbuHamil.length}
				</span>
			</div>
		</div>

		<div class="w-full overflow-x-auto">
			<table class="w-full min-w-[900px] border-collapse text-left text-sm text-gray-600">
				<thead
					class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-500 uppercase"
				>
					<tr>
						<th class="w-16 px-6 py-5 text-center">No</th>
						<th class="px-6 py-5">NIK / Nama Ibu Hamil</th>
						<th class="px-6 py-5 text-center">Kontak (HP)</th>
						<th class="px-6 py-5 text-center">HPL (Perkiraan)</th>
						<th class="px-6 py-5 text-center">Usia Kandungan</th>
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
					{:else if filteredIbuHamil.length === 0}
						<tr>
							<td colspan="7" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<Search size={32} class="mb-3 opacity-50" />
									<p class="text-base font-medium">Tidak ada data Ibu Hamil ditemukan.</p>
								</div>
							</td>
						</tr>
					{/if}

					{#each filteredIbuHamil as ibu, index (ibu.id)}
						<tr class="transition-colors hover:bg-[#f8fdfb]">
							<td class="px-6 py-4 text-center font-bold text-gray-400">{index + 1}</td>
							<td class="px-6 py-4">
								<p class="text-base font-bold text-[#1e293b]">{ibu.nama_ibu}</p>
								<p class="mt-0.5 font-mono text-xs text-gray-400">{ibu.nik}</p>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="font-medium text-gray-700">{ibu.no_hp}</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="font-bold text-gray-700">{ibu.hpl ? ibu.hpl.split('T')[0] : '-'}</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="text-lg font-black text-[#0f6456]">{hitungUsiaKandungan(ibu.hpl)}</span
								>
								<span class="text-xs font-medium text-gray-500"> Mgg</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`inline-flex items-center justify-center rounded-lg border px-3 py-1 text-xs font-bold ${getStatusColor(ibu.status_kehamilan)}`}
								>
									{#if ibu.status_kehamilan === 'Risiko Tinggi'}
										<AlertCircle size={14} class="mr-1.5" />
									{/if}
									{ibu.status_kehamilan}
								</span>
							</td>
							<td class="px-6 py-4">
								<div class="flex items-center justify-center gap-2">
									<button
										onclick={() => handleEditIbuHamil(ibu)}
										class="cursor-pointer rounded-lg bg-blue-50 p-2 text-blue-600 transition-colors hover:bg-blue-100"
										title="Edit Data"
									>
										<Edit2 size={18} />
									</button>
									<button
										onclick={() => handleHapusIbuHamil(ibu.id)}
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

		<div class="relative z-10 w-full max-w-[500px] overflow-hidden rounded-2xl bg-white shadow-2xl">
			<div class="border-b border-gray-100 bg-gray-50/50 px-6 py-4">
				<h2 class="text-lg font-black text-gray-800">
					{isEditMode ? 'Edit Data Ibu Hamil' : 'Registrasi Ibu Hamil Baru'}
				</h2>
				<p class="text-xs text-gray-500">
					{isEditMode
						? 'Ubah informasi rekam medis kehamilan.'
						: 'Masukkan identitas lengkap untuk pendataan ibu hamil.'}
				</p>
			</div>

			<form onsubmit={handleSubmit} class="space-y-4 p-6">
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">NIK Ibu Hamil</label>
					<input
						type="text"
						required
						bind:value={formIbuHamil.nik}
						placeholder="Contoh: 3509123456781001"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
					/>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Nama Lengkap Ibu</label>
					<input
						type="text"
						required
						bind:value={formIbuHamil.nama_ibu}
						placeholder="Contoh: Siti Aminah"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
					/>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">HPL (Perkiraan Lahir)</label
						>
						<input
							type="date"
							required
							bind:value={formIbuHamil.hpl}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						/>
					</div>
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Status Kehamilan</label>
						<select
							bind:value={formIbuHamil.status_kehamilan}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						>
							<option value="Normal">Normal</option>
							<option value="Risiko Tinggi">Risiko Tinggi</option>
							<option value="Pemantauan">Pemantauan Khusus</option>
						</select>
					</div>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Tanggal Lahir Ibu</label>
					<input
						type="date"
						required
						bind:value={formIbuHamil.tanggal_lahir}
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
					/>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Alamat Lengkap</label>
					<textarea
						required
						bind:value={formIbuHamil.alamat}
						placeholder="Contoh: Jl. Mawar No. 10, RT 01 RW 02"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						rows="2"
					></textarea>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700"
						>Nomor HP / WhatsApp (Aktif)</label
					>
					<input
						type="tel"
						required
						disabled={isEditMode}
						bind:value={formIbuHamil.no_hp}
						placeholder="Contoh: 081234567890"
						class="w-full rounded-xl border border-gray-200 bg-gray-50 p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456] disabled:cursor-not-allowed disabled:opacity-60"
					/>
				</div>

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
						{isSubmitting
							? 'Menyimpan...'
							: isEditMode
								? 'Simpan Perubahan'
								: 'Daftarkan Ibu Hamil'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
