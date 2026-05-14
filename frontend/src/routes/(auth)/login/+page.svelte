<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { goto } from '$app/navigation';

	// Menggunakan Svelte 5 Runes untuk reaktivitas state formulir
	let username = $state('');
	let password = $state('');
	let errorMessage = $state('');
	let isLoading = $state(false);

	async function handleLogin(event: Event) {
		// Mencegah perilaku bawaan peramban memuat ulang halaman
		event.preventDefault();
		isLoading = true;
		errorMessage = '';

		try {
			const res = await fetch('http://localhost:8080/api/v1/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			const data = await res.json();

			if (!res.ok) {
				errorMessage = data.error || 'Username atau password salah.';
				isLoading = false;
				return;
			}

			// Menyimpan token JWT dan identitas pengguna di penyimpanan lokal peramban
			localStorage.setItem('token', data.token);
			localStorage.setItem('user', JSON.stringify(data.user));

			// Alihkan pengguna ke halaman utama dasbor operasional
			goto('/dashboard');
		} catch {
			errorMessage = 'Tidak dapat terhubung ke peladen backend. Pastikan peladen Go aktif.';
		} finally {
			isLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Masuk - Sistem Posyandu Terintegrasi</title>
</svelte:head>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 to-teal-50 p-4"
>
	<div class="w-full max-w-md rounded-2xl border border-gray-100 bg-white p-8 shadow-xl">
		<div class="text-center">
			<h1 class="text-3xl font-bold text-gray-900">Sistem Posyandu</h1>
			<p class="mt-2 text-sm text-gray-500">
				Masuk sebagai Admin atau Kader untuk mengelola rekam medis
			</p>
		</div>

		{#if errorMessage}
			<div class="mt-6 rounded-md border-l-4 border-red-500 bg-red-50 p-4 transition-all">
				<p class="text-sm font-medium text-red-700">{errorMessage}</p>
			</div>
		{/if}

		<form onsubmit={handleLogin} class="mt-6 space-y-4">
			<div>
				<label for="username" class="block text-sm font-medium text-gray-700"
					>Username / Akun Petugas</label
				>
				<input
					id="username"
					type="text"
					bind:value={username}
					required
					placeholder="Masukkan username (contoh: admin)"
					class="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-3 shadow-sm transition-all focus:border-teal-500 focus:ring-2 focus:ring-teal-500 focus:outline-none"
				/>
			</div>

			<div>
				<label for="password" class="block text-sm font-medium text-gray-700">Password</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					placeholder="••••••••"
					class="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-3 shadow-sm transition-all focus:border-teal-500 focus:ring-2 focus:ring-teal-500 focus:outline-none"
				/>
			</div>

			<button
				type="submit"
				disabled={isLoading}
				class="flex w-full cursor-pointer items-center justify-center rounded-lg border border-transparent bg-teal-600 px-4 py-3 text-sm font-medium text-white shadow-md transition-all hover:bg-teal-700 focus:ring-2 focus:ring-teal-500 focus:ring-offset-2 focus:outline-none disabled:opacity-50"
			>
				{#if isLoading}
					<svg
						class="mr-3 -ml-1 h-5 w-5 animate-spin text-white"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
					>
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
						></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						></path>
					</svg>
					Memverifikasi...
				{:else}
					Masuk ke Sistem
				{/if}
			</button>
		</form>
	</div>
</div>
