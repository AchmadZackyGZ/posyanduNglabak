<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import { Baby, LineChart, Calendar, Award, User, Sparkles, AlertCircle } from 'lucide-svelte';

	interface Pemeriksaan {
		id: string;
		tanggal_periksa: string;
		berat_badan: number;
		tinggi_badan: number;
		lingkar_kepala: number;
		status_gizi: string;
		catatan: string;
	}

	interface BalitaKMS {
		id: string;
		nama_balita: string;
		nik: string;
		tanggal_lahir: string;
		jenis_kelamin: string;
		nama_orang_tua: string;
		pemeriksaans: Pemeriksaan[];
	}

	// 1. Interface untuk data mentah dari API (pemeriksaans bisa null)
	interface BalitaResponse {
		id: string;
		nama_balita: string;
		nik: string;
		tanggal_lahir: string;
		jenis_kelamin: string;
		nama_orang_tua: string;
		pemeriksaans: Pemeriksaan[] | null;
	}

	let dataKMS = $state<BalitaKMS[]>([]);
	let selectedAnakIndex = $state(0);
	let isLoading = $state(true);
	let errorMessage = $state('');

	let chartCanvas = $state<HTMLCanvasElement | null>(null);

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	let chartInstance: any = null;

	async function loadKMSData() {
		try {
			const res = await fetchAPI('/warga/kms');

			// TypeScript sekarang tahu persis tipe data 'anak' dan hasil return-nya
			dataKMS = (res.data || []).map(
				(anak: BalitaResponse): BalitaKMS => ({
					...anak,
					pemeriksaans: anak.pemeriksaans || []
				})
			);
		} catch (error: unknown) {
			errorMessage =
				error instanceof Error
					? error.message || 'Gagal memuat rekam medis KMS Anda.'
					: 'Gagal memuat rekam medis KMS Anda.';
			console.error(error);
		} finally {
			isLoading = false;
		}
	}

	function hitungUmurBulan(tanggalLahirStr: string, tanggalPeriksaStr: string): number {             
		const lahir = new Date(tanggalLahirStr);
		const periksa = new Date(tanggalPeriksaStr);
		return (
			(periksa.getFullYear() - lahir.getFullYear()) * 12 + (periksa.getMonth() - lahir.getMonth())
		);
	}

	function renderChart(anak: BalitaKMS) {
		if (!chartCanvas || !anak.pemeriksaans || anak.pemeriksaans.length === 0) return;

		if (chartInstance) {
			chartInstance.destroy();
		}

		const labels = anak.pemeriksaans.map((p) => {
			const bulan = hitungUmurBulan(anak.tanggal_lahir, p.tanggal_periksa);
			return `${bulan} Bln`;
		});

		const dataBB = anak.pemeriksaans.map((p) => p.berat_badan);
		const dataTB = anak.pemeriksaans.map((p) => p.tinggi_badan);

		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		chartInstance = new (window as any).Chart(chartCanvas, {
			type: 'line',
			data: {
				labels: labels,
				datasets: [
					{
						label: 'Berat Badan (kg)',
						data: dataBB,
						borderColor: '#117064',
						backgroundColor: 'rgba(17, 112, 100, 0.1)',
						borderWidth: 3,
						tension: 0.3,
						fill: true,
						yAxisID: 'yBB'
					},
					{
						label: 'Tinggi Badan (cm)',
						data: dataTB,
						borderColor: '#3b82f6',
						backgroundColor: 'transparent',
						borderWidth: 3,
						tension: 0.3,
						yAxisID: 'yTB'
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					legend: {
						position: 'top',
						labels: { font: { weight: 'bold', family: 'Plus Jakarta Sans' } }
					}
				},
				scales: {
					yBB: {
						type: 'linear',
						position: 'left',
						title: { display: true, text: 'Berat Badan (kg)', font: { weight: 'bold' } }
					},
					yTB: {
						type: 'linear',
						position: 'right',
						title: { display: true, text: 'Tinggi Badan (cm)', font: { weight: 'bold' } },
						grid: { drawOnChartArea: false }
					}
				}
			}
		});
	}

	onMount(async () => {
		await loadKMSData();
	});

	$effect(() => {
		if (dataKMS.length > 0 && dataKMS[selectedAnakIndex]) {
			setTimeout(() => renderChart(dataKMS[selectedAnakIndex]), 50);
		}
	});

	let anakTerpilih = $derived(dataKMS[selectedAnakIndex] || null);

	// FIX: Pelindung untuk pemeriksaan terakhir agar tidak error jika array kosong
	let pemeriksaanTerakhir = $derived(
		anakTerpilih?.pemeriksaans && anakTerpilih.pemeriksaans.length > 0
			? anakTerpilih.pemeriksaans[anakTerpilih.pemeriksaans.length - 1]
			: null
	);
</script>

<svelte:head>
	<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
	<title>KMS Digital Bunda - POSYANDU Sehat Bersama</title>
</svelte:head>

<div
	class="min-h-screen bg-gradient-to-br from-teal-50/50 via-white to-blue-50/30 p-4 font-sans md:p-8"
>
	<div
		class="mb-8 flex flex-col items-start justify-between gap-4 rounded-3xl border border-gray-100 bg-white p-6 shadow-md shadow-teal-900/5 sm:flex-row sm:items-center"
	>
		<div class="flex items-center gap-4">
			<div
				class="flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-tr from-[#117064] to-teal-400 text-white shadow-md"
			>
				<Baby size={28} />
			</div>
			<div>
				<h1 class="flex items-center gap-2 text-xl font-extrabold text-gray-800">
					Selamat Datang Bunda! <Sparkles class="text-amber-400" size={20} />
				</h1>
				<p class="text-sm text-gray-500">
					Pantau Kartu Menuju Sehat (KMS) Digital buah hati Anda secara real-time.
				</p>
			</div>
		</div>

		{#if dataKMS.length > 1}
			<div class="flex w-full items-center gap-2 sm:w-auto">
				<label class="shrink-0 text-xs font-bold text-gray-500 uppercase" for="select-anak"
					>Pilih Anak:</label
				>
				<select
					id="select-anak"
					bind:value={selectedAnakIndex}
					class="w-full rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm font-bold text-gray-700 focus:border-[#117064] focus:ring-2 focus:ring-teal-100 focus:outline-none"
				>
					{#each dataKMS as anak, index (anak.id)}
						<option value={index}>{anak.nama_balita}</option>
					{/each}
				</select>
			</div>
		{/if}
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-20 font-medium text-gray-400">
			Memuat KMS Digital Bunda...
		</div>
	{:else if errorMessage}
		<div class="flex gap-3 rounded-2xl border border-red-200 bg-red-50 p-4 text-red-600">
			<AlertCircle size={20} class="shrink-0" />
			<p class="text-sm font-medium">{errorMessage}</p>
		</div>
	{:else if dataKMS.length === 0}
		<div
			class="mx-auto max-w-md rounded-3xl border-2 border-dashed border-gray-200 bg-white p-8 py-16 text-center"
		>
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-gray-400"
			>
				<Baby size={32} />
			</div>
			<h3 class="mb-1 text-lg font-bold text-gray-700">Data Balita Belum Ditautkan</h3>
			<p class="text-sm text-gray-500">
				Akun Anda belum terhubung dengan rekam medis posyandu. Silakan hubungi Kader/Bidan saat
				kunjungan Posyandu berikutnya.
			</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
			<div class="space-y-6">
				<div
					class="relative overflow-hidden rounded-3xl border border-gray-100 bg-white p-6 shadow-md"
				>
					<div class="absolute top-0 right-0 h-24 w-24 rounded-bl-[4rem] bg-teal-500/5"></div>
					<h3
						class="mb-4 flex items-center gap-2 text-sm font-bold tracking-wider text-gray-400 uppercase"
					>
						<User size={16} /> Profil Buah Hati
					</h3>
					<div class="space-y-3">
						<div>
							<div class="text-xs text-gray-400">Nama Lengkap</div>
							<div class="text-lg font-extrabold text-gray-800">{anakTerpilih.nama_balita}</div>
						</div>
						<div class="grid grid-cols-2 gap-4 pt-2">
							<div>
								<div class="text-xs text-gray-400">Jenis Kelamin</div>
								<div class="text-sm font-bold text-gray-700">
									{anakTerpilih.jenis_kelamin === 'L' ? 'Laki-Laki ♂' : 'Perempuan ♀'}
								</div>
							</div>
							<div>
								<div class="text-xs text-gray-400">Tanggal Lahir</div>
								<div class="text-sm font-bold text-gray-700">
									{new Date(anakTerpilih.tanggal_lahir).toLocaleDateString('id-ID', {
										dateStyle: 'medium'
									})}
								</div>
							</div>
						</div>
						<div class="pt-2">
							<div class="text-xs text-gray-400">NIK</div>
							<div class="font-mono text-sm font-bold text-gray-600">{anakTerpilih.nik}</div>
						</div>
					</div>
				</div>

				<div
					class="rounded-3xl bg-gradient-to-tr from-[#117064] to-teal-600 p-6 text-white shadow-xl"
				>
					<h3
						class="mb-4 flex items-center gap-2 text-xs font-bold tracking-wider text-teal-200 uppercase"
					>
						<Award size={16} /> Status Timbang Terakhir
					</h3>
					{#if pemeriksaanTerakhir}
						<div class="py-2 text-center">
							<div class="text-xs text-teal-100">Kondisi Status Gizi</div>
							<div class="mt-1 animate-pulse text-3xl font-black tracking-wide">
								{pemeriksaanTerakhir.status_gizi}
							</div>
						</div>
						<div class="mt-4 grid grid-cols-3 gap-2 border-t border-white/20 pt-4 text-center">
							<div>
								<div class="text-[10px] text-teal-100">Berat</div>
								<div class="text-base font-extrabold">{pemeriksaanTerakhir.berat_badan} kg</div>
							</div>
							<div>
								<div class="text-[10px] text-teal-100">Tinggi</div>
								<div class="text-base font-extrabold">{pemeriksaanTerakhir.tinggi_badan} cm</div>
							</div>
							<div>
								<div class="text-[10px] text-teal-100">L. Kepala</div>
								<div class="text-base font-extrabold">
									{pemeriksaanTerakhir.lingkar_kepala || '-'} cm
								</div>
							</div>
						</div>
						{#if pemeriksaanTerakhir.catatan}
							<div class="mt-4 rounded-xl border border-white/10 bg-white/10 p-3 text-xs">
								<span class="font-bold">Pesan Bidan/Kader:</span>
								{pemeriksaanTerakhir.catatan}
							</div>
						{/if}
					{:else}
						<div class="py-6 text-center text-sm text-teal-100">
							Belum ada rekam jejak pemeriksaan masuk. Kader belum menginput data timbang.
						</div>
					{/if}
				</div>
			</div>

			<div class="space-y-6 lg:col-span-2">
				<div class="rounded-3xl border border-gray-100 bg-white p-6 shadow-md">
					<div class="mb-4 flex items-center justify-between">
						<h3
							class="flex items-center gap-2 text-sm font-bold tracking-wider text-gray-700 uppercase"
						>
							<LineChart size={18} class="text-[#117064]" /> Kurva Tumbuh Kembang Anak
						</h3>
						<span class="text-xs font-medium text-gray-400">Grafik BB & TB Real-Time</span>
					</div>

					<div class="relative h-80 w-full">
						{#if anakTerpilih.pemeriksaans && anakTerpilih.pemeriksaans.length > 0}
							<canvas bind:this={chartCanvas}></canvas>
						{:else}
							<div
								class="flex h-full flex-col items-center justify-center rounded-2xl border border-dashed border-gray-200 bg-gray-50 p-6 text-center text-sm text-gray-400"
							>
								<LineChart size={32} class="mb-2 text-gray-300" />
								<p>Grafik akan muncul setelah Kader mengisi minimal 1 kali penimbangan.</p>
							</div>
						{/if}
					</div>
				</div>

				<div class="overflow-hidden rounded-3xl border border-gray-100 bg-white p-6 shadow-md">
					<h3
						class="mb-4 flex items-center gap-2 text-sm font-bold tracking-wider text-gray-700 uppercase"
					>
						<Calendar size={18} class="text-blue-500" /> Histori Timbangan Bulanan
					</h3>
					<div class="overflow-x-auto">
						<table class="w-full text-left text-sm">
							<thead class="bg-gray-50 text-xs font-bold text-gray-500 uppercase">
								<tr>
									<th class="px-4 py-3">Usia</th>
									<th class="px-4 py-3">Tanggal Timbang</th>
									<th class="px-4 py-3 text-center">Berat (kg)</th>
									<th class="px-4 py-3 text-center">Tinggi (cm)</th>
									<th class="px-4 py-3 text-center">Status</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-gray-100">
								{#if anakTerpilih.pemeriksaans && anakTerpilih.pemeriksaans.length > 0}
									{#each [...anakTerpilih.pemeriksaans].reverse() as p (p.id)}
										<tr class="transition hover:bg-gray-50/50">
											<td class="px-4 py-3 font-bold text-[#117064]"
												>{hitungUmurBulan(anakTerpilih.tanggal_lahir, p.tanggal_periksa)} Bulan</td
											>
											<td class="px-4 py-3 text-gray-500"
												>{new Date(p.tanggal_periksa).toLocaleDateString('id-ID', {
													dateStyle: 'medium'
												})}</td
											>
											<td class="px-4 py-3 text-center font-bold text-gray-700">{p.berat_badan}</td>
											<td class="px-4 py-3 text-center font-bold text-gray-700">{p.tinggi_badan}</td
											>
											<td class="px-4 py-3 text-center">
												<span
													class="inline-flex rounded-full border border-teal-100 bg-teal-50 px-2.5 py-0.5 text-xs font-bold text-[#117064]"
												>
													{p.status_gizi}
												</span>
											</td>
										</tr>
									{/each}
								{:else}
									<tr>
										<td colspan="5" class="py-6 text-center text-gray-400"
											>Belum ada riwayat timbangan yang dicatat.</td
										>
									</tr>
								{/if}
							</tbody>
						</table>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
