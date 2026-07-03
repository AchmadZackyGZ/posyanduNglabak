<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import { FileText, FileSpreadsheet, Printer } from 'lucide-svelte';

	// DTO LAPORAN
	interface DetailLaporan {
		no: number;
		// Parameter Khusus Balita
		nama_balita?: string;
		usia?: string; // Dipakai Balita & Lansia
		jk?: string;
		berat_kg?: number; // Dipakai Balita & Bumil
		tinggi_cm?: number;
		status_gizi: string;

		// Parameter Khusus Ibu Hamil
		nama_ibu?: string;
		usia_kandungan?: number;
		tekanan_darah?: string; // Dipakai Bumil & Lansia
		catatan?: string; // Dipakai Bumil & Lansia

		// Parameter Khusus Lansia
		nama_lansia?: string;
		gula_darah?: number;
		kolesterol?: number;
	}

	// --- STATE MANAGEMENT ---
	let selectedJenisLaporan = $state('balita'); // 'balita' | 'ibu-hamil' | 'lansia'

	// Default periode: Bulan ini (YYYY-MM)
	let selectedPeriode = $state(new Date().toISOString().slice(0, 7));
	let isLoading = $state(false);

	// State Data
	let summaryStats = $state({
		label1: 'Total',
		val1: 0,
		label2: 'Diperiksa',
		val2: 0,
		label3: 'Naik BB',
		val3: 0,
		label4: 'Turun BB',
		val4: 0
	});
	let laporanList = $state<DetailLaporan[]>([]);

	// --- FUNGSI AMBIL DATA DARI API GOLANG ---
	async function loadLaporan() {
		if (!selectedPeriode) return;
		isLoading = true;
		try {
			const res = await fetchAPI(`/laporan/${selectedJenisLaporan}?periode=${selectedPeriode}`);
			if (res.data) {
				laporanList = res.data.data || [];

				// Pemetaan Ringkasan Dinamis
				if (selectedJenisLaporan === 'balita') {
					summaryStats = {
						label1: 'Total Balita',
						val1: res.data.summary.total_balita,
						label2: 'Balita Ditimbang',
						val2: res.data.summary.balita_ditimbang,
						label3: 'Balita Naik BB',
						val3: res.data.summary.balita_naik_bb,
						label4: 'Balita Turun BB',
						val4: res.data.summary.balita_turun_bb
					};
				} else if (selectedJenisLaporan === 'ibu-hamil') {
					summaryStats = {
						label1: 'Total Ibu Hamil',
						val1: res.data.summary.total_ibu_hamil,
						label2: 'Ibu Diperiksa',
						val2: res.data.summary.ibu_hamil_diperiksa,
						label3: '-',
						val3: 0,
						label4: '-',
						val4: 0
					};
				} else if (selectedJenisLaporan === 'lansia') {
					summaryStats = {
						label1: 'Total Lansia',
						val1: res.data.summary.total_lansia,
						label2: 'Lansia Diperiksa',
						val2: res.data.summary.lansia_diperiksa,
						label3: '-',
						val3: 0,
						label4: '-',
						val4: 0
					};
				}
			}
		} catch (error) {
			console.error('Gagal memuat laporan:', error);
			alert('Gagal mengambil data laporan dari server.');
		} finally {
			isLoading = false;
		}
	}

	// Fetch pertama kali saat komponen dimuat
	onMount(() => {
		loadLaporan();
	});

	// --- FITUR EKSPOR RAKITAN SENDIRI (NATIVE) ---

	// 1. Export ke Excel (CSV format)
	function exportToExcel() {
		if (laporanList.length === 0) {
			alert('Tidak ada data untuk diekspor!');
			return;
		}

		// A. Buat Header yang rapi (Bukan key JSON mentah)
		let headerArray: string[] = [];
		if (selectedJenisLaporan === 'balita') {
			headerArray = ['No', 'Nama Balita', 'Usia', 'Jenis Kelamin', 'Berat (Kg)', 'Tinggi (Cm)', 'Status Gizi'];
		} else if (selectedJenisLaporan === 'ibu-hamil') {
			headerArray = ['No', 'Nama Ibu', 'Usia Kandungan (Mgg)', 'Tensi Darah', 'Berat (Kg)', 'Catatan'];
		} else if (selectedJenisLaporan === 'lansia') {
			headerArray = ['No', 'Nama Lansia', 'Usia', 'Tensi Darah', 'Gula Darah', 'Kolesterol', 'Catatan'];
		}
		// Gunakan Titik Koma (;) agar langsung rapi saat dibuka di Excel regional Indonesia
		const headers = headerArray.join(';'); 

		// B. Mapping Isi Data
		const csvRows = laporanList.map((row) => {
			return Object.values(row)
				.map((value) => `"${value || '-'}"`) // Tambahkan fallback '-' jika data kosong
				.join(';'); // Gunakan Titik Koma (;)
		});

		// C. Tambahkan BOM (\ufeff) agar Excel membaca encoding UTF-8 dengan sempurna
		const csvData = '\ufeff' + [headers, ...csvRows].join('\n');
		const blob = new Blob([csvData], { type: 'text/csv;charset=utf-8;' });
		const url = window.URL.createObjectURL(blob);

		const a = document.createElement('a');
		a.setAttribute('href', url);
		a.setAttribute('download', `Laporan_${selectedJenisLaporan}_${selectedPeriode}.csv`);
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
	}

	// 2. Export ke PDF (via Native Print Browser)
	function exportToPDF() {
		window.print();
	}

	// Helper warna badge status gizi
	function getStatusColor(status: string) {
		const s = status ? status.toUpperCase() : '';
		if (s.includes('NORMAL') || s.includes('BAIK'))
			return 'bg-green-50 text-green-600 border-green-200';
		if (s.includes('KURANG')) return 'bg-amber-50 text-amber-600 border-amber-200';
		if (s.includes('BURUK') || s.includes('STUNTING'))
			return 'bg-red-50 text-red-600 border-red-200';
		return 'bg-blue-50 text-blue-600 border-blue-200';
	}
