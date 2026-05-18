<script lang="ts">
	import { FileText, FileSpreadsheet } from 'lucide-svelte';

	// --- STATE MANAGEMENT ---
	let selectedJenisLaporan = $state('Rekap Data Balita');
	let selectedPeriode = $state('Mei 2024');

	// DATA DUMMY: Ringkasan Metrik
	let summaryStats = $state({
		totalBalita: 86,
		balitaDitimbang: 62,
		balitaNaik: 40,
		balitaTurun: 5
	});

	// DATA DUMMY: Tabel Laporan
	let laporanList = $state([
		{
			id: 1,
			nama: 'Aisyah Putri',
			usia: '1 th 2 bln',
			jk: 'Perempuan',
			berat: 9.2,
			tinggi: 74,
			statusGizi: 'Normal'
		},
		{
			id: 2,
			nama: 'Muhammad Zaki',
			usia: '2 th 1 bln',
			jk: 'Laki-laki',
			berat: 11.5,
			tinggi: 85,
			statusGizi: 'Normal'
		},
		{
			id: 3,
			nama: 'Qonita Nurul',
			usia: '8 bln',
			jk: 'Perempuan',
			berat: 6.8,
			tinggi: 65,
			statusGizi: 'Kurang'
		},
		{
			id: 4,
			nama: 'Fathan Alfarizi',
			usia: '3 th 5 bln',
			jk: 'Laki-laki',
			berat: 14.2,
			tinggi: 96,
			statusGizi: 'Normal'
		}
	]);

	// Helper warna badge status gizi
	function getStatusColor(status: string) {
		if (status === 'Normal') return 'bg-green-50 text-green-600 border-green-200';
		if (status === 'Kurang') return 'bg-amber-50 text-amber-600 border-amber-200';
		if (status === 'Buruk') return 'bg-red-50 text-red-600 border-red-200';
		return 'bg-gray-50 text-gray-600 border-gray-200';
	}
</script>

<svelte:head>
	<title>Laporan - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="border-b border-gray-100 pb-4">
		<h1 class="text-2xl font-black text-gray-900">Laporan</h1>
	</div>

	<div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
		<div class="mb-6 flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
			<h2 class="text-base font-bold text-gray-800">Filter Laporan</h2>
			<div class="flex items-center gap-3">
				<button
					class="flex cursor-pointer items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-bold text-gray-700 transition hover:bg-gray-50 active:scale-95"
				>
					<FileText size={16} class="text-gray-500" />
					PDF
				</button>
				<button
					class="flex cursor-pointer items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-bold text-gray-700 transition hover:bg-gray-50 active:scale-95"
				>
					<FileSpreadsheet size={16} class="text-gray-500" />
					Excel
				</button>
			</div>
		</div>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
			<div>
				<label class="mb-2 block text-xs font-bold text-gray-600">Jenis Laporan</label>
				<select
					bind:value={selectedJenisLaporan}
					class="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-700 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
				>
					<option value="Rekap Data Balita">Rekap Data Balita</option>
					<option value="Rekap Pemeriksaan">Rekap Pemeriksaan</option>
					<option value="Status Gizi">Status Gizi</option>
					<option value="Kegiatan Posyandu">Kegiatan Posyandu</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-xs font-bold text-gray-600">Periode</label>
				<select
					bind:value={selectedPeriode}
					class="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-700 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
				>
					<option value="Mei 2024">Mei 2024</option>
					<option value="April 2024">April 2024</option>
					<option value="Maret 2024">Maret 2024</option>
				</select>
			</div>
		</div>
	</div>

	<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
		<div class="border-b border-gray-100 p-6">
			<h2 class="text-lg font-black text-gray-800">{selectedJenisLaporan} – {selectedPeriode}</h2>
		</div>

		<div class="p-6">
			<div class="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
				<div class="rounded-xl border border-green-100 bg-green-50/50 p-5">
					<p class="text-xs font-bold text-green-700">Total Balita</p>
					<p class="mt-1 text-3xl font-black text-green-600">{summaryStats.totalBalita}</p>
					<p class="text-[11px] font-medium text-green-600/70">Orang</p>
				</div>

				<div class="rounded-xl border border-blue-100 bg-blue-50/50 p-5">
					<p class="text-xs font-bold text-blue-700">Balita Ditimbang</p>
					<p class="mt-1 text-3xl font-black text-blue-600">{summaryStats.balitaDitimbang}</p>
					<p class="text-[11px] font-medium text-blue-600/70">Orang</p>
				</div>

				<div class="rounded-xl border border-amber-100 bg-amber-50/40 p-5">
					<p class="text-xs font-bold text-amber-700">Balita Naik BB</p>
					<p class="mt-1 text-3xl font-black text-amber-600">{summaryStats.balitaNaik}</p>
					<p class="text-[11px] font-medium text-amber-600/70">Orang</p>
				</div>

				<div class="rounded-xl border border-red-100 bg-red-50/50 p-5">
					<p class="text-xs font-bold text-red-700">Balita Turun BB</p>
					<p class="mt-1 text-3xl font-black text-red-600">{summaryStats.balitaTurun}</p>
					<p class="text-[11px] font-medium text-red-600/70">Orang</p>
				</div>
			</div>

			<div class="w-full overflow-x-auto rounded-xl border border-gray-100">
				<table class="w-full min-w-[900px] text-left text-sm text-gray-600">
					<thead
						class="border-b border-gray-100 bg-gray-50 text-xs font-bold tracking-wider text-gray-500 uppercase"
					>
						<tr>
							<th class="w-16 px-6 py-4 text-center">NO</th>
							<th class="px-6 py-4 text-center">NAMA BALITA</th>
							<th class="px-6 py-4 text-center">USIA</th>
							<th class="px-6 py-4 text-center">JK</th>
							<th class="px-6 py-4 text-center">BERAT (KG)</th>
							<th class="px-6 py-4 text-center">TINGGI (CM)</th>
							<th class="px-6 py-4 text-center">STATUS GIZI</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-50 bg-white">
						{#each laporanList as row, index (row.id)}
							<tr class="transition-colors hover:bg-gray-50/50">
								<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
								<td class="px-6 py-4 text-center font-bold text-gray-800">{row.nama}</td>
								<td class="px-6 py-4 text-center text-gray-600">{row.usia}</td>
								<td class="px-6 py-4 text-center text-gray-600">{row.jk}</td>
								<td class="px-6 py-4 text-center font-medium text-gray-700">{row.berat}</td>
								<td class="px-6 py-4 text-center font-medium text-gray-700">{row.tinggi}</td>
								<td class="px-6 py-4 text-center">
									<span
										class={`rounded-md border px-2.5 py-1 text-xs font-bold ${getStatusColor(row.statusGizi)}`}
									>
										{row.statusGizi}
									</span>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>
