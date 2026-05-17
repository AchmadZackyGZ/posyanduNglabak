<script lang="ts">
	import { Search, Plus, Edit2, Trash2 } from 'lucide-svelte';

	// --- STATE MANAGEMENT ---
	let searchQuery = $state('');
	let filterStatus = $state('Semua Status');

	// DATA DUMMY: Jadwal Kegiatan (Sesuai Referensi Desain)
	let jadwalList = $state([
		{
			id: 1,
			kegiatan: 'Posyandu Balita',
			tanggal: '25 Mei 2024',
			jam: '08.00-11.00 WIB',
			lokasi: 'Balai RW 04',
			status: 'Akan Datang'
		},
		{
			id: 2,
			kegiatan: 'Posyandu Balita',
			tanggal: '8 Jun 2024',
			jam: '08.00-11.00 WIB',
			lokasi: 'Balai RW 04',
			status: 'Akan Datang'
		},
		{
			id: 3,
			kegiatan: 'Penyuluhan Gizi',
			tanggal: '15 Apr 2024',
			jam: '09.00-12.00 WIB',
			lokasi: 'Posyandu Mawar',
			status: 'Selesai'
		},
		{
			id: 4,
			kegiatan: 'Imunisasi Massal',
			tanggal: '20 Mar 2024',
			jam: '08.00-13.00 WIB',
			lokasi: 'Puskesmas Setempat',
			status: 'Selesai'
		}
	]);

	// --- LOGIKA FILTER REAKTIF SVELTE 5 ---
	let filteredJadwal = $derived(
		jadwalList.filter((j) => {
			const matchCari =
				j.kegiatan.toLowerCase().includes(searchQuery.toLowerCase()) ||
				j.lokasi.toLowerCase().includes(searchQuery.toLowerCase());
			const matchStatus = filterStatus === 'Semua Status' || j.status === filterStatus;
			return matchCari && matchStatus;
		})
	);

	// Helper warna badge status
	function getStatusColor(status: string) {
		if (status === 'Akan Datang') return 'bg-blue-50 text-blue-600 border-blue-200';
		if (status === 'Selesai') return 'bg-green-50 text-green-600 border-green-200';
		if (status === 'Dibatalkan') return 'bg-red-50 text-red-600 border-red-200';
		return 'bg-gray-50 text-gray-600 border-gray-200';
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
					placeholder="Cari kegiatan..."
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
					{#if filteredJadwal.length === 0}
						<tr>
							<td colspan="7" class="py-12 text-center text-gray-500">
								Tidak ada jadwal yang cocok dengan filter pencarian.
							</td>
						</tr>
					{/if}

					{#each filteredJadwal as jadwal, index (jadwal.id)}
						<tr class="transition-colors hover:bg-gray-50/50">
							<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
							<td class="px-6 py-4 text-center font-medium text-gray-800">{jadwal.kegiatan}</td>
							<td class="px-6 py-4 text-center">{jadwal.tanggal}</td>
							<td class="px-6 py-4 text-center text-gray-500">{jadwal.jam}</td>
							<td class="px-6 py-4 text-center text-gray-500">{jadwal.lokasi}</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`rounded-full border px-3 py-1 text-xs font-bold ${getStatusColor(jadwal.status)}`}
								>
									{jadwal.status}
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
				</tbody>
			</table>
		</div>

		<div
			class="flex items-center justify-between border-t border-gray-100 p-5 text-sm font-medium text-gray-400"
		>
			<p>Menampilkan 1-{filteredJadwal.length} dari 10 data</p>
			<div class="flex gap-1">
				<button
					class="flex h-8 w-8 items-center justify-center rounded-lg border border-[#117064] bg-[#117064] font-bold text-white shadow-sm"
					>1</button
				>
				<button
					class="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg border border-gray-200 bg-white font-bold text-gray-600 transition hover:bg-gray-50"
					>2</button
				>
				<button
					class="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg border border-gray-200 bg-white font-bold text-gray-600 transition hover:bg-gray-50"
					>›</button
				>
			</div>
		</div>
	</div>
</div>
