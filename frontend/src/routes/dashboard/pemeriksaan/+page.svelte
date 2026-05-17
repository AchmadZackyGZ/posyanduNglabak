<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, Edit2, Trash2, Calculator, X } from 'lucide-svelte';
	import Chart from 'chart.js/auto';

	let chartCanvas: HTMLCanvasElement;
	let chartInstance: Chart | null = null;

	// --- STATE MODAL & KALKULATOR ---
	let isModalOpen = $state(false);
	let beratInput = $state<number | null>(null);
	let tinggiInput = $state<number | null>(null);
	let hasilStatusGizi = $state('');

	// DATA DUMMY Pemeriksaan Balita (Kolom Lengkap)
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

	function hitungGizi() {
		if (!beratInput || !tinggiInput) return;
		const score = beratInput / (tinggiInput / 10);
		if (score < 1.1) hasilStatusGizi = 'Gizi Buruk';
		else if (score < 1.3) hasilStatusGizi = 'Gizi Kurang';
		else if (score < 1.8) hasilStatusGizi = 'Normal';
		else hasilStatusGizi = 'Gizi Lebih';
	}

	function resetKalkulator() {
		beratInput = null;
		tinggiInput = null;
		hasilStatusGizi = '';
		isModalOpen = false;
	}

	onMount(() => {
		if (chartCanvas) {
			chartInstance = new Chart(chartCanvas, {
				type: 'doughnut',
				data: {
					labels: ['Normal', 'Gizi Kurang', 'Gizi Lebih', 'Gizi Buruk'],
					datasets: [
						{
							data: [62, 15, 8, 2],
							backgroundColor: ['#22c55e', '#f59e0b', '#3b82f6', '#ef4444'],
							borderWidth: 0
						}
					]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false,
					cutout: '70%',
					plugins: { legend: { position: 'right' } }
				}
			});
		}
		return () => {
			if (chartInstance) chartInstance.destroy();
		};
	});

	function getBadgeColor(status: string) {
		if (status === 'Normal') return 'bg-green-50 text-green-600 border-green-200';
		if (status === 'Kurang' || status === 'Gizi Kurang')
			return 'bg-amber-50 text-amber-600 border-amber-200';
		if (status === 'Gizi Buruk') return 'bg-red-50 text-red-600 border-red-200';
		return 'bg-blue-50 text-blue-600 border-blue-200';
	}
</script>

<svelte:head>
	<title>Pemeriksaan - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div
		class="flex flex-col items-start justify-between border-b border-gray-100 pb-4 md:flex-row md:items-end"
	>
		<div>
			<h1 class="text-3xl font-black text-gray-900">Pemeriksaan Balita</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Catat hasil pengukuran fisik dan status gizi balita hari ini.
			</p>
		</div>
		<div class="mt-4 flex gap-3 md:mt-0">
			<button
				onclick={() => (isModalOpen = true)}
				class="flex cursor-pointer items-center gap-2 rounded-xl border-2 border-[#117064] bg-white px-4 py-2.5 text-sm font-bold text-[#117064] transition hover:bg-teal-50"
			>
				<Calculator size={18} />
				Kalkulator Gizi
			</button>
			<button
				class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#117064] px-4 py-2.5 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43]"
			>
				<Plus size={18} strokeWidth={3} />
				Tambah Data
			</button>
		</div>
	</div>

	<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
		<div class="overflow-x-auto">
			<table class="w-full min-w-[1000px] text-left text-sm text-gray-600">
				<thead
					class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-400 uppercase"
				>
					<tr>
						<th class="px-6 py-5 text-center">NO</th>
						<th class="px-6 py-5">NAMA BALITA</th>
						<th class="px-6 py-5 text-center">TANGGAL</th>
						<th class="px-6 py-5 text-center">BERAT (KG)</th>
						<th class="px-6 py-5 text-center">TINGGI (CM)</th>
						<th class="px-6 py-5 text-center">LINGKAR KEPALA</th>
						<th class="px-6 py-5 text-center">STATUS GIZI</th>
						<th class="px-6 py-5">CATATAN</th>
						<th class="px-6 py-5 text-center">AKSI</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50">
					{#each pemeriksaanList as item, index (item.id)}
						<tr class="transition-colors hover:bg-gray-50/50">
							<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
							<td class="px-6 py-4 font-bold text-gray-800">{item.nama}</td>
							<td class="px-6 py-4 text-center">{item.tanggal}</td>
							<td class="px-6 py-4 text-center font-bold text-gray-700">{item.berat}</td>
							<td class="px-6 py-4 text-center font-bold text-gray-700">{item.tinggi}</td>
							<td class="px-6 py-4 text-center text-gray-600">{item.lingkarKepala} cm</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`rounded-md border px-2.5 py-1 text-xs font-bold uppercase ${getBadgeColor(item.statusGizi)}`}
								>
									{item.statusGizi}
								</span>
							</td>
							<td class="px-6 py-4 text-gray-500">{item.catatan}</td>
							<td class="px-6 py-4">
								<div class="flex justify-center gap-2">
									<button
										class="cursor-pointer rounded-lg p-1.5 text-blue-500 transition hover:bg-blue-50"
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
	</div>

	<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
		<div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
			<h2 class="mb-4 text-base font-bold text-gray-800">Sebaran Gizi Balita</h2>
			<div class="relative flex h-64 w-full items-center justify-center">
				<canvas bind:this={chartCanvas}></canvas>
			</div>
		</div>

		<div
			class="flex flex-col items-center justify-center rounded-2xl border border-gray-100 bg-teal-50/40 p-8 text-center shadow-sm"
		>
			<div class="mb-4 rounded-full bg-teal-100 p-4 text-[#117064]">
				<Calculator size={40} />
			</div>
			<h3 class="text-lg font-black text-gray-800">Kalkulator Gizi Pintar</h3>
			<p class="mt-2 mb-6 max-w-sm text-sm text-gray-500">
				Gunakan alat ini untuk menghitung status gizi balita secara instan berdasarkan parameter
				berat dan tinggi badan standar.
			</p>
			<button
				onclick={() => (isModalOpen = true)}
				class="cursor-pointer rounded-xl bg-[#117064] px-8 py-3.5 text-sm font-bold text-white shadow-lg transition hover:bg-[#0c4e43] active:scale-95"
			>
				Buka Kalkulator
			</button>
		</div>
	</div>
