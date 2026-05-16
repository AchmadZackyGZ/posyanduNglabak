<script lang="ts">
	import { Search, Plus, Edit2, Trash2, HeartPulse, AlertCircle } from 'lucide-svelte';

	// State pencarian reaktif Svelte 5
	let searchQuery = $state('');

	// DATA DUMMY: Simulasi respon dari backend API Golang (FR-04)
	let ibuHamilList = $state([
		{
			id: 1,
			nik: '3509123456781001',
			namaIbu: 'Siti Aminah',
			noHp: '081234567890',
			hpl: '2026-08-15',
			usiaKandungan: 24,
			status: 'Normal'
		},
		{
			id: 2,
			nik: '3509123456781002',
			namaIbu: 'Nurul Huda',
			noHp: '081298765432',
			hpl: '2026-06-10',
			usiaKandungan: 34,
			status: 'Risiko Tinggi'
		},
		{
			id: 3,
			nik: '3509123456781003',
			namaIbu: 'Dewi Lestari',
			noHp: '085712341234',
			hpl: '2026-11-20',
			usiaKandungan: 10,
			status: 'Normal'
		},
		{
			id: 4,
			nik: '3509123456781004',
			namaIbu: 'Rina Wati',
			noHp: '087856785678',
			hpl: '2026-07-05',
			usiaKandungan: 30,
			status: 'Pemantauan'
		},
		{
			id: 5,
			nik: '3509123456781005',
			namaIbu: 'Sri Wahyuni',
			noHp: '081911223344',
			hpl: '2026-09-12',
			usiaKandungan: 20,
			status: 'Normal'
		}
	]);

	// RUNES $derived: Memfilter tabel secara otomatis (Real-time Search)
	let filteredIbuHamil = $derived(
		ibuHamilList.filter(
			(ibu) =>
				ibu.namaIbu.toLowerCase().includes(searchQuery.toLowerCase()) ||
				ibu.nik.includes(searchQuery)
		)
	);

	// Fungsi helper untuk menentukan warna lencana (badge) status kehamilan
	function getStatusColor(status: string) {
		switch (status) {
			case 'Risiko Tinggi':
				return 'bg-red-50 text-red-600 border-red-200';
			case 'Pemantauan':
				return 'bg-amber-50 text-amber-600 border-amber-200';
			default:
				return 'bg-teal-50 text-teal-600 border-teal-200';
		}
	}
</script>

<svelte:head>
	<title>Data Ibu Hamil - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-end">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-[#1e293b]">Manajemen Ibu Hamil</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Kelola pendaftaran kandungan, taksiran persalinan, dan rekam medis klinis.
			</p>
		</div>
		<button
			class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#0f6456] px-5 py-3 text-sm font-bold text-white shadow-md transition-all hover:bg-[#0c4e43] active:scale-95"
		>
			<Plus size={18} strokeWidth={3} />
			Tambah Ibu Hamil
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
					placeholder="Cari nama ibu atau NIK..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-3 pr-4 pl-11 text-sm font-medium text-gray-700 shadow-xs transition-all focus:border-[#0f6456] focus:bg-white focus:ring-2 focus:ring-[#0f6456]/20 focus:outline-none"
				/>
			</div>

			<div
				class="flex flex-shrink-0 items-center gap-3 rounded-xl border border-gray-100 bg-gray-50 px-4 py-2.5 text-sm font-medium text-gray-500"
			>
				<HeartPulse size={20} class="text-[#14a38b]" />
				<span>Total Kandungan:</span>
				<span class="rounded-lg bg-[#e6f6f4] px-2.5 py-0.5 font-black text-[#0f6456]"
					>{filteredIbuHamil.length}</span
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
						<th class="px-6 py-5">NIK / Nama Ibu Hamil</th>
						<th class="px-6 py-5 text-center">Kontak (HP)</th>
						<th class="px-6 py-5 text-center">HPL (Perkiraan)</th>
						<th class="px-6 py-5 text-center">Usia Kandungan</th>
						<th class="px-6 py-5 text-center">Status</th>
						<th class="w-32 px-6 py-5 text-center">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50 bg-white">
					{#if filteredIbuHamil.length === 0}
						<tr>
							<td colspan="7" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<Search size={32} class="mb-3 opacity-50" />
									<p class="text-base font-medium">Tidak ada data Ibu Hamil ditemukan.</p>
									<p class="text-xs">
										Kata kunci "{searchQuery}" tidak cocok dengan NIK atau Nama mana pun.
									</p>
								</div>
							</td>
						</tr>
					{/if}

					{#each filteredIbuHamil as ibu, index (ibu.id)}
						<tr class="transition-colors hover:bg-[#f8fdfb]">
							<td class="px-6 py-4 text-center font-bold text-gray-400">{index + 1}</td>
							<td class="px-6 py-4">
								<p class="text-base font-bold text-[#1e293b]">{ibu.namaIbu}</p>
								<p class="mt-0.5 font-mono text-xs text-gray-400">{ibu.nik}</p>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="font-medium text-gray-700">{ibu.noHp}</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="font-bold text-gray-700">{ibu.hpl}</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="text-lg font-black text-[#0f6456]">{ibu.usiaKandungan}</span>
								<span class="text-xs font-medium text-gray-500"> Mgg</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`inline-flex items-center justify-center rounded-lg border px-3 py-1 text-xs font-bold ${getStatusColor(ibu.status)}`}
								>
									{#if ibu.status === 'Risiko Tinggi'}
										<AlertCircle size={14} class="mr-1.5" />
									{/if}
									{ibu.status}
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
