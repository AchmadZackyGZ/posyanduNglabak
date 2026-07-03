<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { goto } from '$app/navigation';
	import { fetchAPI } from '$lib/api';
	import { ShieldCheck, Lock, User, AlertCircle, ArrowRight } from 'lucide-svelte';

	let username = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let errorMessage = $state('');

	async function handleAdminLogin(event: Event) {
		event.preventDefault();
		isLoading = true;
		errorMessage = '';

		try {
			const res = await fetchAPI('/auth/login', {
				method: 'POST',
				body: JSON.stringify({ username, password })
			});

			// PENGAMANAN GANDA: Cek apakah role benar-benar ADMIN
			const userRole = (res.user.Role || res.user.role || '').toUpperCase();

			if (userRole !== 'ADMIN') {
				errorMessage =
					'Akses Ditolak: Kredensial valid, namun Anda tidak memiliki otoritas Administrator.';
				return;
			}

			// Jika lolos, simpan token dan arahkan ke Dashboard
			localStorage.setItem('token', res.token);
			localStorage.setItem('user', JSON.stringify(res.user));

			// FIX 1: Tambahkan await sebelum goto
			await goto('/dashboard');
		} catch (error: unknown) {
			// FIX 2: Ganti 'any' menjadi 'unknown'
			// FIX 3: Validasi tipe error sebelum mengambil error.message
			errorMessage =
				error instanceof Error ? error.message : 'Gagal terhubung ke server autentikasi.';
		} finally {
			isLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Portal Administrator - POSYANDU Sehat Bersama</title>
</svelte:head>

<!-- Tampilan Login Admin dibuat lebih gelap/profesional (Slate Theme) -->
<div
	class="relative flex min-h-screen items-center justify-center overflow-hidden bg-slate-900 px-4 sm:px-6 lg:px-8"
>
	<!-- Ornamen Background -->
	<div
		class="pointer-events-none absolute top-0 left-0 z-0 h-full w-full overflow-hidden opacity-10"
	>
		<div
			class="absolute -top-[10%] -left-[10%] h-[50%] w-[50%] rounded-full bg-teal-500 blur-[120px]"
		></div>
		<div
			class="absolute top-[60%] right-[5%] h-[40%] w-[30%] rounded-full bg-blue-500 blur-[100px]"
		></div>
	</div>

	<div class="relative z-10 w-full max-w-md space-y-8">
		<div class="text-center">
			<div
				class="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl border border-slate-700 bg-slate-800 shadow-2xl"
			>
				<ShieldCheck size={36} class="text-teal-400" />
			</div>
			<h2 class="mt-6 text-3xl font-black tracking-tight text-white uppercase">Portal Sistem</h2>
			<p class="mt-2 text-sm font-medium text-slate-400">Otorisasi khusus Administrator Posyandu</p>
		</div>

		<div
			class="mt-8 rounded-3xl border border-slate-700 bg-slate-800 p-8 shadow-2xl backdrop-blur-sm"
		>
			{#if errorMessage}
				<div
					class="mb-6 flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-500/10 p-4"
				>
					<AlertCircle size={20} class="mt-0.5 shrink-0 text-red-400" />
					<p class="text-sm font-medium text-red-200">{errorMessage}</p>
				</div>
			{/if}

			<form onsubmit={handleAdminLogin} class="space-y-6">
				<div>
					<label
						for="username"
						class="mb-2 block text-xs font-bold tracking-wider text-slate-300 uppercase"
					>
						ID Administrator
					</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
							<User size={18} class="text-slate-500" />
						</div>
						<input
							id="username"
							type="text"
							required
							bind:value={username}
							placeholder="Masukkan username admin..."
							class="block w-full rounded-xl border border-slate-600 bg-slate-900/50 py-3.5 pr-4 pl-11 text-white placeholder-slate-500 transition focus:border-teal-400 focus:bg-slate-900 focus:ring-1 focus:ring-teal-400 focus:outline-none sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<label
						for="password"
						class="mb-2 block text-xs font-bold tracking-wider text-slate-300 uppercase"
					>
						Kata Sandi Keamanan
					</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
							<Lock size={18} class="text-slate-500" />
						</div>
						<input
							id="password"
							type="password"
							required
							bind:value={password}
							placeholder="••••••••••••"
							class="block w-full rounded-xl border border-slate-600 bg-slate-900/50 py-3.5 pr-4 pl-11 text-white placeholder-slate-500 transition focus:border-teal-400 focus:bg-slate-900 focus:ring-1 focus:ring-teal-400 focus:outline-none sm:text-sm"
						/>
					</div>
				</div>

				<button
					type="submit"
					disabled={isLoading}
					class="group relative flex w-full items-center justify-center gap-2 rounded-xl bg-teal-500 px-4 py-3.5 text-sm font-bold text-slate-900 shadow-lg shadow-teal-500/30 transition-all hover:bg-teal-400 focus:outline-none active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-70"
				>
					{isLoading ? 'MEMVERIFIKASI OTORITAS...' : 'MASUK KE SISTEM'}
					{#if !isLoading}
						<ArrowRight size={18} class="transition-transform group-hover:translate-x-1" />
					{/if}
				</button>
			</form>
		</div>

		<p class="text-center text-xs text-slate-500">
			Halaman ini dilacak dan dilindungi secara ketat.<br />Aktivitas ilegal akan dicatat oleh
			sistem.
		</p>
	</div>
</div>