</script>

<svelte:head>
	<title>Laporan - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="print-container space-y-6">
	<div class="border-b border-gray-100 pb-4">
		<h1 class="text-2xl font-black text-gray-900">Laporan & Arsip</h1>
	</div>

	<!-- FILTER SECTION (Disembunyikan saat cetak PDF) -->
	<div class="no-print rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
		<div class="mb-6 flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
			<h2 class="text-base font-bold text-gray-800">Filter Laporan</h2>
			<div class="flex items-center gap-3">
				<button
					onclick={exportToPDF}
					class="flex cursor-pointer items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-bold text-gray-700 transition hover:border-red-200 hover:bg-red-50 hover:text-red-600 active:scale-95"
				>
					<FileText size={16} /> Cetak PDF
				</button>
				<button
					onclick={exportToExcel}
					class="flex cursor-pointer items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-bold text-gray-700 transition hover:border-green-200 hover:bg-green-50 hover:text-green-600 active:scale-95"
				>
					<FileSpreadsheet size={16} /> Unduh Excel (CSV)
				</button>
			</div>
		</div>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
			<div>
				<label class="mb-2 block text-xs font-bold text-gray-600">Jenis Laporan</label>
				<select
					bind:value={selectedJenisLaporan}
					onchange={loadLaporan}
					class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-700 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
				>
					<option value="balita">Rekap Data Balita</option>
					<option value="ibu-hamil">Rekap Pemeriksaan Ibu Hamil</option>
					<option value="lansia">Rekap Kesehatan Lansia</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-xs font-bold text-gray-600">Periode (Bulan & Tahun)</label>
				<input
					type="month"
					bind:value={selectedPeriode}
					onchange={loadLaporan}
					class="w-full cursor-text rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
				/>
			</div>
		</div>
	</div>

	<!-- HASIL LAPORAN -->
	<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
		<div class="flex items-center justify-between border-b border-gray-100 p-6">
			<h2 class="text-lg font-black tracking-wide text-gray-800 uppercase">
				LAPORAN {selectedJenisLaporan.replace('-', ' ')} – {selectedPeriode}
			</h2>
			<Printer size={20} class="hidden text-gray-300 print:block" />
		</div>

		<div class="p-6">
			<!-- RINGKASAN METRIK DINAMIS -->
			<div class="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
				<div class="rounded-xl border border-green-100 bg-green-50/50 p-5">
					<p class="text-xs font-bold text-green-700">{summaryStats.label1}</p>
					<p class="mt-1 text-3xl font-black text-green-600">
						{isLoading ? '...' : summaryStats.val1}
					</p>
					<p class="text-[11px] font-medium text-green-600/70">Orang</p>
				</div>

				<div class="rounded-xl border border-blue-100 bg-blue-50/50 p-5">
					<p class="text-xs font-bold text-blue-700">{summaryStats.label2}</p>
					<p class="mt-1 text-3xl font-black text-blue-600">
						{isLoading ? '...' : summaryStats.val2}
					</p>
					<p class="text-[11px] font-medium text-blue-600/70">Orang</p>
				</div>

				{#if summaryStats.label3 !== '-'}
					<div class="rounded-xl border border-amber-100 bg-amber-50/40 p-5">
						<p class="text-xs font-bold text-amber-700">{summaryStats.label3}</p>
						<p class="mt-1 text-3xl font-black text-amber-600">
							{isLoading ? '...' : summaryStats.val3}
						</p>
						<p class="text-[11px] font-medium text-amber-600/70">Anak</p>
					</div>
				{/if}

				{#if summaryStats.label4 !== '-'}
					<div class="rounded-xl border border-red-100 bg-red-50/50 p-5">
						<p class="text-xs font-bold text-red-700">{summaryStats.label4}</p>
						<p class="mt-1 text-3xl font-black text-red-600">
							{isLoading ? '...' : summaryStats.val4}
						</p>
						<p class="text-[11px] font-medium text-red-600/70">Anak</p>
					</div>
				{/if}
			</div>

			<!-- TABEL DINAMIS -->
			<div class="w-full overflow-x-auto rounded-xl border border-gray-100">
				<table class="w-full min-w-[900px] text-left text-sm text-gray-600">
					<thead
						class="border-b border-gray-100 bg-gray-50 text-xs font-bold tracking-wider text-gray-500 uppercase"
					>
						<tr>
							<th class="w-16 px-6 py-4 text-center">NO</th>
							{#if selectedJenisLaporan === 'balita'}
								<th class="px-6 py-4">NAMA BALITA</th><th class="px-6 py-4 text-center">USIA</th><th
									class="px-6 py-4 text-center">JK</th
								><th class="px-6 py-4 text-center">BERAT (KG)</th><th class="px-6 py-4 text-center"
									>TINGGI (CM)</th
								><th class="px-6 py-4 text-center">STATUS GIZI</th>
							{:else if selectedJenisLaporan === 'ibu-hamil'}
								<th class="px-6 py-4">NAMA IBU</th><th class="px-6 py-4 text-center"
									>USIA KANDUNGAN</th
								><th class="px-6 py-4 text-center">TENSI</th><th class="px-6 py-4 text-center"
									>BERAT (KG)</th
								><th class="px-6 py-4">CATATAN</th>
							{:else if selectedJenisLaporan === 'lansia'}
								<th class="px-6 py-4">NAMA LANSIA</th><th class="px-6 py-4 text-center">USIA</th><th
									class="px-6 py-4 text-center">TENSI</th
								><th class="px-6 py-4 text-center">GULA DARAH</th><th class="px-6 py-4 text-center"
									>KOLESTEROL</th
								>
							{/if}
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-50 bg-white">
						{#if isLoading}
							<tr
								><td colspan="7" class="px-6 py-8 text-center text-gray-400">Memuat laporan...</td
								></tr
							>
						{:else if laporanList.length === 0}
							<tr
								><td colspan="7" class="px-6 py-8 text-center text-gray-400"
									>Tidak ada data rekam medis pada periode ini.</td
								></tr
							>
						{/if}

						{#each laporanList as row, index (row.no)}
							<tr class="transition-colors hover:bg-gray-50/50">
								<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
								{#if selectedJenisLaporan === 'balita'}
									<td class="px-6 py-4 font-bold text-gray-800">{row.nama_balita}</td>
									<td class="px-6 py-4 text-center text-gray-600">{row.usia}</td>
									<td class="px-6 py-4 text-center text-gray-600">{row.jk}</td>
									<td class="px-6 py-4 text-center font-medium text-gray-700">{row.berat_kg}</td>
									<td class="px-6 py-4 text-center font-medium text-gray-700">{row.tinggi_cm}</td>
									<td class="px-6 py-4 text-center"
										><span
											class={`rounded-md border px-2.5 py-1 text-[10px] font-bold ${getStatusColor(row.status_gizi)}`}
											>{row.status_gizi}</span
										></td
									>
								{:else if selectedJenisLaporan === 'ibu-hamil'}
									<td class="px-6 py-4 font-bold text-gray-800">{row.nama_ibu}</td>
									<td class="px-6 py-4 text-center font-bold text-[#117064]"
										>{row.usia_kandungan} Mgg</td
									>
									<td class="px-6 py-4 text-center font-bold text-red-500">{row.tekanan_darah}</td>
									<td class="px-6 py-4 text-center font-medium text-gray-700">{row.berat_kg}</td>
									<td class="px-6 py-4 text-xs text-gray-500">{row.catatan || '-'}</td>
								{:else if selectedJenisLaporan === 'lansia'}
									<td class="px-6 py-4 font-bold text-gray-800">{row.nama_lansia}</td>
									<td class="px-6 py-4 text-center text-gray-600">{row.usia}</td>
									<td class="px-6 py-4 text-center font-bold text-red-500">{row.tekanan_darah}</td>
									<td class="px-6 py-4 text-center font-bold text-amber-600"
										>{row.gula_darah ? `${row.gula_darah} mg/dL` : '-'}</td
									>
									<td class="px-6 py-4 text-center font-bold text-blue-600"
										>{row.kolesterol ? `${row.kolesterol} mg/dL` : '-'}</td
									>
								{/if}
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>

<!-- CSS KHUSUS UNTUK MODE CETAK (PDF) -->
<style>
	@media print {
		/* Sembunyikan elemen yang tidak perlu dicetak (Filter, Tombol, Sidebar) */
		:global(.no-print) {
			display: none !important;
		}
		
		/* Hilangkan background pembungkus dan bayangan agar tinta hemat dan bersih */
		:global(body) {
			background-color: white !important;
		}
		.print-container {
			box-shadow: none !important;
			border: none !important;
			margin: 0 !important;
			padding: 0 !important;
		}

		/* KUNCI UTAMA: Hilangkan efek scroll (overflow) agar tabel tidak terpotong */
		:global(.overflow-x-auto), :global(.overflow-hidden) {
			overflow: visible !important;
		}

		/* Paksa tabel menggunakan 100% lebar kertas dan tambahkan garis tegas */
		table {
			width: 100% !important;
			border-collapse: collapse !important;
		}
		th, td {
			border: 1px solid #d1d5db !important; /* Tambahkan border abu-abu tegas untuk cetak */
			padding: 8px 12px !important;
			white-space: normal !important; /* Biarkan teks panjang turun ke bawah (wrap) */
		}

		/* Atur kertas menjadi Landscape (Tidur) dan beri margin proporsional */
		@page {
			size: landscape;
			margin: 1.5cm;
		}
	}
</style>
