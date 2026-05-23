<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { User, Dumbbell, Baby, Activity, Eye } from 'lucide-svelte';
	import Chart from 'chart.js/auto';
	import { fetchAPI } from '$lib/api';

	interface Balita {
		id: string;
		nama_balita: string;
		tanggal_lahir: string;
		jenis_kelamin: string;
		nama_orang_tua: string;
	}

	interface Jadwal {
		id: string;
		nama_kegiatan: string;
		tanggal: string;
		jam_mulai: string;
		jam_selesai: string;
		status: string;
	}

	interface GrafikPertumbuhan {
    bulan: string;
    rata_tinggi: number;
    rata_berat: number;
}

	let grafikLabels = $state<string[]>([]);
	let grafikTinggi = $state<number[]>([]);
	let grafikBerat = $state<number[]>([]);

	let chartCanvas = $state<HTMLCanvasElement>();
	let chartInstance: Chart | null = null;

	// --- STATE MANAGEMENT DATA BACKEND ---
	let isLoading = $state(true);

	let summary = $state({
		total_balita: 0,
		balita_ditimbang: 0,
		ibu_hamil_aktif: 0,
		total_lansia: 0
	});

	let jadwalTerdekat = $state<Jadwal[]>([]);
	let balitaTerbaru = $state<Balita[]>([]);

	// Helper untuk menghitung usia dari tanggal lahir
	function hitungUsia(tglLahir: string) {
		const birth = new Date(tglLahir);
		const now = new Date();
		let months = (now.getFullYear() - birth.getFullYear()) * 12;
		months -= birth.getMonth();
		months += now.getMonth();

		const years = Math.floor(months / 12);
		const remainingMonths = months % 12;

		if (years === 0) return `${remainingMonths} bln`;
		return `${years} th ${remainingMonths} bln`;
	}

	// Helper untuk format tanggal jadwal
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

	// 3. FIX: onMount harus Sync, kita bungkus fetch di dalam fungsi async internal
	onMount(() => {
		const fetchDashboardData = async () => {
			try {
				const [resSummary, resBalita, resJadwal] = await Promise.all([
					fetchAPI('/dashboard/summary'),
					fetchAPI('/balita'),
					fetchAPI('/jadwal')
				]);

				summary = resSummary.data;

				// extrak data grafik dari backend
				if (resSummary.data.grafik_pertumbuhan && resSummary.data.grafik_pertumbuhan.length > 0) {
					grafikLabels = resSummary.data.grafik_pertumbuhan.map((g: GrafikPertumbuhan) => g.bulan);
					grafikTinggi = resSummary.data.grafik_pertumbuhan.map((g: GrafikPertumbuhan) => g.rata_tinggi);
					grafikBerat = resSummary.data.grafik_pertumbuhan.map((g: GrafikPertumbuhan) => g.rata_berat);
				}

				balitaTerbaru = resBalita.data.slice(0, 5);
				jadwalTerdekat = resJadwal.data.slice(0, 2);
			} catch (error) {
				console.error('Gagal mengambil data beranda:', error);
			} finally {
				isLoading = false;
				// Beri sedikit jeda agar DOM Svelte selesai merender canvas
				setTimeout(renderChart, 50);
			}
		};

		// Panggil eksekusi
		fetchDashboardData();

		// Return fungsi cleanup ChartJS dengan benar secara TypeScript
		return () => {
			if (chartInstance) chartInstance.destroy();
		};
	});

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
					class="h-8 w-8 animate-spin rounded-full border-4 border-teal-500 border-t-transparent"
				></div>
				<p class="text-sm font-medium text-gray-500">Memuat data dashboard...</p>
			</div>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-teal-500 bg-white p-5 shadow-sm"
			>
				<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-teal-50 text-teal-600">
					<User size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Total Balita</p>
					<p class="text-2xl font-black text-gray-800">{summary.total_balita}</p>
					<p class="text-[10px] text-gray-400">Balita terdaftar</p>
				</div>
			</div>

			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-orange-400 bg-white p-5 shadow-sm"
			>
				<div
					class="flex h-12 w-12 items-center justify-center rounded-lg bg-orange-50 text-orange-500"
				>
					<Dumbbell size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Balita Ditimbang</p>
					<p class="text-2xl font-black text-gray-800">{summary.balita_ditimbang}</p>
					<p class="text-[10px] text-gray-400">Bulan ini</p>
				</div>
			</div>

			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-pink-400 bg-white p-5 shadow-sm"
			>
				<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-pink-50 text-pink-500">
					<Baby size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Ibu Hamil</p>
					<p class="text-2xl font-black text-gray-800">{summary.ibu_hamil_aktif}</p>
					<p class="text-[10px] text-gray-400">Aktif terdaftar</p>
				</div>
			</div>

			<div
				class="flex items-center gap-4 rounded-xl border-t-4 border-t-blue-500 bg-white p-5 shadow-sm"
			>
				<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
					<Activity size={24} />
				</div>
				<div>
					<p class="text-xs font-semibold text-gray-500">Total Lansia</p>
					<p class="text-2xl font-black text-gray-800">{summary.total_lansia}</p>
					<p class="text-[10px] text-gray-400">Lansia terdaftar</p>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
			<div class="rounded-xl border border-gray-100 bg-white p-6 shadow-sm lg:col-span-2">
				<div class="mb-6 flex items-center justify-between">
					<h2 class="text-base font-bold text-gray-800">Grafik Pertumbuhan Balita</h2>
					<select
						class="rounded-lg border border-gray-200 bg-gray-50 px-3 py-1.5 text-xs text-gray-600 outline-none"
					>
						<option>6 Bulan Terakhir</option>
					</select>
				</div>
				<div class="relative h-64 w-full">
					<canvas bind:this={chartCanvas}></canvas>
				</div>
			</div>

			<div class="flex flex-col rounded-xl border border-gray-100 bg-white p-6 shadow-sm">
				<h2 class="mb-4 text-base font-bold text-gray-800">Jadwal Terdekat</h2>

				<div class="flex flex-col gap-4">
					{#if jadwalTerdekat.length === 0}
						<div
							class="flex flex-1 items-center justify-center rounded-lg border-2 border-dashed border-gray-200 p-6 text-center text-sm text-gray-500"
						>
							Belum ada jadwal kegiatan posyandu.
						</div>
					{:else}
						{#each jadwalTerdekat as jadwal (jadwal.id)}
							{@const formatTgl = formatTanggalJadwal(jadwal.tanggal)}
							<div class="flex items-center gap-4">
								<div
									class="flex h-14 w-12 flex-col items-center justify-center rounded-lg bg-[#117064] text-white shadow-sm"
								>
									<span class="text-lg font-black">{formatTgl.tanggalPendek}</span>
									<span class="text-[10px] font-medium tracking-widest uppercase"
										>{formatTgl.bulanPendek}</span
									>
								</div>
								<div class="flex-1">
									<h3 class="text-sm font-bold text-gray-800">{jadwal.nama_kegiatan}</h3>
									<p class="text-xs text-gray-500">
										{formatTgl.hari}, {formatTgl.dmy} • {jadwal.jam_mulai}-{jadwal.jam_selesai}
									</p>
								</div>
								<span
									class="rounded-md border border-[#117064] bg-teal-50 px-2 py-1 text-[10px] font-bold text-[#117064]"
								>
									{jadwal.status.toUpperCase()}
								</span>
							</div>
						{/each}
					{/if}
				</div>

				<a
					href="/dashboard/jadwal"
					class="mt-6 block w-full cursor-pointer rounded-lg border border-[#117064] py-2.5 text-center text-xs font-bold text-[#117064] transition-colors hover:bg-[#117064] hover:text-white"
				>
					LIHAT SEMUA JADWAL
				</a>
			</div>
		</div>

		<div class="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm">
			<div class="border-b border-gray-100 p-5">
				<h2 class="text-base font-bold text-gray-800">Balita Terakhir Ditambahkan</h2>
			</div>
			<div class="overflow-x-auto">
				<table class="w-full text-left text-sm text-gray-600">
					<thead class="bg-gray-50 text-xs font-bold tracking-wider text-gray-500 uppercase">
						<tr>
							<th class="px-6 py-4">NAMA BALITA</th>
							<th class="px-6 py-4 text-center">USIA</th>
							<th class="px-6 py-4 text-center">JENIS KELAMIN</th>
							<th class="px-6 py-4">ORANG TUA</th>
							<th class="px-6 py-4 text-center">AKSI</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-50 bg-white">
						{#if balitaTerbaru.length === 0}
							<tr>
								<td colspan="5" class="px-6 py-12 text-center text-sm font-medium text-gray-500">
									Belum ada data balita yang terdaftar.
								</td>
							</tr>
						{:else}
							{#each balitaTerbaru as balita (balita.id)}
								<tr class="transition-colors hover:bg-gray-50/50">
									<td class="px-6 py-4 font-bold text-gray-800">{balita.nama_balita}</td>
									<td class="bg-teal-50/30 px-6 py-4 text-center font-medium text-[#117064]"
										>{hitungUsia(balita.tanggal_lahir)}</td
									>
									<td class="px-6 py-4 text-center">
										{#if balita.jenis_kelamin === 'L'}
											<span
												class="rounded border border-blue-100 bg-blue-50 px-2 py-1 text-xs font-bold text-blue-600"
												>Laki-laki</span
											>
										{:else}
											<span
												class="rounded border border-pink-100 bg-pink-50 px-2 py-1 text-xs font-bold text-pink-600"
												>Perempuan</span
											>
										{/if}
									</td>
									<td class="px-6 py-4 font-medium">{balita.nama_orang_tua}</td>
									<td class="px-6 py-4 text-center">
										<button
											class="cursor-pointer rounded-lg p-2 text-[#117064] transition-colors hover:bg-teal-50"
											title="Lihat Detail"
										>
											<Eye size={18} />
										</button>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>
