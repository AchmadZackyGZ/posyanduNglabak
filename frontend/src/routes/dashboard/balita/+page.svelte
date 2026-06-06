<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import { Search, Plus, Edit2, Trash2, Baby } from 'lucide-svelte';

	interface Balita {
		id: string;
		nik: string;
		nama_balita: string;
		tanggal_lahir: string;
		jenis_kelamin: string;
		nama_orang_tua: string;
		alamat: string;
		no_hp: string;
	}

	let searchQuery = $state('');
	let balitaList = $state<Balita[]>([]);
	let isLoading = $state(true);

	// --- STATE MODAL & FORM ---
	let isModalOpen = $state(false);
	let isSubmitting = $state(false);

	// State khusus untuk mode Edit
	let isEditMode = $state(false);
	let editId = $state('');

	let formBalita = $state({
		nik: '',
		nama_balita: '',
		tanggal_lahir: '',
		jenis_kelamin: 'L',
		nama_orang_tua: '',
		alamat: '',
		no_hp: ''
	});

	function hitungUsiaBulan(tglLahir: string) {
		const birth = new Date(tglLahir);
		const now = new Date();
		let months = (now.getFullYear() - birth.getFullYear()) * 12;
		months -= birth.getMonth();
		months += now.getMonth();
		return months <= 0 ? 0 : months;
	}

	async function loadDataBalita() {
		try {
			const response = await fetchAPI('/balita');
			if (response.data) {
				balitaList = response.data;
			}
		} catch (error) {
			console.error('Gagal memuat data balita:', error);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadDataBalita();
	});

	let filteredBalita = $derived(
		balitaList.filter((b) => {
			if (!b || !b.nama_balita) return false;
			return (
				b.nama_balita.toLowerCase().includes(searchQuery.toLowerCase()) ||
				(b.nik && b.nik.includes(searchQuery))
			);
		})
	);

	// --- FUNGSI MEMBUKA MODAL TAMBAH BARU ---
	function openAddModal() {
		isEditMode = false;
		editId = '';
		formBalita = {
			nik: '',
			nama_balita: '',
			tanggal_lahir: '',
			jenis_kelamin: 'L',
			nama_orang_tua: '',
			alamat: '',
			no_hp: ''
		};
		isModalOpen = true;
	}

	// --- FUNGSI MEMBUKA MODAL EDIT ---
	function handleEditBalita(balita: Balita) {
		isEditMode = true;
		editId = balita.id;

		// Backend Golang mengirim waktu lengkap (misal: 2026-04-07T00:00:00Z).
		// Kita potong (split) agar hanya mengambil 'YYYY-MM-DD' untuk input type="date"
		const tglLahirSaja = balita.tanggal_lahir ? balita.tanggal_lahir.split('T')[0] : '';

		formBalita = {
			nik: balita.nik,
			nama_balita: balita.nama_balita,
			tanggal_lahir: tglLahirSaja,
			jenis_kelamin: balita.jenis_kelamin,
			nama_orang_tua: balita.nama_orang_tua,
			alamat: balita.alamat,
			no_hp: balita.no_hp
		};
		isModalOpen = true;
	}

	// --- FUNGSI HAPUS DATA BALITA ---
	async function handleHapusBalita(id: string) {
		if (
			!confirm(
				'Apakah Anda yakin ingin menghapus data balita ini? Data yang dihapus tidak dapat dikembalikan.'
			)
		)
			return;

		try {
			// Tembak API DELETE ke Golang
			await fetchAPI(`/balita/${id}`, { method: 'DELETE' });
			await loadDataBalita(); // Segarkan tabel
			alert('Data balita berhasil dihapus!');
		} catch (error) {
			console.error('Gagal menghapus data balita:', error);
			alert('Gagal menghapus data balita.');
		}
	}

	// --- FUNGSI SUBMIT (Bisa POST untuk Tambah, Bisa PUT untuk Edit) ---
	async function handleSubmit(event: Event) {
		event.preventDefault();
		isSubmitting = true;

		try {
			if (isEditMode) {
				// Jalur Update (PUT)
				await fetchAPI(`/balita/${editId}`, {
					method: 'PUT',
					body: JSON.stringify(formBalita)
				});
				alert('Data balita berhasil diperbarui!');
			} else {
				// Jalur Tambah Baru (POST)
				await fetchAPI('/balita/register', {
					method: 'POST',
					body: JSON.stringify(formBalita)
				});
				alert('Data balita berhasil didaftarkan!');
			}

			await loadDataBalita();
			isModalOpen = false;
		} catch (error) {
			if (error instanceof Error) {
				alert(`Gagal ${isEditMode ? 'memperbarui' : 'mendaftarkan'} balita: ` + error.message);
			} else {
				alert('Terjadi kesalahan tidak terduga.');
			}
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Data Balita - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-end">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-[#1e293b]">Manajemen Data Balita</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Kelola pendaftaran, identitas, dan riwayat penimbangan balita.
			</p>
		</div>
		<button
			onclick={openAddModal}
			class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#0f6456] px-5 py-3 text-sm font-bold text-white shadow-md transition-all hover:bg-[#0c4e43] active:scale-95"
		>
			<Plus size={18} strokeWidth={3} />
			Tambah Balita
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
					placeholder="Cari nama atau NIK balita..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-3 pr-4 pl-11 text-sm font-medium text-gray-700 shadow-xs transition-all focus:border-[#0f6456] focus:bg-white focus:ring-2 focus:ring-[#0f6456]/20 focus:outline-none"
				/>
			</div>

			<div
				class="flex flex-shrink-0 items-center gap-3 rounded-xl border border-gray-100 bg-gray-50 px-4 py-2.5 text-sm font-medium text-gray-500"
			>
				<Baby size={20} class="text-[#14a38b]" />
				<span>Total Data:</span>
				<span class="rounded-lg bg-[#e6f6f4] px-2.5 py-0.5 font-black text-[#0f6456]">
					{filteredBalita.length}
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
						<th class="px-6 py-5">NIK / Nama Balita</th>
						<th class="w-24 px-6 py-5 text-center">L/P</th>
						<th class="w-32 px-6 py-5 text-center">Usia</th>
						<th class="px-6 py-5">Nama Ibu</th>
						<th class="w-32 px-6 py-5 text-center">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50 bg-white">
					{#if isLoading}
						<tr>
							<td colspan="6" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<div
										class="mb-3 h-8 w-8 animate-spin rounded-full border-4 border-[#14a38b] border-t-transparent"
									></div>
									<p class="text-base font-medium">Memuat data dari server...</p>
								</div>
							</td>
						</tr>
					{:else if filteredBalita.length === 0}
						<tr>
							<td colspan="6" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<Search size={32} class="mb-3 opacity-50" />
									<p class="text-base font-medium">Tidak ada data balita ditemukan.</p>
								</div>
							</td>
						</tr>
					{/if}

					{#each filteredBalita as balita, index (balita.id)}
						<tr class="transition-colors hover:bg-[#f8fdfb]">
							<td class="px-6 py-4 text-center font-bold text-gray-400">{index + 1}</td>
							<td class="px-6 py-4">
								<p class="text-base font-bold text-[#1e293b]">{balita.nama_balita}</p>
								<p class="mt-0.5 font-mono text-xs text-gray-400">{balita.nik}</p>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-xs font-bold {balita.jenis_kelamin ===
									'L'
										? 'bg-blue-50 text-blue-600'
										: 'bg-pink-50 text-pink-600'}"
								>
									{balita.jenis_kelamin}
								</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="font-bold text-gray-700">{hitungUsiaBulan(balita.tanggal_lahir)}</span>
								<span class="text-xs text-gray-500">Bln</span>
							</td>
							<td class="px-6 py-4">
								<span class="font-medium text-gray-700">{balita.nama_orang_tua}</span>
							</td>
							<td class="px-6 py-4">
								<div class="flex items-center justify-center gap-2">
									<button
										onclick={() => handleEditBalita(balita)}
										class="cursor-pointer rounded-lg bg-blue-50 p-2 text-blue-600 transition-colors hover:bg-blue-100"
										title="Edit Data"
									>
										<Edit2 size={18} />
									</button>
									<button
										onclick={() => handleHapusBalita(balita.id)}
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
					{isEditMode ? 'Edit Data Balita' : 'Registrasi Balita Baru'}
				</h2>
				<p class="text-xs text-gray-500">
					{isEditMode
						? 'Ubah informasi profil balita.'
						: 'Masukkan identitas lengkap balita untuk pendataan Posyandu.'}
				</p>
			</div>

			<form onsubmit={handleSubmit} class="space-y-4 p-6">
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700"
						>NIK Balita (Jika ada) / Nomor KIA</label
					>
					<input
						type="text"
						required
						bind:value={formBalita.nik}
						placeholder="Contoh: 3509123456780001"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
					/>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Nama Lengkap Balita</label>
					<input
						type="text"
						required
						bind:value={formBalita.nama_balita}
						placeholder="Contoh: Budi Santoso"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
					/>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Tanggal Lahir</label>
						<input
							type="date"
							required
							bind:value={formBalita.tanggal_lahir}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						/>
					</div>
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Jenis Kelamin</label>
						<select
							bind:value={formBalita.jenis_kelamin}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
						>
							<option value="L">Laki-Laki (L)</option>
							<option value="P">Perempuan (P)</option>
						</select>
					</div>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700"
						>Nama Orang Tua (Ibu/Ayah)</label
					>
					<input
						type="text"
						required
						bind:value={formBalita.nama_orang_tua}
						placeholder="Contoh: Siti Aminah"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
					/>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Alamat Lengkap</label>
					<textarea
						required
						bind:value={formBalita.alamat}
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
						bind:value={formBalita.no_hp}
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
						{isSubmitting ? 'Menyimpan...' : isEditMode ? 'Simpan Perubahan' : 'Daftarkan Balita'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
