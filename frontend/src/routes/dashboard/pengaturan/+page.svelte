<script lang="ts">
	// --- STATE MANAGEMENT FORM ---
	let formProfil = $state({
		namaPosyandu: 'Posyandu Sehat Bersama',
		alamat: 'Jl. Melati No. 5, RW 04, Kelurahan Harapan',
		nomorHp: '0812-3456-7890',
		wilayahKerja: 'Kelurahan Harapan, Kecamatan Sejahtera'
	});

	// --- STATE MANAGEMENT TOGGLE (SVELTE MURNI) ---
	let pengaturanSistem = $state({
		notifJadwal: true,
		notifImunisasi: true,
		alertStunting: true,
		backupOtomatis: true,
		modeGelap: false
	});

	function simpanPerubahan() {
		alert('Perubahan profil berhasil disimpan secara lokal!');
	}
</script>

<svelte:head>
	<title>Pengaturan - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<div class="border-b border-gray-100 pb-4">
		<h1 class="text-2xl font-black text-gray-900">Pengaturan</h1>
	</div>

	<div class="grid grid-cols-1 items-start gap-6 md:grid-cols-2">
		<div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
			<h2 class="mb-6 border-b border-gray-100 pb-4 text-base font-bold text-gray-800">
				Profil Posyandu
			</h2>

			<div class="space-y-5">
				<div>
					<label for="namaPosyandu" class="mb-2 block text-sm font-bold text-gray-700"
						>Nama Posyandu</label
					>
					<input
						type="text"
						id="namaPosyandu"
						bind:value={formProfil.namaPosyandu}
						class="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
					/>
				</div>

				<div>
					<label for="alamat" class="mb-2 block text-sm font-bold text-gray-700">Alamat</label>
					<textarea
						id="alamat"
						bind:value={formProfil.alamat}
						rows="3"
						class="w-full resize-none rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
					></textarea>
				</div>

				<div>
					<label for="nomorHp" class="mb-2 block text-sm font-bold text-gray-700"
						>Nomor HP Kader</label
					>
					<input
						type="text"
						id="nomorHp"
						bind:value={formProfil.nomorHp}
						class="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
					/>
				</div>

				<div>
					<label for="wilayahKerja" class="mb-2 block text-sm font-bold text-gray-700"
						>Wilayah Kerja</label
					>
					<input
						type="text"
						id="wilayahKerja"
						bind:value={formProfil.wilayahKerja}
						class="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-800 transition outline-none focus:border-[#117064] focus:ring-1 focus:ring-[#117064]"
					/>
				</div>

				<div class="pt-2">
					<button
						onclick={simpanPerubahan}
						class="cursor-pointer rounded-xl bg-[#117064] px-6 py-3 text-sm font-bold text-white shadow-md transition hover:bg-[#0c4e43] active:scale-95"
					>
						Simpan Perubahan
					</button>
				</div>
			</div>
		</div>

		<div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
			<h2 class="mb-6 border-b border-gray-100 pb-4 text-base font-bold text-gray-800">
				Notifikasi & Sistem
			</h2>

			<div class="mb-8">
				<h3 class="mb-4 text-xs font-black tracking-widest text-gray-500 uppercase">Notifikasi</h3>
				<div class="space-y-6">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm font-bold text-gray-800">Notifikasi Jadwal</p>
							<p class="text-xs font-medium text-gray-500">Ingatkan jadwal posyandu mendatang</p>
						</div>
						<button
							type="button"
							role="switch"
							aria-label="Notifikasi Jadwal"
							aria-checked={pengaturanSistem.notifJadwal}
							onclick={() => (pengaturanSistem.notifJadwal = !pengaturanSistem.notifJadwal)}
							class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none {pengaturanSistem.notifJadwal
								? 'bg-[#117064]'
								: 'bg-gray-200'}"
						>
							<span
								class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {pengaturanSistem.notifJadwal
									? 'translate-x-5'
									: 'translate-x-0'}"
							></span>
						</button>
					</div>

					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm font-bold text-gray-800">Notifikasi Imunisasi</p>
							<p class="text-xs font-medium text-gray-500">Ingatkan jadwal imunisasi balita</p>
						</div>
						<button
							type="button"
							role="switch"
							aria-label="Notifikasi Imunisasi"
							aria-checked={pengaturanSistem.notifImunisasi}
							onclick={() => (pengaturanSistem.notifImunisasi = !pengaturanSistem.notifImunisasi)}
							class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none {pengaturanSistem.notifImunisasi
								? 'bg-[#117064]'
								: 'bg-gray-200'}"
						>
							<span
								class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {pengaturanSistem.notifImunisasi
									? 'translate-x-5'
									: 'translate-x-0'}"
							></span>
						</button>
					</div>

					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm font-bold text-gray-800">Alert Stunting</p>
							<p class="text-xs font-medium text-gray-500">Notifikasi balita berisiko stunting</p>
						</div>
						<button
							type="button"
							role="switch"
							aria-label="Alert Stunting"
							aria-checked={pengaturanSistem.alertStunting}
							onclick={() => (pengaturanSistem.alertStunting = !pengaturanSistem.alertStunting)}
							class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none {pengaturanSistem.alertStunting
								? 'bg-[#117064]'
								: 'bg-gray-200'}"
						>
							<span
								class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {pengaturanSistem.alertStunting
									? 'translate-x-5'
									: 'translate-x-0'}"
							></span>
						</button>
					</div>
				</div>
			</div>

			<div>
				<h3 class="mb-4 text-xs font-black tracking-widest text-gray-500 uppercase">Sistem</h3>
				<div class="space-y-6">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm font-bold text-gray-800">Backup Otomatis</p>
							<p class="text-xs font-medium text-gray-500">Backup data setiap minggu</p>
						</div>
						<button
							type="button"
							role="switch"
							aria-label="Backup Otomatis"
							aria-checked={pengaturanSistem.backupOtomatis}
							onclick={() => (pengaturanSistem.backupOtomatis = !pengaturanSistem.backupOtomatis)}
							class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none {pengaturanSistem.backupOtomatis
								? 'bg-[#117064]'
								: 'bg-gray-200'}"
						>
							<span
								class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {pengaturanSistem.backupOtomatis
									? 'translate-x-5'
									: 'translate-x-0'}"
							></span>
						</button>
					</div>

					<div class="flex items-center justify-between opacity-50">
						<div>
							<p class="text-sm font-bold text-gray-800">Mode Gelap</p>
							<p class="text-xs font-medium text-gray-500">Tampilan gelap (coming soon)</p>
						</div>
						<button
							type="button"
							role="switch"
							aria-label="Mode Gelap"
							aria-checked={pengaturanSistem.modeGelap}
							disabled
							class="relative inline-flex h-6 w-11 shrink-0 cursor-not-allowed items-center rounded-full border-2 border-transparent bg-gray-200 transition-colors duration-200 ease-in-out focus:outline-none"
						>
							<span
								class="pointer-events-none inline-block h-5 w-5 translate-x-0 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
							></span>
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
