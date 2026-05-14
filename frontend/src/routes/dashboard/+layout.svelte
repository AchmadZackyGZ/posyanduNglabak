<script>
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	/** @type {any} */
	let user = $state(null);
	let isSidebarOpen = $state(false);

	// Logika Proteksi Klien: Mencegat peramban saat komponen dimuat
	onMount(() => {
		const token = localStorage.getItem('token');
		const userData = localStorage.getItem('user');

		// Jika gembok sesi (JWT) tidak ada, lempar kembali ke luar
		if (!token) {
			goto('/login');
			return;
		}

		if (userData) {
			try {
				user = JSON.parse(userData);
			} catch {
				user = { NamaLengkap: 'Petugas Aktif', Role: 'KADER' };
			}
		}
	});

	// Pembersihan sesi saat tombol keluar ditekan
	function handleLogout() {
		localStorage.removeItem('token');
		localStorage.removeItem('user');
		goto('/login');
	}

	// Variabel reaktif untuk menyalakan indikator link aktif di sidebar menggunakan Svelte 5 Runes
	let currentPath = $derived($page.url.pathname);
</script>

<div class="flex min-h-screen flex-col bg-gray-50 md:flex-row">
	<div class="flex items-center justify-between bg-teal-700 p-4 text-white shadow-md md:hidden">
		<span class="text-lg font-bold tracking-wide">Posyandu v1.0</span>
		<button
			on:click={() => (isSidebarOpen = !isSidebarOpen)}
			aria-label="Buka Tutup Navigasi"
			class="focus:outline-none"
		>
			<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d={isSidebarOpen ? 'M6 18L18 6M6 6l12 12' : 'M4 6h16M4 12h16M4 18h16'}
				></path>
			</svg>
		</button>
	</div>

	<aside
		class="{isSidebarOpen
			? 'block'
			: 'hidden'} w-full flex-shrink-0 border-r border-gray-200 bg-white shadow-sm transition-all duration-300 md:block md:w-64"
	>
		<div class="hidden border-b border-gray-100 p-6 md:block">
			<h2 class="text-xl font-bold tracking-wide text-teal-700">Posyandu v1.0</h2>
			<p class="mt-1 text-xs text-gray-400">Sistem Informasi Medis</p>
		</div>

		<div class="border-b border-gray-100 bg-teal-50/50 p-4 md:hidden">
			<p class="text-xs text-gray-500">Petugas Aktif:</p>
			<p class="text-sm font-semibold text-teal-800">{user?.NamaLengkap || 'Memuat...'}</p>
			<span
				class="mt-1 inline-block rounded-full bg-teal-600 px-2 py-0.5 text-[10px] font-medium text-white"
				>{user?.Role || 'KADER'}</span
			>
		</div>

		<nav class="space-y-1.5 p-4">
			<a
				href="/dashboard"
				class="flex items-center rounded-lg px-4 py-3 text-sm font-medium transition-all {currentPath ===
				'/dashboard'
					? 'bg-teal-50 font-semibold text-teal-700'
					: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
			>
				<svg
					class="mr-3 h-5 w-5 {currentPath === '/dashboard' ? 'text-teal-600' : 'text-gray-400'}"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
					></path></svg
				>
				Beranda
			</a>

			<a
				href="/dashboard/balita"
				class="flex items-center rounded-lg px-4 py-3 text-sm font-medium transition-all {currentPath.startsWith(
					'/dashboard/balita'
				)
					? 'bg-teal-50 font-semibold text-teal-700'
					: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
			>
				<svg
					class="mr-3 h-5 w-5 {currentPath.startsWith('/dashboard/balita')
						? 'text-teal-600'
						: 'text-gray-400'}"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
					></path></svg
				>
				Data Balita
			</a>

			<a
				href="/dashboard/ibu-hamil"
				class="flex items-center rounded-lg px-4 py-3 text-sm font-medium transition-all {currentPath.startsWith(
					'/dashboard/ibu-hamil'
				)
					? 'bg-teal-50 font-semibold text-teal-700'
					: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
			>
				<svg
					class="mr-3 h-5 w-5 {currentPath.startsWith('/dashboard/ibu-hamil')
						? 'text-teal-600'
						: 'text-gray-400'}"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
					></path></svg
				>
				Ibu Hamil
			</a>

			<a
				href="/dashboard/lansia"
				class="flex items-center rounded-lg px-4 py-3 text-sm font-medium transition-all {currentPath.startsWith(
					'/dashboard/lansia'
				)
					? 'bg-teal-50 font-semibold text-teal-700'
					: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
			>
				<svg
					class="mr-3 h-5 w-5 {currentPath.startsWith('/dashboard/lansia')
						? 'text-teal-600'
						: 'text-gray-400'}"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
					></path></svg
				>
				Data Lansia
			</a>
		</nav>

		<div class="mt-auto border-t border-gray-100 p-4">
			<button
				on:click={handleLogout}
				class="flex w-full items-center justify-center rounded-lg px-4 py-2.5 text-sm font-medium text-red-600 transition-colors hover:bg-red-50"
			>
				<svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
					></path></svg
				>
				Keluar Sistem
			</button>
		</div>
	</aside>

	<main class="flex min-w-0 flex-1 flex-col overflow-y-auto">
		<header
			class="hidden h-16 items-center justify-between border-b border-gray-200 bg-white px-8 shadow-xs md:flex"
		>
			<div class="text-sm font-medium text-gray-500">
				<span>Sistem Pengelolaan Medis</span>
			</div>

			<div class="flex items-center space-x-3">
				<div class="text-right">
					<p class="text-sm font-semibold text-gray-800">{user?.NamaLengkap || 'Petugas Aktif'}</p>
					<p class="text-xs font-bold text-teal-600">{user?.Role || 'KADER'}</p>
				</div>
				<div
					class="flex h-10 w-10 items-center justify-center rounded-full border border-teal-200 bg-teal-100 font-bold text-teal-700"
				>
					{user?.NamaLengkap ? user.NamaLengkap.charAt(0).toUpperCase() : 'P'}
				</div>
			</div>
		</header>

		<div class="flex-1 p-6 md:p-8">
			<slot />
		</div>
	</main>
</div>
