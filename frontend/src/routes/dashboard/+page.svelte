<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	// FIX 1: Hapus ikon yang tidak terpakai
	import {
		User,
		Dumbbell,
		Baby,
		Activity,
		AlertTriangle,
		Info,
		BellRing,
		MapPin
	} from 'lucide-svelte';
	import Chart from 'chart.js/auto';
	import { fetchAPI } from '$lib/api';

	// --- INTERFACES ---
	interface GrafikPertumbuhan {
		bulan: string;
		rata_tinggi: number;
		rata_berat: number;
	}

	interface JadwalSingkat {
		id: string;
		nama_kegiatan: string;
		tanggal: string;
		jam_mulai: string;
		lokasi: string;
	}

	interface Notifikasi {
		tipe: string;
		pesan: string;
	}

	// --- STATE MANAGEMENT ---
	let isLoading = $state(true);

	let kpi = $state({
		total_balita: 0,
		balita_ditimbang_bulan_ini: 0,
		balita_ditimbang_bulan_lalu: 0,
		ibu_hamil_aktif: 0,
		total_lansia: 0
	});

	let jadwalTerdekat = $state<JadwalSingkat[]>([]);
	let notifikasiList = $state<Notifikasi[]>([]);

	// Chart States
	let grafikLabels = $state<string[]>([]);
	let grafikTinggi = $state<number[]>([]);
	let grafikBerat = $state<number[]>([]);
	let chartCanvas = $state<HTMLCanvasElement>();
	let chartInstance: Chart | null = null;

	// --- HELPER FUNCTIONS ---

	function formatTanggalJadwal(tgl: string) {
		const date = new Date(tgl);
		const hari = date.toLocaleDateString('id-ID', { weekday: 'long' });
		const dmy = date.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		});
		const tanggalPendek = date.toLocaleDateString('id-ID', { day: '2-digit' });
		const bulanPendek = date.toLocaleDateString('id-ID', { month: 'short' }).toUpperCase();
		return { hari, dmy, tanggalPendek, bulanPendek };
	}

	function getTrend(ini: number, lalu: number) {
		if (lalu === 0 && ini === 0)
			return { status: 'neutral', text: 'Sama dengan bulan lalu', color: 'text-gray-400' };
		if (lalu === 0 && ini > 0)
			return { status: 'up', text: '↑ Naik 100% dari bulan lalu', color: 'text-green-500' };

		const diff = ((ini - lalu) / lalu) * 100;
		if (diff > 0)
			return {
				status: 'up',
				text: `↑ Naik ${diff.toFixed(0)}% dari bulan lalu`,
				color: 'text-green-500'
			};
		if (diff < 0)
			return {
				status: 'down',
				text: `↓ Turun ${Math.abs(diff).toFixed(0)}% dari bulan lalu`,
				color: 'text-red-500'
			};

		return { status: 'neutral', text: 'Sama dengan bulan lalu', color: 'text-gray-400' };
	}

	// FIX 2: Pindahkan perhitungan tren ke script menggunakan $derived agar reaktif tanpa melanggar aturan HTML Svelte 5
	let trendPartisipasi = $derived(
		getTrend(kpi.balita_ditimbang_bulan_ini, kpi.balita_ditimbang_bulan_lalu)
	);

	// --- FETCH DATA (SINGLE REQUEST) ---
	onMount(() => {
		const fetchDashboardData = async () => {
			try {
				const res = await fetchAPI('/dashboard/summary');
				const data = res.data;

				kpi = data.kpi;
				jadwalTerdekat = data.jadwal_terdekat || [];
				notifikasiList = data.notifikasi || [];

				if (data.grafik_pertumbuhan && data.grafik_pertumbuhan.length > 0) {
					grafikLabels = data.grafik_pertumbuhan.map((g: GrafikPertumbuhan) => g.bulan);
					grafikTinggi = data.grafik_pertumbuhan.map((g: GrafikPertumbuhan) => g.rata_tinggi);
					grafikBerat = data.grafik_pertumbuhan.map((g: GrafikPertumbuhan) => g.rata_berat);
				}
			} catch (error) {
				console.error('Gagal mengambil data beranda:', error);
			} finally {
				isLoading = false;
				setTimeout(renderChart, 50);
			}
		};

		fetchDashboardData();

		return () => {
			if (chartInstance) chartInstance.destroy();
		};
	});

	// --- RENDER CHART ---
	function renderChart() {
		if (chartCanvas) {
			chartInstance = new Chart(chartCanvas, {
				type: 'line',
				data: {
					labels: grafikLabels.length > 0 ? grafikLabels : ['Belum ada data'],
					datasets: [
						{
							label: 'Tinggi Badan (cm)',
							data: grafikTinggi.length > 0 ? grafikTinggi : [0],
							borderColor: '#14a38b',
							backgroundColor: 'transparent',
							tension: 0.4,
							borderWidth: 2,
							pointBackgroundColor: '#14a38b'
						},
						{
							label: 'Berat Badan (kg)',
							data: grafikBerat.length > 0 ? grafikBerat : [0],
							borderColor: '#0f6456',
							backgroundColor: 'transparent',
							tension: 0.4,
							borderWidth: 2,
							pointBackgroundColor: '#0f6456'
						}
					]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false,
					plugins: { legend: { position: 'top', align: 'start' } },
					scales: { y: { beginAtZero: false } }
				}
			});
		}
	}
