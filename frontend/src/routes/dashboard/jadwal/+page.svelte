<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { Search, Plus, Edit2, Trash2, MapPin } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';

	// --- 1. INTERFACE TYPESCRIPT ---
	interface Jadwal {
		id: string;
		nama_kegiatan: string;
		tanggal: string;
		jam_mulai: string;
		jam_selesai: string;
		lokasi: string;
		status: string;
	}

	// --- 2. STATE MANAGEMENT ---
	let searchQuery = $state('');
	let filterStatus = $state('Semua Status');
	let jadwalList = $state<Jadwal[]>([]);
	let isLoading = $state(true);

	// --- 3. STATE MODAL & FORM TAMBAH JADWAL ---
	let isModalOpen = $state(false);
	let isSubmitting = $state(false);
	let formJadwal = $state({
		nama_kegiatan: '',
		tanggal: '',
		jam_mulai: '',
		jam_selesai: '',
		lokasi: ''
	});

	// State untuk loading lokasi
	let isFetchingLocation = $state(false);

	// --- 4. LOGIKA FILTER REAKTIF SVELTE 5 ---
	let filteredJadwal = $derived(
		jadwalList.filter((j) => {
			const matchCari =
				j.nama_kegiatan.toLowerCase().includes(searchQuery.toLowerCase()) ||
				j.lokasi.toLowerCase().includes(searchQuery.toLowerCase());
			const matchStatus =
				filterStatus === 'Semua Status' || j.status.toUpperCase() === filterStatus.toUpperCase();
			return matchCari && matchStatus;
		})
	);

	// --- 5. HELPER FORMAT TANGGAL & WARNA ---
	function formatTanggal(tglString: string) {
		const date = new Date(tglString);
		return date.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		});
	}

	// Fungsi menarik lokasi saat ini (HTML5 Geolocation + Reverse Geocoding)
	async function dapatkanLokasiSaatIni() {
		// Cek apakah browser mendukung fitur lokasi
		if (!navigator.geolocation) {
			alert('Browser Anda tidak mendukung fitur pelacakan lokasi.');
			return;
		}

		isFetchingLocation = true;

		navigator.geolocation.getCurrentPosition(
			async (position) => {
				const lat = position.coords.latitude;
				const lon = position.coords.longitude;

				try {
					// Menggunakan OpenStreetMap (Nominatim) gratis untuk mengubah koordinat jadi nama tempat
					const response = await fetch(
						`https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lon}`
					);
					const data = await response.json();

					// Ambil data yang paling masuk akal (desa/kelurahan, kecamatan, atau nama tempat)
					const alamatDetail =
						data.address.village || data.address.suburb || data.address.city || data.display_name;

					// Masukkan ke dalam input form
					formJadwal.lokasi = alamatDetail;
				} catch (error) {
					console.error('Gagal mendapatkan detail alamat:', error);
					// Jika API gagal, fallback gunakan koordinat mentah
					formJadwal.lokasi = `${lat}, ${lon}`;
					alert('Gagal membaca nama jalan, menggunakan titik koordinat sebagai gantinya.');
				} finally {
					isFetchingLocation = false;
				}
			},
			(error) => {
				console.error('Error geolokasi:', error);
				alert(
					'Akses lokasi ditolak atau gagal. Pastikan Anda memberikan izin akses lokasi di browser.'
				);
				isFetchingLocation = false;
			},
			{ enableHighAccuracy: true } // Minta akurasi GPS tertinggi
		);
	}

	function getStatusColor(status: string) {
		const safeStatus = status.toUpperCase();
		if (safeStatus === 'AKAN DATANG') return 'bg-blue-50 text-blue-600 border-blue-200';
		if (safeStatus === 'SELESAI') return 'bg-green-50 text-green-600 border-green-200';
		if (safeStatus === 'DIBATALKAN') return 'bg-red-50 text-red-600 border-red-200';
		return 'bg-gray-50 text-gray-600 border-gray-200';
	}

	// --- 6. FUNGSI GET API (DIPANGGIL SAAT HALAMAN DIBUKA) ---
	onMount(async () => {
		try {
			const response = await fetchAPI('/jadwal');
			if (response.data) {
				jadwalList = response.data;
			}
		} catch (error) {
			console.error('Gagal memuat jadwal:', error);
		} finally {
			isLoading = false;
		}
	});

	// --- 7. FUNGSI POST API (SIMPAN JADWAL BARU) ---
	async function handleTambahJadwal(event: Event) {
		event.preventDefault(); // Mencegah halaman reload
		isSubmitting = true;

		try {
			// Tembak API POST ke backend Golang
			const response = await fetchAPI('/jadwal', {
				method: 'POST',
				body: JSON.stringify(formJadwal)
			});

			// Tambahkan data baru ke tabel secara real-time
			jadwalList = [...jadwalList, response.data];

			// Tutup modal dan bersihkan form
			isModalOpen = false;
			formJadwal = { nama_kegiatan: '', tanggal: '', jam_mulai: '', jam_selesai: '', lokasi: '' };

			alert('Jadwal berhasil ditambahkan!');
		} catch (error) {
			if (error instanceof Error) {
				alert('Gagal menambah jadwal: ' + error.message);
			} else {
				alert('Gagal menambah jadwal: Terjadi kesalahan tidak terduga.');
			}
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Jadwal Kegiatan - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="border-b border-gray-100 pb-4">
		<h1 class="text-2xl font-black text-gray-900">Jadwal Kegiatan</h1>
	</div>

	<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
		<div
			class="flex flex-col gap-4 border-b border-gray-100 p-6 md:flex-row md:items-center md:justify-between"
		>
			<div class="relative w-full md:max-w-md">
				<div
					class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-gray-400"
				>
					<Search size={18} />
				</div>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Cari kegiatan atau lokasi..."
					class="block w-full rounded-xl border border-gray-200 bg-white py-2.5 pr-4 pl-11 text-sm font-medium text-gray-700 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
				/>
			</div>

			<div class="flex items-center gap-3">
				<select
					bind:value={filterStatus}
					class="rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-600 transition outline-none focus:border-[#117064]"
				>
					<option value="Semua Status">Semua Status</option>
					<option value="Akan Datang">Akan Datang</option>
					<option value="Selesai">Selesai</option>
					<option value="Dibatalkan">Dibatalkan</option>
				</select>

				<button
					onclick={() => (isModalOpen = true)}
					class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#117064] px-5 py-2.5 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43]"
				>
					<Plus size={18} strokeWidth={3} />
					Tambah Jadwal
				</button>
			</div>
		</div>

		<div class="w-full overflow-x-auto">
			<table class="w-full min-w-[1000px] text-left text-sm text-gray-600">
				<thead
					class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-500 uppercase"
				>
					<tr>
						<th class="w-16 px-6 py-5 text-center">NO</th>
						<th class="px-6 py-5 text-center">KEGIATAN</th>
						<th class="px-6 py-5 text-center">TANGGAL</th>
						<th class="px-6 py-5 text-center">JAM</th>
						<th class="px-6 py-5 text-center">LOKASI</th>
						<th class="px-6 py-5 text-center">STATUS</th>
						<th class="w-24 px-6 py-5 text-center">AKSI</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50 bg-white">
					{#if isLoading}
						<tr>
							<td colspan="7" class="py-12 text-center text-sm font-medium text-gray-500">
								Sedang memuat data jadwal dari server...
							</td>
						</tr>
					{:else if filteredJadwal.length === 0}
						<tr>
							<td colspan="7" class="py-12 text-center text-gray-500">
								Tidak ada jadwal kegiatan.
							</td>
						</tr>
					{:else}
						{#each filteredJadwal as jadwal, index (jadwal.id)}
							<tr class="transition-colors hover:bg-gray-50/50">
								<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
								<td class="px-6 py-4 text-center font-medium text-gray-800"
									>{jadwal.nama_kegiatan}</td
								>
								<td class="px-6 py-4 text-center">{formatTanggal(jadwal.tanggal)}</td>
								<td class="px-6 py-4 text-center text-gray-500"
									>{jadwal.jam_mulai} - {jadwal.jam_selesai}</td
								>
								<td class="px-6 py-4 text-center text-gray-500">{jadwal.lokasi}</td>
								<td class="px-6 py-4 text-center">
									<span
										class={`rounded-full border px-3 py-1 text-xs font-bold ${getStatusColor(jadwal.status)}`}
									>
										{jadwal.status.toUpperCase()}
									</span>
								</td>
								<td class="px-6 py-4">
									<div class="flex justify-center gap-2">
										<button
											class="cursor-pointer rounded-lg p-1.5 text-amber-500 transition hover:bg-amber-50"
											title="Edit"><Edit2 size={16} /></button
										>
										<button
											class="cursor-pointer rounded-lg p-1.5 text-red-500 transition hover:bg-red-50"
											title="Hapus"><Trash2 size={16} /></button
										>
									</div>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>

		<div
			class="flex items-center justify-between border-t border-gray-100 p-5 text-sm font-medium text-gray-400"
		>
			<p>Menampilkan {filteredJadwal.length > 0 ? 1 : 0}-{filteredJadwal.length} data</p>
			<div class="flex gap-1">
				<button
					class="flex h-8 w-8 items-center justify-center rounded-lg border border-[#117064] bg-[#117064] font-bold text-white shadow-sm"
					>1</button
				>
			</div>
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
				<h2 class="text-lg font-black text-gray-800">Tambah Jadwal Baru</h2>
				<p class="text-xs text-gray-500">Masukkan detail agenda Posyandu mendatang.</p>
			</div>

			<form onsubmit={handleTambahJadwal} class="space-y-4 p-6">
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Nama Kegiatan</label>
					<input
						type="text"
						required
						bind:value={formJadwal.nama_kegiatan}
						placeholder="Contoh: Posyandu Balita RW 04"
						class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
					/>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Tanggal</label>
						<input
							type="date"
							required
							bind:value={formJadwal.tanggal}
							class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
						/>
					</div>
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-700">Lokasi</label>
						<div class="flex gap-2">
							<input
								type="text"
								required
								bind:value={formJadwal.lokasi}
								placeholder="Contoh: Balai Desa"
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
							<button
								type="button"
								onclick={dapatkanLokasiSaatIni}
								disabled={isFetchingLocation}
								title="Deteksi Lokasi Saat Ini"
								class="flex shrink-0 cursor-pointer items-center justify-center rounded-xl border border-gray-200 bg-gray-50 px-3.5 text-gray-500 transition hover:bg-teal-50 hover:text-[#117064] disabled:cursor-not-allowed disabled:opacity-50"
							>
								{#if isFetchingLocation}
									<div
										class="h-5 w-5 animate-spin rounded-full border-2 border-[#117064] border-t-transparent"
									></div>
								{:else}
									<MapPin size={20} />
								{/if}
							</button>
						</div>
					</div>

					<div class="grid grid-cols-2 gap-4">
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Jam Mulai</label>
							<input
								type="time"
								required
								bind:value={formJadwal.jam_mulai}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Jam Selesai</label>
							<input
								type="time"
								required
								bind:value={formJadwal.jam_selesai}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
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
							class="cursor-pointer rounded-xl bg-[#117064] px-5 py-2.5 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43] disabled:cursor-not-allowed disabled:opacity-70"
						>
							{isSubmitting ? 'Menyimpan...' : 'Simpan Jadwal'}
						</button>
					</div>
				</div>
			</form>
		</div>
	</div>
{/if}
