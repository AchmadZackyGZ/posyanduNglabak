<script lang="ts">
	import { onMount } from 'svelte';
	import { User, Dumbbell, Baby, CalendarDays, Eye } from 'lucide-svelte';
	import Chart from 'chart.js/auto';

	let chartCanvas: HTMLCanvasElement;
	let chartInstance: Chart | null = null;

	onMount(() => {
		if (chartCanvas) {
			chartInstance = new Chart(chartCanvas, {
				type: 'line',
				data: {
					labels: ['Des 2023', 'Jan 2024', 'Feb 2024', 'Mar 2024', 'Apr 2024', 'Mei 2024'],
					datasets: [
						{
							label: 'Tinggi Badan (cm)',
							data: [68, 71, 73, 74, 73, 74],
							borderColor: '#14a38b',
							backgroundColor: 'transparent',
							tension: 0.4,
							borderWidth: 2,
							pointBackgroundColor: '#14a38b'
						},
						{
							label: 'Berat Badan (kg)',
							data: [8.5, 8.8, 9.0, 9.1, 9.0, 9.2],
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

		return () => {
			if (chartInstance) chartInstance.destroy();
		};
	});
</script>

<svelte:head>
	<title>Beranda - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
		<div
			class="flex items-center gap-4 rounded-xl border-t-4 border-t-teal-500 bg-white p-5 shadow-sm"
		>
			<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-teal-50 text-teal-600">
				<User size={24} />
			</div>
			<div>
				<p class="text-xs font-semibold text-gray-500">Total Balita</p>
				<p class="text-2xl font-black text-gray-800">6</p>
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
				<p class="text-2xl font-black text-gray-800">62</p>
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
				<p class="text-2xl font-black text-gray-800">12</p>
				<p class="text-[10px] text-gray-400">Orang terdaftar</p>
			</div>
		</div>

		<div
			class="flex items-center gap-4 rounded-xl border-t-4 border-t-blue-500 bg-white p-5 shadow-sm"
		>
			<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
				<CalendarDays size={24} />
			</div>
			<div>
				<p class="text-xs font-semibold text-gray-500">Kegiatan Bulan Ini</p>
				<p class="text-2xl font-black text-gray-800">2</p>
				<p class="text-[10px] text-gray-400">Kegiatan</p>
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
					<option>1 Tahun Terakhir</option>
				</select>
			</div>
			<div class="relative h-64 w-full">
				<canvas bind:this={chartCanvas}></canvas>
			</div>
		</div>

		<div class="flex flex-col rounded-xl border border-gray-100 bg-white p-6 shadow-sm">
			<h2 class="mb-4 text-base font-bold text-gray-800">Jadwal Terdekat</h2>

			<div class="flex flex-col gap-4">
				<div class="flex items-center gap-4">
					<div
						class="flex h-14 w-12 flex-col items-center justify-center rounded-lg bg-[#117064] text-white"
					>
						<span class="text-lg font-black">25</span>
						<span class="text-[10px] font-medium tracking-widest uppercase">Mei</span>
					</div>
					<div class="flex-1">
						<h3 class="text-sm font-bold text-gray-800">Posyandu Balita</h3>
						<p class="text-xs text-gray-500">Sabtu, 25 Mei 2024 • 08.00-11.00 WIB</p>
					</div>
					<span class="rounded-md bg-[#117064] px-2 py-1 text-[10px] font-bold text-white"
						>AKAN DATANG</span
					>
				</div>

				<div class="flex items-center gap-4">
					<div
						class="flex h-14 w-12 flex-col items-center justify-center rounded-lg bg-orange-400 text-white"
					>
						<span class="text-lg font-black">08</span>
						<span class="text-[10px] font-medium tracking-widest uppercase">Jun</span>
					</div>
					<div class="flex-1">
						<h3 class="text-sm font-bold text-gray-800">Posyandu Balita</h3>
						<p class="text-xs text-gray-500">Sabtu, 8 Jun 2024 • 08.00-11.00 WIB</p>
					</div>
					<span class="rounded-md bg-[#117064] px-2 py-1 text-[10px] font-bold text-white"
						>AKAN DATANG</span
					>
				</div>
			</div>

			<button
				class="mt-6 mt-auto w-full cursor-pointer rounded-lg border border-[#117064] py-2 text-xs font-bold text-[#117064] transition-colors hover:bg-[#f0fdf4]"
			>
				LIHAT SEMUA JADWAL
			</button>
		</div>
	</div>

	<div class="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm">
		<div class="border-b border-gray-100 p-5">
			<h2 class="text-base font-bold text-gray-800">Balita Terakhir Ditambahkan</h2>
		</div>
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm text-gray-600">
				<thead class="bg-gray-50 text-xs font-bold text-gray-500 uppercase">
					<tr>
						<th class="px-6 py-4">NAMA BALITA</th>
						<th class="px-6 py-4 text-center">USIA</th>
						<th class="px-6 py-4 text-center">JENIS KELAMIN</th>
						<th class="px-6 py-4">ORANG TUA</th>
						<th class="px-6 py-4 text-center">AKSI</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100 bg-white">
					<tr class="transition-colors hover:bg-gray-50/50">
						<td class="px-6 py-4 font-bold text-gray-800">Aisyah Putri</td>
						<td class="px-6 py-4 text-center font-medium">3 th 2 bln</td>
						<td class="px-6 py-4 text-center">Perempuan</td>
						<td class="px-6 py-4 font-medium">Siti Aminah</td>
						<td class="px-6 py-4 text-center">
							<button class="text-[#117064] hover:text-[#0c4e43]"><Eye size={18} /></button>
						</td>
					</tr>
					<tr class="transition-colors hover:bg-gray-50/50">
						<td class="px-6 py-4 font-bold text-gray-800">Muhammad Zaki</td>
						<td class="px-6 py-4 text-center font-medium">4 th 1 bln</td>
						<td class="px-6 py-4 text-center">Laki-laki</td>
						<td class="px-6 py-4 font-medium">Rudi Hartono</td>
						<td class="px-6 py-4 text-center">
							<button class="text-[#117064] hover:text-[#0c4e43]"><Eye size={18} /></button>
						</td>
					</tr>
					<tr class="transition-colors hover:bg-gray-50/50">
						<td class="px-6 py-4 font-bold text-gray-800">Qonita Nurul</td>
						<td class="px-6 py-4 text-center font-medium">2 th 8 bln</td>
						<td class="px-6 py-4 text-center">Perempuan</td>
						<td class="px-6 py-4 font-medium">Dewi Lestari</td>
						<td class="px-6 py-4 text-center">
							<button class="text-[#117064] hover:text-[#0c4e43]"><Eye size={18} /></button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>
