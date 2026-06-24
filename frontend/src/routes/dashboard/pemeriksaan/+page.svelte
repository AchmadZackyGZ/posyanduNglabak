<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import {
		Plus,
		Edit2,
		Trash2,
		Calculator,
		X,
		Baby,
		HeartPulse,
		Activity,
		Save
	} from 'lucide-svelte';
	import Chart from 'chart.js/auto';

	// --- 1. INTERFACE ---
	interface Pemeriksaan_balita {
		id: string;
		balita_id: string;
		tanggal_periksa: string;
		berat_badan: number;
		tinggi_badan: number;
		lingkar_kepala: number;
		status_gizi: string;
		catatan: string;
		nama_pasien?: string;
	}
	interface Pemeriksaan_ibu_hamil {
		id: string;
		ibu_hamil_id: string;
		tanggal_periksa: string;
		usia_kehamilan: number;
		tekanan_darah: string;
		berat_badan: number;
		catatan: string;
		nama_pasien?: string;
	}
	interface Pemeriksaan_lansia {
		id: string;
		lansia_id: string;
		tanggal_periksa: string;
		tekanan_darah: string;
		gula_darah: number;
		kolesterol: number;
		catatan: string;
		nama_pasien?: string;
	}

	// --- INTERFACE DATA MASTER ---
	interface MasterBalita {
		id: string;
		nik: string;
		nama_balita: string;
		tanggal_lahir: string;
		jenis_kelamin: string;
		nama_orang_tua: string;
	}

	interface MasterIbuHamil {
		id: string;
		nik: string;
		nama_ibu: string;
		tanggal_lahir: string;
		hpl: string;
		status_kehamilan: string;
	}

	interface MasterLansia {
		id: string;
		nik: string;
		nama_lengkap: string;
		tanggal_lahir: string;
		jenis_kelamin: string;
	}

	// --- 2. STATE TAB & LOADING ---
	let activeTab = $state('balita'); // 'balita' | 'ibu_hamil' | 'lansia'
	let isLoading = $state(true);

	// --- 3. STATE DATA PEMERIKSAAN (TABLE) ---
	let pemeriksaanBalita = $state<Pemeriksaan_balita[]>([]);
	let pemeriksaanIbuHamil = $state<Pemeriksaan_ibu_hamil[]>([]);
	let pemeriksaanLansia = $state<Pemeriksaan_lansia[]>([]);

	// --- 4. STATE DATA MASTER (UNTUK DROPDOWN MODAL) ---
	let listMasterBalita = $state<MasterBalita[]>([]);
	let listMasterIbuHamil = $state<MasterIbuHamil[]>([]);
	let listMasterLansia = $state<MasterLansia[]>([]);

	// --- 5. STATE KALKULATOR BALITA ---
	let isCalcModalOpen = $state(false);
	let beratInput = $state<number | null>(null);
	let tinggiInput = $state<number | null>(null);
	let hasilStatusGizi = $state('');

	// --- 6. STATE "SUPER MODAL" TAMBAH REKAM MEDIS ---
	let isAddModalOpen = $state(false);
	let isSubmitting = $state(false);

	// Form spesifik per entitas
	let formBalita = $state({
		balita_id: '',
		berat_badan: '',
		tinggi_badan: '',
		lingkar_kepala: '',
		status_gizi: 'NORMAL',
		catatan: ''
	});
	let formIbuHamil = $state({
		ibu_hamil_id: '',
		usia_kehamilan: '',
		tekanan_darah: '',
		berat_badan: '',
		catatan: ''
	});
	let formLansia = $state({
		lansia_id: '',
		tekanan_darah: '',
		gula_darah: '',
		kolesterol: '',
		catatan: ''
	});

	// --- FUNGSI AMBIL SEMUA DATA (FRONTEND JOIN & MASTER LIST) ---
	async function loadSemuaPemeriksaan() {
		isLoading = true;
		try {
			// 1. Balita
			const [resBalita, resPeriksaBalita] = await Promise.all([
				fetchAPI('/balita'),
				fetchAPI('/balita/pemeriksaan')
			]);
			if (resBalita.data) listMasterBalita = resBalita.data; // Simpan untuk dropdown
			if (resPeriksaBalita.data && resBalita.data) {
				pemeriksaanBalita = resPeriksaBalita.data.map((p: Pemeriksaan_balita) => {
					const master = resBalita.data.find(
						(b: { id: string; nama_balita: string }) => b.id === p.balita_id
					);
					return { ...p, nama_pasien: master ? master.nama_balita : 'Tidak Diketahui' };
				});
			}

			// 2. Ibu Hamil
			const [resIbuHamil, resPeriksaIbu] = await Promise.all([
				fetchAPI('/ibu-hamil'),
				fetchAPI('/ibu-hamil/pemeriksaan')
			]);
			if (resIbuHamil.data) listMasterIbuHamil = resIbuHamil.data;
			if (resPeriksaIbu.data && resIbuHamil.data) {
				pemeriksaanIbuHamil = resPeriksaIbu.data.map((p: Pemeriksaan_ibu_hamil) => {
					const master = resIbuHamil.data.find(
						(i: { id: string; nama_ibu: string }) => i.id === p.ibu_hamil_id
					);
					return { ...p, nama_pasien: master ? master.nama_ibu : 'Tidak Diketahui' };
				});
			}

			// 3. Lansia
			const [resLansia, resPeriksaLansia] = await Promise.all([
				fetchAPI('/lansia'),
				fetchAPI('/lansia/pemeriksaan')
			]);
			if (resLansia.data) listMasterLansia = resLansia.data;
			if (resPeriksaLansia.data && resLansia.data) {
				pemeriksaanLansia = resPeriksaLansia.data.map((p: Pemeriksaan_lansia) => {
					const master = resLansia.data.find(
						(l: { id: string; nama_lengkap: string }) => l.id === p.lansia_id
					);
					return { ...p, nama_pasien: master ? master.nama_lengkap : 'Tidak Diketahui' };
				});
			}
		} catch (error) {
			console.error('Gagal memuat data pemeriksaan:', error);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadSemuaPemeriksaan();
	});

	// --- FUNGSI SUBMIT REKAM MEDIS ---
	async function handleAddPemeriksaan(event: Event) {
		event.preventDefault();
		isSubmitting = true;

		try {
			if (activeTab === 'balita') {
				const payload = {
					balita_id: formBalita.balita_id,
					berat_badan: parseFloat(formBalita.berat_badan),
					tinggi_badan: parseFloat(formBalita.tinggi_badan),
					lingkar_kepala: formBalita.lingkar_kepala ? parseFloat(formBalita.lingkar_kepala) : 0,
					status_gizi: formBalita.status_gizi,
					catatan: formBalita.catatan
				};
				await fetchAPI('/balita/timbang', { method: 'POST', body: JSON.stringify(payload) });
			} else if (activeTab === 'ibu_hamil') {
				const payload = {
					ibu_hamil_id: formIbuHamil.ibu_hamil_id,
					usia_kehamilan: parseInt(formIbuHamil.usia_kehamilan),
					tekanan_darah: formIbuHamil.tekanan_darah,
					berat_badan: parseFloat(formIbuHamil.berat_badan),
					catatan: formIbuHamil.catatan
				};
				await fetchAPI('/ibu-hamil/periksa', { method: 'POST', body: JSON.stringify(payload) });
			} else if (activeTab === 'lansia') {
				const payload = {
					lansia_id: formLansia.lansia_id,
					tekanan_darah: formLansia.tekanan_darah,
					gula_darah: formLansia.gula_darah ? parseFloat(formLansia.gula_darah) : 0,
					kolesterol: formLansia.kolesterol ? parseFloat(formLansia.kolesterol) : 0,
					catatan: formLansia.catatan
				};
				await fetchAPI('/lansia/periksa', { method: 'POST', body: JSON.stringify(payload) });
			}

			alert('Rekam medis berhasil dicatat!');
			isAddModalOpen = false;

			// Reset Form
			formBalita = {
				balita_id: '',
				berat_badan: '',
				tinggi_badan: '',
				lingkar_kepala: '',
				status_gizi: 'NORMAL',
				catatan: ''
			};
			formIbuHamil = {
				ibu_hamil_id: '',
				usia_kehamilan: '',
				tekanan_darah: '',
				berat_badan: '',
				catatan: ''
			};
			formLansia = {
				lansia_id: '',
				tekanan_darah: '',
				gula_darah: '',
				kolesterol: '',
				catatan: ''
			};

			// Refresh Tabel
			await loadSemuaPemeriksaan();
		} catch (error: unknown) {
			console.error('Gagal menyimpan rekam medis:', error);
			// Pengecekan tipe error secara aman
			const errorMessage = error instanceof Error ? error.message : 'Terjadi kesalahan sistem';
			alert('Gagal menyimpan rekam medis: ' + errorMessage);
		} finally {
			isSubmitting = false;
		}
	}

	// --- HELPER UI ---
	function chartAction(node: HTMLCanvasElement) {
		const chart = new Chart(node, {
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
		return {
			destroy() {
				chart.destroy();
			}
		};
	}

	function hitungGizi() {
		if (!beratInput || !tinggiInput) return;
		const tinggiMeter = tinggiInput / 100;
		const imt = beratInput / (tinggiMeter * tinggiMeter);
		if (imt < 14) hasilStatusGizi = 'Gizi Buruk';
		else if (imt >= 14 && imt < 17) hasilStatusGizi = 'Gizi Kurang';
		else if (imt >= 17 && imt <= 19) hasilStatusGizi = 'Normal';
		else hasilStatusGizi = 'Gizi Lebih';
	}

	function resetKalkulator() {
		beratInput = null;
		tinggiInput = null;
		hasilStatusGizi = '';
		isCalcModalOpen = false;
	}

	function getBadgeColor(status: string) {
		const s = status ? status.toUpperCase() : '';
		if (s.includes('NORMAL')) return 'bg-green-50 text-green-600 border-green-200';
		if (s.includes('KURANG')) return 'bg-amber-50 text-amber-600 border-amber-200';
		if (s.includes('BURUK') || s.includes('STUNTING'))
			return 'bg-red-50 text-red-600 border-red-200';
		return 'bg-blue-50 text-blue-600 border-blue-200';
	}

	function formatDate(isoString: string) {
		if (!isoString) return '-';
		return new Date(isoString).toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Pemeriksaan Klinis - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div
		class="flex flex-col items-start justify-between border-b border-gray-100 pb-4 md:flex-row md:items-end"
	>
		<div>
			<h1 class="text-3xl font-black text-gray-900">Pemeriksaan & Rekam Medis</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Catat dan pantau hasil pengukuran klinis pasien Posyandu.
			</p>
		</div>
		<div class="mt-4 flex gap-3 md:mt-0">
			<button
				onclick={() => (isAddModalOpen = true)}
				class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#117064] px-5 py-2.5 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43]"
			>
				<Plus size={18} strokeWidth={3} /> Tambah Rekam Medis
			</button>
		</div>
	</div>

	<div class="flex gap-2 overflow-x-auto border-b border-gray-100 pb-px">
		<button
			onclick={() => (activeTab = 'balita')}
			class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-bold whitespace-nowrap transition-colors {activeTab ===
			'balita'
				? 'border-[#117064] text-[#117064]'
				: 'border-transparent text-gray-500 hover:border-gray-200 hover:text-gray-700'}"
		>
			<Baby size={18} /> Pemeriksaan Balita
		</button>
		<button
			onclick={() => (activeTab = 'ibu_hamil')}
			class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-bold whitespace-nowrap transition-colors {activeTab ===
			'ibu_hamil'
				? 'border-[#117064] text-[#117064]'
				: 'border-transparent text-gray-500 hover:border-gray-200 hover:text-gray-700'}"
		>
			<HeartPulse size={18} /> Kontrol Ibu Hamil
		</button>
		<button
			onclick={() => (activeTab = 'lansia')}
			class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-bold whitespace-nowrap transition-colors {activeTab ===
			'lansia'
				? 'border-[#117064] text-[#117064]'
				: 'border-transparent text-gray-500 hover:border-gray-200 hover:text-gray-700'}"
		>
			<Activity size={18} /> Tensi & Cek Lansia
		</button>
	</div>

	{#if activeTab === 'balita'}
		<div class="flex justify-end">
			<button
				onclick={() => (isCalcModalOpen = true)}
				class="flex cursor-pointer items-center gap-2 rounded-xl border-2 border-[#117064] bg-white px-4 py-2 text-sm font-bold text-[#117064] transition hover:bg-teal-50"
			>
				<Calculator size={18} /> Kalkulator Gizi
			</button>
		</div>
		<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
			<div class="overflow-x-auto">
				<table class="w-full min-w-[1000px] text-left text-sm text-gray-600">
					<thead
						class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-400 uppercase"
					>
						<tr
							><th class="px-6 py-5 text-center">NO</th><th class="px-6 py-5">NAMA BALITA</th><th
								class="px-6 py-5 text-center">TANGGAL</th
							><th class="px-6 py-5 text-center">BERAT (KG)</th><th class="px-6 py-5 text-center"
								>TINGGI (CM)</th
							><th class="px-6 py-5 text-center">LINGKAR KEPALA</th><th
								class="px-6 py-5 text-center">STATUS GIZI</th
							><th class="px-6 py-5">CATATAN</th><th class="px-6 py-5 text-center">AKSI</th></tr
						>
					</thead>
					<tbody class="divide-y divide-gray-50">
						{#if isLoading}
							<tr
								><td colspan="9" class="px-6 py-10 text-center text-gray-400"
									>Memuat riwayat balita...</td
								></tr
							>
						{:else if pemeriksaanBalita.length === 0}
							<tr
								><td colspan="9" class="px-6 py-10 text-center text-gray-400"
									>Belum ada data pemeriksaan balita.</td
								></tr
							>
						{/if}
						{#each pemeriksaanBalita as item, index (item.id)}
							<tr class="transition-colors hover:bg-gray-50/50">
								<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
								<td class="px-6 py-4 font-bold text-gray-800">{item.nama_pasien}</td>
								<td class="px-6 py-4 text-center">{formatDate(item.tanggal_periksa)}</td>
								<td class="px-6 py-4 text-center font-bold text-gray-700">{item.berat_badan}</td>
								<td class="px-6 py-4 text-center font-bold text-gray-700">{item.tinggi_badan}</td>
								<td class="px-6 py-4 text-center text-gray-600">{item.lingkar_kepala} cm</td>
								<td class="px-6 py-4 text-center"
									><span
										class={`rounded-md border px-2.5 py-1 text-xs font-bold uppercase ${getBadgeColor(item.status_gizi)}`}
										>{item.status_gizi}</span
									></td
								>
								<td class="px-6 py-4 text-gray-500">{item.catatan || '-'}</td>
								<td class="px-6 py-4"
									><div class="flex justify-center gap-2">
										<button
											class="cursor-pointer rounded-lg p-1.5 text-blue-500 transition hover:bg-blue-50"
											><Edit2 size={16} /></button
										><button
											class="cursor-pointer rounded-lg p-1.5 text-red-500 transition hover:bg-red-50"
											><Trash2 size={16} /></button
										>
									</div></td
								>
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
					<canvas use:chartAction></canvas>
				</div>
			</div>
			<div
				class="flex flex-col items-center justify-center rounded-2xl border border-gray-100 bg-teal-50/40 p-8 text-center shadow-sm"
			>
				<div class="mb-4 rounded-full bg-teal-100 p-4 text-[#117064]"><Calculator size={40} /></div>
				<h3 class="text-lg font-black text-gray-800">Kalkulator Gizi Pintar</h3>
				<p class="mt-2 mb-6 max-w-sm text-sm text-gray-500">
					Hitung status gizi instan berdasarkan parameter berat dan tinggi badan standar.
				</p>
				<button
					onclick={() => (isCalcModalOpen = true)}
					class="cursor-pointer rounded-xl bg-[#117064] px-8 py-3.5 text-sm font-bold text-white shadow-lg transition hover:bg-[#0c4e43] active:scale-95"
					>Buka Kalkulator</button
				>
			</div>
		</div>
	{/if}

	{#if activeTab === 'ibu_hamil'}
		<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
			<div class="overflow-x-auto">
				<table class="w-full min-w-[900px] text-left text-sm text-gray-600">
					<thead
						class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-400 uppercase"
					>
						<tr
							><th class="px-6 py-5 text-center">NO</th><th class="px-6 py-5">NAMA IBU</th><th
								class="px-6 py-5 text-center">TANGGAL KONTROL</th
							><th class="px-6 py-5 text-center">USIA KANDUNGAN</th><th
								class="px-6 py-5 text-center">TENSI DARAH</th
							><th class="px-6 py-5 text-center">BERAT (KG)</th><th class="px-6 py-5">CATATAN</th
							><th class="px-6 py-5 text-center">AKSI</th></tr
						>
					</thead>
					<tbody class="divide-y divide-gray-50">
						{#if isLoading}
							<tr
								><td colspan="8" class="px-6 py-10 text-center text-gray-400"
									>Memuat riwayat ibu hamil...</td
								></tr
							>
						{:else if pemeriksaanIbuHamil.length === 0}
							<tr
								><td colspan="8" class="px-6 py-10 text-center text-gray-400"
									>Belum ada data pemeriksaan ibu hamil.</td
								></tr
							>
						{/if}
						{#each pemeriksaanIbuHamil as item, index (item.id)}
							<tr class="transition-colors hover:bg-gray-50/50">
								<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
								<td class="px-6 py-4 font-bold text-gray-800">{item.nama_pasien}</td>
								<td class="px-6 py-4 text-center">{formatDate(item.tanggal_periksa)}</td>
								<td class="px-6 py-4 text-center font-bold text-[#117064]"
									>{item.usia_kehamilan} Mgg</td
								>
								<td class="px-6 py-4 text-center font-bold text-red-500">{item.tekanan_darah}</td>
								<td class="px-6 py-4 text-center font-bold text-gray-700">{item.berat_badan}</td>
								<td class="px-6 py-4 text-gray-500">{item.catatan || '-'}</td>
								<td class="px-6 py-4"
									><div class="flex justify-center gap-2">
										<button
											class="cursor-pointer rounded-lg p-1.5 text-blue-500 transition hover:bg-blue-50"
											><Edit2 size={16} /></button
										><button
											class="cursor-pointer rounded-lg p-1.5 text-red-500 transition hover:bg-red-50"
											><Trash2 size={16} /></button
										>
									</div></td
								>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}

	{#if activeTab === 'lansia'}
		<div class="overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm">
			<div class="overflow-x-auto">
				<table class="w-full min-w-[900px] text-left text-sm text-gray-600">
					<thead
						class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-400 uppercase"
					>
						<tr
							><th class="px-6 py-5 text-center">NO</th><th class="px-6 py-5">NAMA LANSIA</th><th
								class="px-6 py-5 text-center">TANGGAL CEK</th
							><th class="px-6 py-5 text-center">TENSI DARAH</th><th class="px-6 py-5 text-center"
								>GULA DARAH</th
							><th class="px-6 py-5 text-center">KOLESTEROL</th><th class="px-6 py-5">CATATAN</th
							><th class="px-6 py-5 text-center">AKSI</th></tr
						>
					</thead>
					<tbody class="divide-y divide-gray-50">
						{#if isLoading}
							<tr
								><td colspan="8" class="px-6 py-10 text-center text-gray-400"
									>Memuat riwayat lansia...</td
								></tr
							>
						{:else if pemeriksaanLansia.length === 0}
							<tr
								><td colspan="8" class="px-6 py-10 text-center text-gray-400"
									>Belum ada data pemeriksaan lansia.</td
								></tr
							>
						{/if}
						{#each pemeriksaanLansia as item, index (item.id)}
							<tr class="transition-colors hover:bg-gray-50/50">
								<td class="px-6 py-4 text-center font-medium text-gray-400">{index + 1}</td>
								<td class="px-6 py-4 font-bold text-gray-800">{item.nama_pasien}</td>
								<td class="px-6 py-4 text-center">{formatDate(item.tanggal_periksa)}</td>
								<td class="px-6 py-4 text-center font-bold text-red-500">{item.tekanan_darah}</td>
								<td class="px-6 py-4 text-center font-bold text-amber-600"
									>{item.gula_darah ? `${item.gula_darah} mg/dL` : '-'}</td
								>
								<td class="px-6 py-4 text-center font-bold text-blue-600"
									>{item.kolesterol ? `${item.kolesterol} mg/dL` : '-'}</td
								>
								<td class="px-6 py-4 text-gray-500">{item.catatan || '-'}</td>
								<td class="px-6 py-4"
									><div class="flex justify-center gap-2">
										<button
											class="cursor-pointer rounded-lg p-1.5 text-blue-500 transition hover:bg-blue-50"
											><Edit2 size={16} /></button
										><button
											class="cursor-pointer rounded-lg p-1.5 text-red-500 transition hover:bg-red-50"
											><Trash2 size={16} /></button
										>
									</div></td
								>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>

{#if isCalcModalOpen}
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
					title="Tutup"><X size={20} strokeWidth={3} /></button
				>
			</div>
			<div class="space-y-4 p-6">
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Berat Badan (kg)</label><input
						type="number"
						bind:value={beratInput}
						placeholder="Contoh: 9.5"
						class="w-full rounded-xl border border-gray-200 bg-gray-50 p-3 text-sm text-gray-800 shadow-sm transition outline-none focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064]"
					/>
				</div>
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Tinggi Badan (cm)</label
					><input
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
						>Hitung Status Gizi</button
					>
					<p class="mt-3 text-center text-[10px] font-medium text-gray-400">
						Berdasarkan rumus IMT dummy sementara.
					</p>
				</div>
			</div>
		</div>
	</div>
{/if}

{#if isAddModalOpen}
	<div
		class="fixed top-0 left-0 z-[100] flex h-screen w-screen items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
	>
		<div class="absolute inset-0 cursor-pointer" onclick={() => (isAddModalOpen = false)}></div>
		<div
			class="relative z-10 w-full max-w-[500px] overflow-hidden rounded-[24px] bg-white shadow-2xl"
		>
			<div class="flex items-center justify-between border-b border-gray-100 bg-gray-50 px-6 py-4">
				<div class="flex items-center gap-3 text-gray-800">
					<Save size={20} class="text-[#117064]" />
					<h2 class="text-base font-black">
						{#if activeTab === 'balita'}
							Timbang Balita
						{:else if activeTab === 'ibu_hamil'}
							Kontrol Ibu Hamil
						{:else}
							Cek Kesehatan Lansia
						{/if}
					</h2>
				</div>
				<button
					onclick={() => (isAddModalOpen = false)}
					class="cursor-pointer text-gray-400 transition hover:text-red-500"
					title="Tutup"
				>
					<X size={20} strokeWidth={3} />
				</button>
			</div>

			<form onsubmit={handleAddPemeriksaan} class="max-h-[75vh] space-y-4 overflow-y-auto p-6">
				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-700">Pilih Pasien</label>
					{#if activeTab === 'balita'}
						<select
							required
							bind:value={formBalita.balita_id}
							class="w-full rounded-xl border border-gray-200 bg-gray-50 p-3 text-sm text-gray-800 outline-none focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064]"
						>
							<option value="" disabled selected>-- Pilih Balita --</option>
							{#each listMasterBalita as b (b.id)}
								<option value={b.id}>{b.nama_balita} (NIK: {b.nik})</option>
							{/each}
						</select>
					{:else if activeTab === 'ibu_hamil'}
						<select
							required
							bind:value={formIbuHamil.ibu_hamil_id}
							class="w-full rounded-xl border border-gray-200 bg-gray-50 p-3 text-sm text-gray-800 outline-none focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064]"
						>
							<option value="" disabled selected>-- Pilih Ibu Hamil --</option>
							{#each listMasterIbuHamil as i (i.id)}
								<option value={i.id}>{i.nama_ibu} (NIK: {i.nik})</option>
							{/each}
						</select>
					{:else}
						<select
							required
							bind:value={formLansia.lansia_id}
							class="w-full rounded-xl border border-gray-200 bg-gray-50 p-3 text-sm text-gray-800 outline-none focus:border-[#117064] focus:bg-white focus:ring-1 focus:ring-[#117064]"
						>
							<option value="" disabled selected>-- Pilih Lansia --</option>
							{#each listMasterLansia as l (l.id)}
								<option value={l.id}>{l.nama_lengkap} (NIK: {l.nik})</option>
							{/each}
						</select>
					{/if}
				</div>

				<div class="grid grid-cols-2 gap-4">
					{#if activeTab === 'balita'}
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Berat (kg)</label><input
								type="number"
								step="0.1"
								required
								bind:value={formBalita.berat_badan}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Tinggi (cm)</label><input
								type="number"
								step="0.1"
								required
								bind:value={formBalita.tinggi_badan}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Lingkar Kepala (cm)</label
							><input
								type="number"
								step="0.1"
								bind:value={formBalita.lingkar_kepala}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Status Gizi</label>
							<select
								required
								bind:value={formBalita.status_gizi}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							>
								<option value="NORMAL">Normal</option><option value="GIZI KURANG"
									>Gizi Kurang</option
								><option value="GIZI BURUK">Gizi Buruk</option><option value="GIZI LEBIH"
									>Gizi Lebih</option
								>
							</select>
						</div>
						<div class="col-span-2">
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Catatan Bidan/Kader</label
							><textarea
								bind:value={formBalita.catatan}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
								rows="2"
							></textarea>
						</div>
					{:else if activeTab === 'ibu_hamil'}
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700"
								>Usia Kandungan (Minggu)</label
							><input
								type="number"
								required
								bind:value={formIbuHamil.usia_kehamilan}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Tensi Darah</label><input
								type="text"
								required
								bind:value={formIbuHamil.tekanan_darah}
								placeholder="120/80"
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div class="col-span-2">
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Berat Badan (kg)</label
							><input
								type="number"
								step="0.1"
								required
								bind:value={formIbuHamil.berat_badan}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div class="col-span-2">
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Catatan Keluhan</label
							><textarea
								bind:value={formIbuHamil.catatan}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
								rows="2"
							></textarea>
						</div>
					{:else}
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Tensi Darah</label><input
								type="text"
								required
								bind:value={formLansia.tekanan_darah}
								placeholder="130/90"
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Gula Darah (mg/dL)</label
							><input
								type="number"
								step="0.1"
								bind:value={formLansia.gula_darah}
								placeholder="Opsional"
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div class="col-span-2">
							<label class="mb-1.5 block text-xs font-bold text-gray-700">Kolesterol (mg/dL)</label
							><input
								type="number"
								step="0.1"
								bind:value={formLansia.kolesterol}
								placeholder="Opsional"
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
							/>
						</div>
						<div class="col-span-2">
							<label class="mb-1.5 block text-xs font-bold text-gray-700"
								>Catatan Keluhan Klinis</label
							><textarea
								bind:value={formLansia.catatan}
								class="w-full rounded-xl border border-gray-200 bg-white p-3 text-sm outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
								rows="2"
							></textarea>
						</div>
					{/if}
				</div>

				<div class="mt-6 flex justify-end gap-3 border-t border-gray-50 pt-4">
					<button
						type="button"
						onclick={() => (isAddModalOpen = false)}
						class="cursor-pointer rounded-xl px-5 py-2.5 text-sm font-bold text-gray-500 transition hover:bg-gray-100"
						>Batal</button
					>
					<button
						type="submit"
						disabled={isSubmitting}
						class="cursor-pointer rounded-xl bg-[#117064] px-5 py-2.5 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43] disabled:opacity-70"
					>
						{isSubmitting ? 'Menyimpan...' : 'Simpan Rekam Medis'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
