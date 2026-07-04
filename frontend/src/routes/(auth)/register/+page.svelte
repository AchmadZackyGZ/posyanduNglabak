<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */ // Supaya tidak muncul warning saat menggunakan goto() untuk navigasi
	import { goto } from '$app/navigation';
	import { fetchAPI } from '$lib/api';
	import {
		UserCircle,
		Phone,
		Lock,
		HeartPulse,
		ArrowRight,
		AlertCircle,
		CheckCircle2
	} from 'lucide-svelte';

	let nama_lengkap = $state('');
	let username = $state('');
	let password = $state('');

	let isLoading = $state(false);
	let errorMessage = $state('');
	let successMessage = $state('');

	async function handleRegister(event: Event) {
		event.preventDefault();
		isLoading = true;
		errorMessage = '';
		successMessage = '';

		try {
			await fetchAPI('/auth/register', {
				method: 'POST',
				body: JSON.stringify({ nama_lengkap, username, password })
			});

			successMessage = 'Yey! Pendaftaran Bunda berhasil. Bersiap menuju halaman masuk...';

			// Jeda 2 detik agar user bisa membaca pesan sukses, lalu alihkan ke halaman login
			setTimeout(() => {
				goto('/login'); // Arahkan ke rute yang baru Anda buat
			}, 2000);
		} catch (error: unknown) {
			errorMessage = (error as Error).message || 'Maaf, terjadi kesalahan saat mendaftar.';
		} finally {
			isLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Daftar Posyandu - Sehat Bersama</title>
</svelte:head>

<!-- BACKGROUND MEGAH & MODERN -->
<div
	class="relative flex min-h-screen items-center justify-center overflow-hidden bg-gradient-to-br from-teal-50 via-white to-emerald-100 p-4 font-sans"
>
	<!-- Elemen Dekorasi Latar Belakang (Blob) -->
	<div class="absolute -top-20 -left-20 h-72 w-72 rounded-full bg-teal-300/20 blur-3xl"></div>
	<div
		class="absolute -right-32 -bottom-32 h-96 w-96 rounded-full bg-emerald-400/20 blur-3xl"
	></div>

	<!-- KARTU REGISTRASI (GLASSMORPHISM) -->
	<div
		class="relative z-10 w-full max-w-md overflow-hidden rounded-[2rem] border border-white/60 bg-white/70 p-8 shadow-2xl shadow-teal-900/5 backdrop-blur-xl sm:p-10"
	>
		<!-- HEADER KARTU -->
		<div class="mb-8 text-center">
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-tr from-[#117064] to-teal-400 text-white shadow-lg shadow-teal-500/30"
			>
				<HeartPulse size={32} strokeWidth={2.5} />
			</div>
			<h1
				class="mb-2 flex items-center justify-center gap-2 text-2xl font-extrabold tracking-tight text-gray-800"
			>
				<span>Daftar Posyandu</span>
			</h1>
			<p class="text-sm text-gray-500">
				Mari bergabung untuk pantau tumbuh kembang si kecil dengan lebih mudah dan menyenangkan.
			</p>
		</div>

		<!-- NOTIFIKASI ERROR / SUKSES -->
		{#if errorMessage}
			<div
				class="animate-in fade-in slide-in-from-top-2 mb-6 flex items-start gap-3 rounded-2xl border border-red-200 bg-red-50 p-4 text-red-600"
			>
				<AlertCircle size={20} class="mt-0.5 shrink-0" />
				<p class="text-sm leading-relaxed font-medium">{errorMessage}</p>
			</div>
		{/if}

		{#if successMessage}
			<div
				class="animate-in fade-in slide-in-from-top-2 mb-6 flex items-start gap-3 rounded-2xl border border-teal-200 bg-teal-50 p-4 text-[#117064]"
			>
				<CheckCircle2 size={20} class="mt-0.5 shrink-0" />
				<p class="text-sm leading-relaxed font-bold">{successMessage}</p>
			</div>
		{/if}

		<!-- FORMULIR -->
		<form onsubmit={handleRegister} class="space-y-5">
			<!-- Input Nama Lengkap -->
			<div class="space-y-1.5">
				<label class="ml-1 text-xs font-bold tracking-wide text-gray-600 uppercase" for="nama"
					>Nama Lengkap Bunda</label
				>
				<div class="relative">
					<div
						class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-gray-400"
					>
						<UserCircle size={20} />
					</div>
					<input
						id="nama"
						type="text"
						required
						bind:value={nama_lengkap}
						disabled={isLoading || successMessage !== ''}
						placeholder="Misal: Siti Aminah"
						class="block w-full rounded-2xl border border-white bg-white/80 py-3.5 pr-4 pl-12 text-sm text-gray-800 shadow-sm transition placeholder:text-gray-400 focus:border-[#117064] focus:bg-white focus:ring-4 focus:ring-teal-500/10 focus:outline-none disabled:opacity-60"
					/>
				</div>
			</div>

			<!-- Input Nomor HP (Username) -->
			<div class="space-y-1.5">
				<label class="ml-1 text-xs font-bold tracking-wide text-gray-600 uppercase" for="hp"
					>Nomor HP / WhatsApp</label
				>
				<div class="relative">
					<div
						class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-gray-400"
					>
						<Phone size={20} />
					</div>
					<input
						id="hp"
						type="tel"
						required
						bind:value={username}
						disabled={isLoading || successMessage !== ''}
						placeholder="081234567890"
						class="block w-full rounded-2xl border border-white bg-white/80 py-3.5 pr-4 pl-12 text-sm text-gray-800 shadow-sm transition placeholder:text-gray-400 focus:border-[#117064] focus:bg-white focus:ring-4 focus:ring-teal-500/10 focus:outline-none disabled:opacity-60"
					/>
				</div>
			</div>

			<!-- Input Password -->
			<div class="space-y-1.5">
				<label class="ml-1 text-xs font-bold tracking-wide text-gray-600 uppercase" for="password"
					>Kata Sandi (Password)</label
				>
				<div class="relative">
					<div
						class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4 text-gray-400"
					>
						<Lock size={20} />
					</div>
					<input
						id="password"
						type="password"
						required
						minlength="6"
						bind:value={password}
						disabled={isLoading || successMessage !== ''}
						placeholder="Minimal 6 karakter rahasia"
						class="block w-full rounded-2xl border border-white bg-white/80 py-3.5 pr-4 pl-12 text-sm text-gray-800 shadow-sm transition placeholder:text-gray-400 focus:border-[#117064] focus:bg-white focus:ring-4 focus:ring-teal-500/10 focus:outline-none disabled:opacity-60"
					/>
				</div>
			</div>

			<!-- Tombol Submit -->
			<button
				type="submit"
				disabled={isLoading || successMessage !== ''}
				class="group mt-8 flex w-full cursor-pointer items-center justify-center gap-2 rounded-2xl bg-gradient-to-r from-[#117064] to-teal-500 py-3.5 text-sm font-bold text-white shadow-lg shadow-teal-600/30 transition-all hover:scale-[1.02] hover:shadow-teal-600/40 active:scale-[0.98] disabled:pointer-events-none disabled:opacity-70"
			>
				{#if isLoading}
					<span class="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
					></span>
					Mendaftarkan...
				{:else if successMessage}
					Berhasil!
				{:else}
					Daftar Sekarang <ArrowRight
						size={18}
						class="transition-transform group-hover:translate-x-1"
					/>
				{/if}
			</button>
		</form>

		<!-- Link ke Login -->
		<p class="mt-8 text-center text-sm text-gray-500">
			Sudah punya akun Posyandu?
			<a href="/" class="font-bold text-[#117064] transition hover:text-teal-500 hover:underline">
				Masuk di sini
			</a>
		</p>
	</div>
</div>