</div>

{#if isModalOpen}
	<div
		class="fixed top-0 left-0 z-[100] flex h-screen w-screen items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
	>
		<div class="absolute inset-0" onclick={resetKalkulator}></div>

		<div
			class="relative z-10 w-full max-w-[400px] overflow-hidden rounded-[24px] bg-white shadow-2xl"
		>
			<div class="flex items-center justify-between border-b border-teal-100 bg-teal-50 px-6 py-4">
				<div class="flex items-center gap-3 text-[#117064]">
					<Calculator size={20} />
					<h2 class="text-base font-black">Kalkulator Gizi</h2>
				</div>
				<button
					onclick={resetKalkulator}
					class="cursor-pointer text-teal-600/50 transition hover:text-red-500"
					title="Tutup"
				>
					<X size={20} strokeWidth={3} />
				</button>
			</div>

			<div class="space-y-4 p-6">
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Berat Badan (kg)</label>
					<input
						type="number"
						bind:value={beratInput}
						placeholder="Contoh: 9.5"
						class="w-full rounded-xl border border-gray-200 bg-gray-50 p-3 text-sm text-gray-800 shadow-sm transition outline-none focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064]"
					/>
				</div>

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Tinggi Badan (cm)</label>
					<input
						type="number"
						bind:value={tinggiInput}
						placeholder="Contoh: 75"
						class="w-full rounded-xl border border-gray-200 bg-gray-50 p-3 text-sm text-gray-800 shadow-sm transition outline-none focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064]"
					/>
				</div>

				{#if hasilStatusGizi}
					<div
						class="mt-4 rounded-xl border-2 border-dashed border-[#117064]/30 bg-[#f0fdf4] p-4 text-center"
					>
						<p class="text-[10px] font-bold tracking-widest text-gray-500 uppercase">
							Hasil Analisis
						</p>
						<p class="mt-0.5 text-xl font-black text-[#117064]">{hasilStatusGizi}</p>
					</div>
				{/if}

				<div class="pt-2">
					<button
						onclick={hitungGizi}
						disabled={!beratInput || !tinggiInput}
						class="w-full cursor-pointer rounded-xl bg-[#117064] py-3 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43] active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
					>
						Hitung Status Gizi
					</button>
					<p class="mt-3 text-center text-[10px] font-medium text-gray-400">
						Berdasarkan standar dummy antropometri Posyandu.
					</p>
				</div>
			</div>
		</div>
	</div>
{/if}
