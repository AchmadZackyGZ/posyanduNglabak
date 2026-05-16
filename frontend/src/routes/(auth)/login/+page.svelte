<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { goto } from '$app/navigation';
	
	// 1. Mengimpor berkas gambar logo fisik dari folder assets
	import logoPosyandu from '$lib/assets/logo-posyandu.png';
	
	// 2. Mengimpor ikon elegan dari lucide-svelte
	import { User, ClipboardList, Stethoscope } from 'lucide-svelte';

	let username = $state('');
	let password = $state('');
	let selectedRole = $state('User'); // Tab default kita ubah ke User
	let rememberMe = $state(true);
	let errorMessage = $state('');
	let isLoading = $state(false);

	async function handleLogin(event: Event) {
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
				errorMessage = data.error || 'Username atau kredensial akses tidak valid.';
				isLoading = false;
				return;
			}

			localStorage.setItem('token', data.token);
			localStorage.setItem('user', JSON.stringify(data.user));

			goto('/dashboard');
		} catch {
			errorMessage = 'Koneksi ditolak peladen. Pastikan Go backend aktif dan izin CORS telah terpasang.';
		} finally {
			isLoading = false;
		}
	}

	// 3. Admin kita keluarkan, diganti dengan User sesuai permintaan
	const roles = [
		{ id: 'User', label: 'User' },
		{ id: 'Kader', label: 'Kader' },
		{ id: 'Bidan', label: 'Bidan' }
	];
</script>

<svelte:head>
	<title>Masuk - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-[#14a38b] p-4">
	<div class="w-full max-w-[420px] rounded-[28px] bg-white px-8 py-10 shadow-2xl">
		
		<div class="mx-auto flex h-24 items-center justify-center">
			<img src={logoPosyandu} alt="Logo Posyandu Sehat Bersama" class="h-full w-auto object-contain drop-shadow-sm" />
		</div>

		<div class="mt-4 text-center">
			<p class="mt-0.5 text-sm font-bold text-gray-500">Aplikasi Sehat Bersama</p>
		</div>

		<p class="mt-6 text-center text-xs font-medium text-gray-400">Masuk sebagai</p>

		<div class="mt-2 grid grid-cols-3 gap-2.5">
			{#each roles as role (role.id)}
				<button
					type="button"
					onclick={() => (selectedRole = role.id)}
					class="flex cursor-pointer flex-col items-center justify-center rounded-xl border py-2.5 transition-all {selectedRole === role.id ? 'border-[#0f6456] bg-[#f0fdf4] text-[#0f6456] shadow-xs' : 'border-gray-200 bg-white text-gray-600 hover:bg-gray-50'}"
				>
					{#if role.id === 'User'}
						<User class="mb-1 h-5 w-5" strokeWidth="2.5" />
					{:else if role.id === 'Kader'}
						<ClipboardList class="mb-1 h-5 w-5" strokeWidth="2.5" />
					{:else}
						<Stethoscope class="mb-1 h-5 w-5" strokeWidth="2.5" />
					{/if}
					
					<span class="text-xs font-bold tracking-tight">{role.label}</span>
				</button>
			{/each}
		</div>

		{#if errorMessage}
			<div class="mt-4 rounded-lg border border-red-100 bg-red-50 p-3 text-center">
				<p class="text-xs font-semibold text-red-600">{errorMessage}</p>
			</div>
		{/if}

		<form onsubmit={handleLogin} class="mt-5 space-y-4">
			<div>
				<label for="username" class="block text-xs font-bold text-gray-600">Username</label>
				<input
					id="username"
					type="text"
					bind:value={username}
					required
					placeholder="Masukkan username"
					class="mt-1.5 block w-full rounded-xl border border-gray-300 px-3.5 py-2.5 text-sm shadow-xs outline-none transition-all focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
				/>
			</div>

			<div>
				<label for="password" class="block text-xs font-bold text-gray-600">Password</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					placeholder="••••••••"
					class="mt-1.5 block w-full rounded-xl border border-gray-300 px-3.5 py-2.5 text-sm shadow-xs outline-none transition-all focus:border-[#0f6456] focus:ring-1 focus:ring-[#0f6456]"
				/>
			</div>

			<div class="flex items-center justify-between pt-1">
				<label class="flex cursor-pointer items-center space-x-2">
					<input
						type="checkbox"
						bind:checked={rememberMe}
						class="h-4 w-4 rounded border-gray-300 text-[#0f6456] focus:ring-[#0f6456]"
					/>
					<span class="text-xs font-semibold text-gray-600">Ingat saya</span>
				</label>

				<a href="#lupa" class="text-xs font-bold text-[#0f6456] transition-colors hover:underline">
					Lupa password?
				</a>
			</div>

			<button
				type="submit"
				disabled={isLoading}
				class="mt-2 flex w-full cursor-pointer items-center justify-center rounded-xl bg-[#0f6456] py-3 text-sm font-bold text-white shadow-md transition-all hover:bg-[#0c4e43] active:scale-[0.99] disabled:opacity-60"
			>
				{#if isLoading}
					<svg class="mr-2 h-4 w-4 animate-spin text-white" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Memverifikasi...
				{:else}
					Masuk
				{/if}
			</button>
		</form>

		<div class="mt-8 border-t border-gray-100 pt-4 text-center">
			<p class="text-[10px] text-gray-400">
				© 2026 Posyandu Sehat Bersama. All rights reserved.
			</p>
		</div>

	</div>
</div>