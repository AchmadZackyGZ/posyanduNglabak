<script lang="ts">
	import { onMount } from 'svelte';
	// Menggunakan ikon Lucide yang elegan
	import { Baby, HeartPulse, Activity, ActivitySquare } from 'lucide-svelte';
	// Mengimpor Chart.js murni
	import Chart from 'chart.js/auto';

	let userName = $state('');

	// Referensi elemen kanvas untuk Chart.js
	let chartCanvas: HTMLCanvasElement;
	let chartInstance: Chart | null = null;

	// DATA DUMMY METRIK
	let stats = $state({
		totalBalita: 142,
		totalIbuHamil: 38,
		totalLansia: 85
	});

	onMount(() => {
		// 1. Ambil identitas pengguna
		const userData = localStorage.getItem('user');
		if (userData) {
			try {
				userName = JSON.parse(userData).NamaLengkap;
			} catch {
				console.log('Error memuat data profil');
			}
		}

		// 2. Render Grafik Kunjungan Dummy
		if (chartCanvas) {
			chartInstance = new Chart(chartCanvas, {
				type: 'line',
				data: {
					labels: ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun'],
					datasets: [
						{
							label: 'Kunjungan Balita',
							data: [65, 78, 90, 81, 105, 120],
							borderColor: '#0f6456', // Teal Gelap
							backgroundColor: 'rgba(15, 100, 86, 0.1)',
							borderWidth: 3,
							tension: 0.4,
							fill: true
						},
						{
							label: 'Kunjungan Lansia',
							data: [40, 45, 55, 50, 65, 70],
							borderColor: '#14a38b', // Teal Terang
							backgroundColor: 'transparent',
							borderWidth: 3,
							borderDash: [5, 5],
							tension: 0.4
						}
					]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false,
					plugins: {
						legend: { position: 'top' }
					},
					scales: {
						y: { beginAtZero: true }
					}
				}
			});
		}

		// Cleanup: Hancurkan grafik saat pindah halaman agar tidak bocor memori
		return () => {
			if (chartInstance) chartInstance.destroy();
		};
	});
</script>

<svelte:head>
	<title>Beranda - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div
		class="flex flex-col items-start justify-between gap-4 rounded-[24px] bg-gradient-to-r from-[#0f6456] to-[#14a38b] p-8 text-white shadow-xl md:flex-row md:items-center"
	>
		<div>
			<h1 class="text-3xl font-black tracking-tight">Selamat Datang, {userName || 'Petugas'}!</h1>
			<p class="mt-2 text-sm font-medium text-teal-100">
				Pantau ringkasan statistik dan aktivitas Posyandu Sehat Bersama hari ini.
			</p>
		</div>
		<span
			class="rounded-full border border-white/30 bg-white/20 px-4 py-2 text-xs font-bold text-white shadow-sm backdrop-blur-md"
		>
			Mode Operasional
		</span>
	</div>

	<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
		<div
			class="flex items-center justify-between rounded-[20px] border border-gray-100 bg-white p-6 shadow-sm transition-all hover:border-[#14a38b]/30 hover:shadow-md"
		>
			<div>
				<p class="text-xs font-bold tracking-wider text-gray-400 uppercase">Total Balita</p>
				<p class="mt-1 text-4xl font-black text-[#1e293b]">{stats.totalBalita}</p>
				<p class="mt-1 text-xs font-semibold text-[#14a38b]">+12 bulan ini</p>
			</div>
			<div
				class="flex h-16 w-16 items-center justify-center rounded-2xl bg-[#f0fdf4] text-[#0f6456]"
			>
				<Baby size={32} strokeWidth={2.5} />
			</div>
		</div>

		<div
			class="flex items-center justify-between rounded-[20px] border border-gray-100 bg-white p-6 shadow-sm transition-all hover:border-[#14a38b]/30 hover:shadow-md"
		>
			<div>
				<p class="text-xs font-bold tracking-wider text-gray-400 uppercase">Ibu Hamil</p>
				<p class="mt-1 text-4xl font-black text-[#1e293b]">{stats.totalIbuHamil}</p>
				<p class="mt-1 text-xs font-semibold text-[#14a38b]">Status Aktif</p>
			</div>
			<div
				class="flex h-16 w-16 items-center justify-center rounded-2xl bg-[#f0fdf4] text-[#0f6456]"
			>
				<HeartPulse size={32} strokeWidth={2.5} />
			</div>
		</div>

		<div
			class="flex items-center justify-between rounded-[20px] border border-gray-100 bg-white p-6 shadow-sm transition-all hover:border-[#14a38b]/30 hover:shadow-md"
		>
			<div>
				<p class="text-xs font-bold tracking-wider text-gray-400 uppercase">Warga Lansia</p>
				<p class="mt-1 text-4xl font-black text-[#1e293b]">{stats.totalLansia}</p>
				<p class="mt-1 text-xs font-semibold text-[#14a38b]">Terdaftar Rutin</p>
			</div>
			<div
				class="flex h-16 w-16 items-center justify-center rounded-2xl bg-[#f0fdf4] text-[#0f6456]"
			>
				<Activity size={32} strokeWidth={2.5} />
			</div>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
		<div class="rounded-[24px] border border-gray-100 bg-white p-6 shadow-sm lg:col-span-2">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="text-lg font-bold text-[#1e293b]">Statistik Kunjungan Bulanan</h2>
				<button
					class="rounded-lg bg-gray-50 px-3 py-1.5 text-xs font-bold text-gray-500 hover:bg-gray-100"
					>Semester 1</button
				>
			</div>
			<div class="relative h-64 w-full">
				<canvas bind:this={chartCanvas}></canvas>
			</div>
		</div>

		<div class="flex flex-col gap-4 rounded-[24px] border border-gray-100 bg-white p-6 shadow-sm">
			<h2 class="text-lg font-bold text-[#1e293b]">Aktivitas Terkini</h2>

			<div class="flex items-start gap-4 rounded-xl border border-gray-100 bg-gray-50 p-4">
				<div class="mt-0.5 rounded-full bg-[#14a38b] p-1.5 text-white">
					<ActivitySquare size={16} />
				</div>
				<div>
					<p class="text-sm font-bold text-gray-800">Posyandu Balita Selesai</p>
					<p class="mt-0.5 text-xs text-gray-500">
						Kader Mawar mencatat 15 penimbangan baru hari ini.
					</p>
				</div>
			</div>

			<div class="flex items-start gap-4 rounded-xl border border-gray-100 bg-gray-50 p-4">
				<div class="mt-0.5 rounded-full bg-amber-500 p-1.5 text-white">
					<HeartPulse size={16} />
				</div>
				<div>
					<p class="text-sm font-bold text-gray-800">Jadwal Cek Kandungan</p>
					<p class="mt-0.5 text-xs text-gray-500">
						Besok: 5 Ibu Hamil dijadwalkan periksa dengan Bidan.
					</p>
				</div>
			</div>

			<button
				class="mt-auto w-full rounded-xl border-2 border-[#14a38b] bg-white py-2.5 text-sm font-bold text-[#14a38b] transition-colors hover:bg-[#f0fdf4]"
			>
				Lihat Semua Laporan
			</button>
		</div>
	</div>
</div>
