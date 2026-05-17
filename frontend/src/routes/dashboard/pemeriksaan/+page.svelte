<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, Edit2, Trash2 } from 'lucide-svelte';
	import Chart from 'chart.js/auto';

	// Referensi kanvas grafik Donut
	let chartCanvas: HTMLCanvasElement;
	let chartInstance: Chart | null = null;

	// DATA DUMMY Pemeriksaan Balita
	let pemeriksaanList = $state([
		{
			id: 1,
			nama: 'Aisyah Putri',
			tanggal: '20 Mei 2026',
			berat: 9.2,
			tinggi: 74,
			lingkarKepala: 45,
			statusGizi: 'Normal',
			catatan: 'Sehat'
		},
		{
			id: 2,
			nama: 'Muhammad Zaki',
			tanggal: '19 Mei 2026',
			berat: 11.5,
			tinggi: 85,
			lingkarKepala: 47,
			statusGizi: 'Normal',
			catatan: 'Sehat'
		},
		{
			id: 3,
			nama: 'Qonita Nurul',
			tanggal: '18 Mei 2026',
			berat: 6.8,
			tinggi: 65,
			lingkarKepala: 41,
			statusGizi: 'Kurang',
			catatan: 'Perlu perhatian'
		}
	]);

	onMount(() => {
		// Inisialisasi Grafik Donut Status Gizi
		if (chartCanvas) {
			chartInstance = new Chart(chartCanvas, {
				type: 'doughnut',
				data: {
					labels: ['Normal', 'Gizi Kurang', 'Gizi Lebih', 'Gizi Buruk'],
					datasets: [
						{
							data: [62, 15, 8, 2], // Dummy rasio persebaran status gizi
							backgroundColor: [
								'#22c55e', // Hijau (Normal)
								'#f59e0b', // Oranye (Kurang)
								'#3b82f6', // Biru (Lebih)
								'#ef4444' // Merah (Buruk)
							],
							borderWidth: 0,
							hoverOffset: 4
						}
					]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false,
					cutout: '65%',
					plugins: {
						legend: {
							position: 'right',
							labels: { boxWidth: 20, padding: 15, font: { size: 12, weight: 'bold' } }
						}
					}
				}
			});
		}

		return () => {
			if (chartInstance) chartInstance.destroy();
		};
	});

	// Helper warna badge
	function getBadgeColor(status: string) {
		if (status === 'Normal') return 'bg-green-50 text-green-600 border-green-200';
		if (status === 'Kurang') return 'bg-amber-50 text-amber-600 border-amber-200';
		return 'bg-gray-50 text-gray-600 border-gray-200';
	}
</script>

<svelte:head>
	<title>Pemeriksaan - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between border-b border-gray-100 pb-4">
		<h1 class="text-2xl font-black text-gray-900">Pemeriksaan</h1>
	</div>

	<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
		<div class="flex items-center justify-between border-b border-gray-100 p-6">
			<h2 class="text-base font-bold text-gray-800">Pemeriksaan Balita</h2>
			<button
				class="flex cursor-pointer items-center gap-2 rounded-lg bg-[#117064] px-4 py-2 text-sm font-bold text-white transition hover:bg-[#0c4e43]"
			>
				<Plus size={16} strokeWidth={3} />
				Tambah Pemeriksaan
			</button>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm text-gray-600">
				<thead
					class="border-b border-gray-100 bg-white text-xs font-bold tracking-wider text-gray-500 uppercase"
				>
					<tr>
						<th class="px-6 py-4">NO</th>
						<th class="px-6 py-4">NAMA BALITA</th>
						<th class="px-6 py-4">TANGGAL</th>
						<th class="px-6 py-4">BERAT (KG)</th>
						<th class="px-6 py-4">TINGGI (CM)</th>
						<th class="px-6 py-4">LINGKAR KEPALA</th>
						<th class="px-6 py-4 text-center">STATUS GIZI</th>
						<th class="px-6 py-4">CATATAN</th>
						<th class="px-6 py-4 text-center">AKSI</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50">
					{#each pemeriksaanList as item, index (item.id)}
						<tr class="transition-colors hover:bg-gray-50/50">
							<td class="px-6 py-4 text-gray-400">{index + 1}</td>
							<td class="px-6 py-4 font-bold text-gray-800">{item.nama}</td>
							<td class="px-6 py-4">{item.tanggal}</td>
							<td class="px-6 py-4 font-medium">{item.berat}</td>
							<td class="px-6 py-4 font-medium">{item.tinggi}</td>
							<td class="px-6 py-4">{item.lingkarKepala} cm</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`inline-flex rounded-md border px-2.5 py-1 text-xs font-bold ${getBadgeColor(item.statusGizi)}`}
								>
									{item.statusGizi}
								</span>
							</td>
							<td class="px-6 py-4 text-gray-500">{item.catatan}</td>
							<td class="px-6 py-4">
								<div class="flex justify-center gap-2">
									<button class="rounded p-1.5 text-amber-500 transition hover:bg-amber-50"
										><Edit2 size={16} /></button
									>
									<button class="rounded p-1.5 text-red-500 transition hover:bg-red-50"
										><Trash2 size={16} /></button
									>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div
			class="flex items-center justify-between border-t border-gray-100 p-4 text-sm text-gray-500"
		>
			<p>Menampilkan 1-3 dari 62 data</p>
			<div class="flex gap-1">
				<button
					class="flex h-8 w-8 items-center justify-center rounded border border-[#117064] bg-[#117064] font-bold text-white"
					>1</button
				>
				<button
					class="flex h-8 w-8 items-center justify-center rounded border border-gray-200 bg-white font-bold text-gray-600 hover:bg-gray-50"
					>2</button
				>
				<button
					class="flex h-8 w-8 items-center justify-center rounded border border-gray-200 bg-white font-bold text-gray-600 hover:bg-gray-50"
					>›</button
				>
			</div>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
		<div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
			<h2 class="mb-4 text-base font-bold text-gray-800">Status Gizi Balita</h2>
			<div class="relative flex h-48 w-full items-center justify-center">
				<canvas bind:this={chartCanvas}></canvas>
			</div>
		</div>

		<div class="flex flex-col rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
			<h2 class="mb-6 text-base font-bold text-gray-800">Kalkulator Status Gizi</h2>

			<div class="mb-auto grid grid-cols-2 gap-4">
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Berat Badan (kg)</label>
					<input
						type="number"
						placeholder="Contoh: 9.2"
						class="w-full rounded-lg border border-gray-200 p-3 text-sm transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
					/>
				</div>
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Tinggi Badan (cm)</label>
					<input
						type="number"
						placeholder="Contoh: 74"
						class="w-full rounded-lg border border-gray-200 p-3 text-sm transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
					/>
				</div>
			</div>

			<button
				class="mt-6 w-full cursor-pointer rounded-lg border border-gray-200 bg-gray-50 py-3 text-sm font-bold text-gray-500 transition hover:bg-gray-100"
			>
				Masukkan data untuk menghitung status gizi
			</button>
		</div>
	</div>
</div>
