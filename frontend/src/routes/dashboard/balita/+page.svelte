<script lang="ts">
	import { Search, Plus, Edit2, Trash2, Baby } from 'lucide-svelte';

	// State pencarian reaktif Svelte 5
	let searchQuery = $state('');

	// DATA DUMMY: Simulasi respon dari backend API Golang
	let balitaList = $state([
		{
			id: 1,
			nik: '3509123456780001',
			nama: 'Budi Santoso',
			jk: 'L',
			namaIbu: 'Siti Aminah',
			umurBulan: 14
		},
		{
			id: 2,
			nik: '3509123456780002',
			nama: 'Aisyah Putri',
			jk: 'P',
			namaIbu: 'Nurul Huda',
			umurBulan: 8
		},
		{
			id: 3,
			nik: '3509123456780003',
			nama: 'Cakra Manggala',
			jk: 'L',
			namaIbu: 'Dewi Lestari',
			umurBulan: 22
		},
		{
			id: 4,
			nik: '3509123456780004',
			nama: 'Dinda Kirana',
			jk: 'P',
			namaIbu: 'Rina Wati',
			umurBulan: 5
		},
		{
			id: 5,
			nik: '3509123456780005',
			nama: 'Eko Prayitno',
			jk: 'L',
			namaIbu: 'Sri Wahyuni',
			umurBulan: 30
		}
	]);

	// RUNES $derived: Memfilter data secara otomatis saat pencarian berubah
	let filteredBalita = $derived(
		balitaList.filter(
			(b) => b.nama.toLowerCase().includes(searchQuery.toLowerCase()) || b.nik.includes(searchQuery)
		)
	);
</script>

<svelte:head>
	<title>Data Balita - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-end">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-[#1e293b]">Manajemen Data Balita</h1>
			<p class="mt-1.5 text-sm font-medium text-gray-500">
				Kelola pendaftaran, identitas, dan riwayat penimbangan balita.
			</p>
		</div>
		<button
			class="flex cursor-pointer items-center gap-2 rounded-xl bg-[#0f6456] px-5 py-3 text-sm font-bold text-white shadow-md transition-all hover:bg-[#0c4e43] active:scale-95"
		>
			<Plus size={18} strokeWidth={3} />
			Tambah Balita
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
					placeholder="Cari nama atau NIK balita..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50 py-3 pr-4 pl-11 text-sm font-medium text-gray-700 shadow-xs transition-all focus:border-[#0f6456] focus:bg-white focus:ring-2 focus:ring-[#0f6456]/20 focus:outline-none"
				/>
			</div>

			<div
				class="flex flex-shrink-0 items-center gap-3 rounded-xl border border-gray-100 bg-gray-50 px-4 py-2.5 text-sm font-medium text-gray-500"
			>
				<Baby size={20} class="text-[#14a38b]" />
				<span>Total Data:</span>
				<span class="rounded-lg bg-[#e6f6f4] px-2.5 py-0.5 font-black text-[#0f6456]"
					>{filteredBalita.length}</span
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
						<th class="px-6 py-5">NIK / Nama Balita</th>
						<th class="w-24 px-6 py-5 text-center">L/P</th>
						<th class="w-32 px-6 py-5 text-center">Usia</th>
						<th class="px-6 py-5">Nama Ibu</th>
						<th class="w-32 px-6 py-5 text-center">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-50 bg-white">
					{#if filteredBalita.length === 0}
						<tr>
							<td colspan="6" class="px-6 py-16 text-center">
								<div class="flex flex-col items-center justify-center text-gray-400">
									<Search size={32} class="mb-3 opacity-50" />
									<p class="text-base font-medium">Tidak ada data balita ditemukan.</p>
									<p class="text-xs">
										Kata kunci "{searchQuery}" tidak cocok dengan NIK atau Nama mana pun.
									</p>
								</div>
							</td>
						</tr>
					{/if}

					{#each filteredBalita as balita, index (balita.id)}
						<tr class="transition-colors hover:bg-[#f8fdfb]">
							<td class="px-6 py-4 text-center font-bold text-gray-400">{index + 1}</td>
							<td class="px-6 py-4">
								<p class="text-base font-bold text-[#1e293b]">{balita.nama}</p>
								<p class="mt-0.5 font-mono text-xs text-gray-400">{balita.nik}</p>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-xs font-bold {balita.jk ===
									'L'
										? 'bg-blue-50 text-blue-600'
										: 'bg-pink-50 text-pink-600'}"
								>
									{balita.jk}
								</span>
							</td>
							<td class="px-6 py-4 text-center">
								<span class="font-bold text-gray-700">{balita.umurBulan}</span>
								<span class="text-xs text-gray-500">Bln</span>
							</td>
							<td class="px-6 py-4">
								<span class="font-medium text-gray-700">{balita.namaIbu}</span>
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
