<script>
	import { onMount } from 'svelte';

	// State awal penampil angka ringkasan menggunakan Svelte 5 Runes
	let stats = $state({
		totalBalita: 124,
		totalIbuHamil: 28,
		totalLansia: 65
	});

	let userName = $state('');

	onMount(() => {
		const userData = localStorage.getItem('user');
		if (userData) {
			try {
				userName = JSON.parse(userData).NamaLengkap;
			} catch {
				console.log('error parsing user data');
			}
		}
	});
</script>

<svelte:head>
	<title>Beranda Dasbor - Sistem Posyandu</title>
</svelte:head>

<div class="space-y-6">
	<div
		class="flex flex-col items-start justify-between gap-4 rounded-2xl border border-teal-500/20 bg-gradient-to-r from-teal-600 to-teal-700 p-6 text-white shadow-lg md:flex-row md:items-center"
	>
		<div>
			<h1 class="text-2xl font-bold tracking-tight">Selamat Datang, {userName || 'Petugas'}!</h1>
			<p class="mt-1 text-sm text-teal-100">
				Pantau dan kelola entri rekam medis warga terdaftar secara waktu nyata.
			</p>
		</div>
		<span
			class="rounded-full border border-teal-400/40 bg-teal-500/40 px-3.5 py-1.5 text-xs font-medium text-teal-50 backdrop-blur-xs"
		>
			Akses: Operasional Penuh
		</span>
	</div>

	<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
		<div
			class="flex items-center justify-between rounded-xl border border-gray-100 bg-white p-6 shadow-sm transition-all hover:shadow-md"
		>
			<div>
				<p class="text-xs font-semibold tracking-wider text-gray-400 uppercase">Pasien Balita</p>
				<p class="mt-2 text-3xl font-bold text-gray-900">{stats.totalBalita}</p>
				<p class="mt-1 text-xs font-medium text-teal-600">Status: Pemantauan Aktif</p>
			</div>
			<div class="flex h-14 w-14 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
				<svg class="h-7 w-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
					></path></svg
				>
			</div>
		</div>

		<div
			class="flex items-center justify-between rounded-xl border border-gray-100 bg-white p-6 shadow-sm transition-all hover:shadow-md"
		>
			<div>
				<p class="text-xs font-semibold tracking-wider text-gray-400 uppercase">Ibu Hamil</p>
				<p class="mt-2 text-3xl font-bold text-gray-900">{stats.totalIbuHamil}</p>
				<p class="mt-1 text-xs font-medium text-amber-600">Status: Kehamilan Aktif</p>
			</div>
			<div class="flex h-14 w-14 items-center justify-center rounded-xl bg-amber-50 text-amber-600">
				<svg class="h-7 w-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
					></path></svg
				>
			</div>
		</div>

		<div
			class="flex items-center justify-between rounded-xl border border-gray-100 bg-white p-6 shadow-sm transition-all hover:shadow-md"
		>
			<div>
				<p class="text-xs font-semibold tracking-wider text-gray-400 uppercase">Warga Lansia</p>
				<p class="mt-2 text-3xl font-bold text-gray-900">{stats.totalLansia}</p>
				<p class="mt-1 text-xs font-medium text-purple-600">Status: Terdaftar Rutin</p>
			</div>
			<div
				class="flex h-14 w-14 items-center justify-center rounded-xl bg-purple-50 text-purple-600"
			>
				<svg class="h-7 w-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
					></path></svg
				>
			</div>
		</div>
	</div>

	<div class="rounded-xl border border-gray-100 bg-white p-6 shadow-sm">
		<h2 class="mb-4 text-base font-bold text-gray-900">Pintasan Manajemen Entri</h2>
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
			<div class="rounded-xl border border-gray-100 bg-gray-50 p-4">
				<p class="text-sm font-semibold text-gray-800">Pendaftaran Pasien Balita</p>
				<p class="mt-1 text-xs leading-relaxed text-gray-500">
					Pencatatan identitas NIK dan pengukuran bobot bulanan.
				</p>
			</div>
			<div class="rounded-xl border border-gray-100 bg-gray-50 p-4">
				<p class="text-sm font-semibold text-gray-800">Pencatatan Medis Kandungan</p>
				<p class="mt-1 text-xs leading-relaxed text-gray-500">
					Pengarsipan catatan klinis bidan dan riwayat persalinan.
				</p>
			</div>
			<div class="rounded-xl border border-gray-100 bg-gray-50 p-4">
				<p class="text-sm font-semibold text-gray-800">Pemeriksaan Vital Lansia</p>
				<p class="mt-1 text-xs leading-relaxed text-gray-500">
					Pengawasan kadar glukosa, tensi darah, dan keluhan fisik.
				</p>
			</div>
		</div>
	</div>
</div>
