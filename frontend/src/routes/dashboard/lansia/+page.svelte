<script lang="ts">
	import { Search, Plus, Edit2, Trash2, Activity, AlertCircle } from 'lucide-svelte';

	// State pencarian reaktif Svelte 5 Runes
	let searchQuery = $state('');

	// DATA DUMMY: Simulasi respon dari backend API Golang untuk Lansia (FR-05)
	let lansiaList = $state([
		{
			id: 1,
			nik: '3509123456782001',
			nama: 'Slamet Riyadi',
			jk: 'L',
			umur: 65,
			tensi: '120/80',
			status: 'Sehat'
		},
		{
			id: 2,
			nik: '3509123456782002',
			nama: 'Supartini',
			jk: 'P',
			umur: 70,
			tensi: '160/95',
			status: 'Hipertensi'
		},
		{
			id: 3,
			nik: '3509123456782003',
			nama: 'Bambang Pamungkas',
			jk: 'L',
			umur: 68,
			tensi: '130/85',
			status: 'Pemantauan'
		},
		{
			id: 4,
			nik: '3509123456782004',
			nama: 'Kusuma Wardhani',
			jk: 'P',
			umur: 72,
			tensi: '110/70',
			status: 'Sehat'
		},
		{
			id: 5,
			nik: '3509123456782005',
			nama: 'Agus Santoso',
			jk: 'L',
			umur: 69,
			tensi: '150/90',
			status: 'Hipertensi'
		}
	]);

	// RUNES $derived: Memfilter data tabel secara otomatis
	let filteredLansia = $derived(
		lansiaList.filter(
			(l) => l.nama.toLowerCase().includes(searchQuery.toLowerCase()) || l.nik.includes(searchQuery)
		)
	);

	// Fungsi penentu warna lencana berdasarkan status kesehatan Lansia
	function getStatusColor(status: string) {
		switch (status) {
			case 'Hipertensi':
				return 'bg-red-50 text-red-600 border-red-200';
			case 'Pemantauan':
				return 'bg-amber-50 text-amber-600 border-amber-200';
			default:
				return 'bg-teal-50 text-teal-600 border-teal-200';
		}
	}
</script>

<svelte:head>
	<title>Data Lansia - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-end">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-[#1e293b]">Manajemen Data Lansia</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Kelola pendaftaran, pemantauan tanda vital, dan riwayat kesehatan lansia.
			</p>
		</div>
		<button
			class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#0f6456] px-5 py-3 text-sm font-bold text-white shadow-md transition-all hover:bg-[#0c4e43] active:scale-95"
		>
			<Plus size={18} strokeWidth={3} />
			Tambah Lansia
		</button>
	</div>

	<div class="overflow-hidden rounded-[24px] border border-gray-100 bg-white shadow-sm">
		<div
			class="flex flex-col items-start justify-between gap-4 border-b border-gray-100 bg-white p-6 md:flex-row md:items-center"
		>
			<div class="relative w-full md:w-96">
				<div
					class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-gray-400"
				>
					<Search size={18} />
				</div>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Cari nama lansia atau NIK..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-3 pr-4 pl-11 text-sm font-medium text-gray-700 shadow-xs transition-all focus:border-[#0f6456] focus:bg-white focus:ring-2 focus:ring-[#0f6456]/20 focus:outline-none"
				/>
			</div>

			<div
				class="flex flex-shrink-0 items-center gap-3 rounded-xl border border-gray-100 bg-gray-50 px-4 py-2.5 text-sm font-medium text-gray-500"
			>
				<Activity size={20} class="text-[#14a38b]" />
				<span>Total Terdaftar:</span>
				<span class="rounded-lg bg-[#e6f6f4] px-2.5 py-0.5 font-black text-[#0f6456]"
					>{filteredLansia.length}</span
				>
			</div>
		</div>

		<div class="w-full overflow-x-auto">
			<table class="w-full min-w-[900px] border-collapse text-left text-sm text-gray-600">
				<thead
					class="border-b border-gray-100 bg-gray-50/80 text-xs font-bold tracking-wider text-gray-500 uppercase"
				>
					<tr>
						<th class="w-16 px-6 py-5 text-center">No</th>
						<th class="px-6 py-5">NIK / Nama Lansia</th>
						<th class="w-24 px-6 py-5 text-center">L/P</th>
						<th class="w-32 px-6 py-5 text-center">Usia</th>
						<th class="px-6 py-5 text-center">Tensi Darah</th>
						<th class="px-6 py-5 text-center">Status Klinis</th>
						<th class="w-32 px-6 py-5 text-center">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50 bg-white">
					{#if filteredLansia.length === 0}
						<tr>
							<td colspan="7" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<Search size={32} class="mb-3 opacity-50" />
									<p class="text-base font-medium">Tidak ada data Lansia ditemukan.</p>
									<p class="text-xs">
										Kata kunci "{searchQuery}" tidak cocok dengan NIK atau Nama mana pun.
									</p>
								</div>
							</td>
						</tr>
					{/if}

					{#each filteredLansia as lansia, index (lansia.id)}
						<tr class="transition-colors hover:bg-[#f8fdfb]">
							<td class="px-6 py-4 text-center font-bold text-gray-400">{index + 1}</td>
							<td class="px-6 py-4">
								<p class="text-base font-bold text-[#1e293b]">{lansia.nama}</p>
								<p class="mt-0.5 font-mono text-xs text-gray-400">{lansia.nik}</p>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-xs font-bold {lansia.jk ===
									'L'
										? 'bg-blue-50 text-blue-600'
										: 'bg-pink-50 text-pink-600'}"
								>
									{lansia.jk}
								</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="text-lg font-black text-[#0f6456]">{lansia.umur}</span>
								<span class="text-xs font-medium text-gray-500"> Thn</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="font-bold text-gray-700">{lansia.tensi}</span>
								<span class="text-[10px] text-gray-400"> mmHg</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`inline-flex items-center justify-center rounded-lg border px-3 py-1 text-xs font-bold ${getStatusColor(lansia.status)}`}
								>
									{#if lansia.status === 'Hipertensi'}
										<AlertCircle size={14} class="mr-1.5" />
									{/if}
									{lansia.status}
								</span>
							</td>
							<td class="px-6 py-4">
								<div class="flex items-center justify-center gap-2">
									<button
										class="cursor-pointer rounded-lg bg-blue-50 p-2 text-blue-600 transition-colors hover:bg-blue-100"
										aria-label="Edit"
										title="Edit Data"
									>
										<Edit2 size={18} />
									</button>
									<button
										class="cursor-pointer rounded-lg bg-red-50 p-2 text-red-600 transition-colors hover:bg-red-100"
										aria-label="Hapus"
										title="Hapus Data"
									>
										<Trash2 size={18} />
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