</script>

<svelte:head>
	<title>Beranda - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	{#if isLoading}
		<div
			class="flex h-64 w-full items-center justify-center rounded-xl border border-gray-100 bg-white shadow-sm"
		>
			<div class="flex flex-col items-center gap-3">
				<div
					class="h-8 w-8 animate-spin rounded-full border-4 border-[#117064] border-t-transparent"
				></div>
				<p class="text-sm font-medium text-gray-500">Menyinkronkan data analitik...</p>
			</div>
		</div>
	{:else}
		<!-- 1. KPI CARDS DENGAN INDIKATOR TREN -->
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-teal-500 bg-white p-5 shadow-sm transition hover:shadow-md"
			>
				<div
					class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-teal-50 text-teal-600"
				>
					<User size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Total Balita</p>
					<p class="text-2xl font-black text-gray-800">{kpi.total_balita}</p>
					<p class="mt-0.5 text-[10px] text-gray-400">Total terdaftar di sistem</p>
				</div>
			</div>

			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-orange-400 bg-white p-5 shadow-sm transition hover:shadow-md"
			>
				<div
					class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-orange-50 text-orange-500"
				>
					<Dumbbell size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Partisipasi Penimbangan</p>
					<p class="text-2xl font-black text-gray-800">{kpi.balita_ditimbang_bulan_ini}</p>
					<!-- FIX 2: Panggil variabel $derived langsung, hapus tag @const -->
					<p class={`mt-0.5 text-[10px] font-bold ${trendPartisipasi.color}`}>
						{trendPartisipasi.text}
					</p>
				</div>
			</div>

			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-pink-400 bg-white p-5 shadow-sm transition hover:shadow-md"
			>
				<div
					class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-pink-50 text-pink-500"
				>
					<Baby size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Ibu Hamil Aktif</p>
					<p class="text-2xl font-black text-gray-800">{kpi.ibu_hamil_aktif}</p>
					<p class="mt-0.5 text-[10px] text-gray-400">Dalam masa pantauan</p>
				</div>
			</div>

			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-blue-500 bg-white p-5 shadow-sm transition hover:shadow-md"
			>
				<div
					class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-600"
				>
					<Activity size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Total Lansia Aktif</p>
					<p class="text-2xl font-black text-gray-800">{kpi.total_lansia}</p>
					<p class="mt-0.5 text-[10px] text-gray-400">Lansia terdaftar</p>
				</div>
			</div>
		</div>

		<!-- 2. GRAFIK & JADWAL -->
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
			<div class="rounded-xl border border-gray-100 bg-white p-6 shadow-sm lg:col-span-2">
				<div class="mb-6 flex items-center justify-between">
					<div>
						<h2 class="text-base font-bold text-gray-800">Grafik Rata-rata Pertumbuhan</h2>
						<p class="text-xs text-gray-500">Berdasarkan hasil pengukuran 6 bulan terakhir</p>
					</div>
				</div>
				<div class="relative h-64 w-full">
					<canvas bind:this={chartCanvas}></canvas>
				</div>
			</div>

			<div class="flex flex-col rounded-xl border border-gray-100 bg-white p-6 shadow-sm">
				<h2 class="mb-4 flex items-center gap-2 text-base font-bold text-gray-800">
					<BellRing size={18} class="text-[#117064]" /> Jadwal Terdekat
				</h2>

				<div class="flex flex-col gap-4">
					{#if jadwalTerdekat.length === 0}
						<div
							class="flex flex-1 items-center justify-center rounded-lg border-2 border-dashed border-gray-200 p-6 text-center text-sm text-gray-500"
						>
							Belum ada agenda kegiatan terdekat.
						</div>
					{:else}
						{#each jadwalTerdekat as jadwal (jadwal.id)}
							{@const formatTgl = formatTanggalJadwal(jadwal.tanggal)}
							<div
								class="flex items-start gap-4 rounded-xl border border-gray-50 bg-gray-50/50 p-3"
							>
								<div
									class="flex h-14 w-12 shrink-0 flex-col items-center justify-center rounded-lg bg-[#117064] text-white shadow-sm"
								>
									<span class="text-lg font-black">{formatTgl.tanggalPendek}</span>
									<span class="text-[10px] font-medium tracking-widest uppercase"
										>{formatTgl.bulanPendek}</span
									>
								</div>
								<div class="flex-1">
									<h3 class="line-clamp-1 text-sm font-bold text-gray-800">
										{jadwal.nama_kegiatan}
									</h3>
									<p class="mt-0.5 text-[11px] font-medium text-gray-500">
										{formatTgl.hari}, {formatTgl.dmy} • Pukul {jadwal.jam_mulai}
									</p>
									<p class="mt-1.5 flex items-center gap-1 text-[11px] font-medium text-gray-500">
										<MapPin size={12} class="text-gray-400" />
										{jadwal.lokasi}
									</p>
								</div>
							</div>
						{/each}
					{/if}
				</div>

				<a
					href="/dashboard/jadwal"
					class="mt-6 block w-full cursor-pointer rounded-lg border border-[#117064] py-2.5 text-center text-xs font-bold text-[#117064] transition-colors hover:bg-[#117064] hover:text-white"
				>
					LIHAT SEMUA AGENDA
				</a>
			</div>
		</div>

		<!-- 3. NOTIFIKASI CERDAS (ACTIONABLE ALERTS) -->
		<div class="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm">
			<div class="flex items-center gap-2 border-b border-gray-100 bg-gray-50/50 p-5">
				<AlertTriangle size={18} class="text-amber-500" />
				<h2 class="text-base font-bold text-gray-800">Perhatian & Tindakan Khusus</h2>
			</div>
			<div class="p-5">
				{#if notifikasiList.length === 0}
					<div class="flex flex-col items-center justify-center py-8 text-center">
						<div class="mb-3 rounded-full bg-green-50 p-3 text-green-500"><Info size={24} /></div>
						<p class="text-sm font-bold text-gray-700">Semua Terkendali!</p>
						<p class="text-xs text-gray-500">
							Tidak ada peringatan gizi atau kehamilan kritis bulan ini.
						</p>
					</div>
				{:else}
					<div class="space-y-3">
						<!-- FIX 3: Gunakan index sebagai key, bukan notif.id karena interface Notifikasi tidak memiliki id -->
						{#each notifikasiList as notif, index (index)}
							<div
								class={`flex items-start gap-3 rounded-xl border-l-4 p-4 ${notif.tipe === 'DANGER' ? 'border-red-500 bg-red-50' : notif.tipe === 'WARNING' ? 'border-amber-500 bg-amber-50' : 'border-blue-500 bg-blue-50'}`}
							>
								<div class="mt-0.5 shrink-0">
									{#if notif.tipe === 'DANGER'}
										<AlertTriangle size={18} class="text-red-600" />
									{:else if notif.tipe === 'WARNING'}
										<Info size={18} class="text-amber-600" />
									{:else}
										<Info size={18} class="text-blue-600" />
									{/if}
								</div>
								<div class="flex-1">
									<h4
										class={`text-sm font-bold ${notif.tipe === 'DANGER' ? 'text-red-800' : notif.tipe === 'WARNING' ? 'text-amber-800' : 'text-blue-800'}`}
									>
										{notif.tipe === 'DANGER' ? 'Peringatan Gizi' : 'Info Pemantauan'}
									</h4>
									<p
										class={`mt-0.5 text-xs font-medium ${notif.tipe === 'DANGER' ? 'text-red-600/90' : notif.tipe === 'WARNING' ? 'text-amber-700/90' : 'text-blue-600/90'}`}
									>
										{notif.pesan}
									</p>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
