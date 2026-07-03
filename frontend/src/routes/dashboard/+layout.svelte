<script lang="ts">
	/* eslint-disable svelte/no-navigation-without-resolve */
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	// Tambahkan ikon Package untuk Inventaris dan UserPlus/Users untuk spesifik menu
	import {
		LayoutDashboard,
		Baby,
		HeartPulse,
		Activity,
		Stethoscope,
		CalendarDays,
		FileText,
		Users,
		UserPlus,
		Settings,
		LogOut,
		Menu,
		Bell,
		Package // Ikon untuk Inventaris
	} from 'lucide-svelte';

	import '../layout.css';
	let { children } = $props();

	let user = $state<{ nama_lengkap: string; role: string } | null>(null);
	let isSidebarOpen = $state(false);

	onMount(() => {
		const token = localStorage.getItem('token');
		const userData = localStorage.getItem('user');

		if (!token) {
			goto('/login');
			return;
		}

		if (userData) {
			try {
				user = JSON.parse(userData);
			} catch {
				user = { nama_lengkap: 'Petugas Aktif', role: 'KADER' };
			}
		} else {
			user = { nama_lengkap: 'Petugas Aktif', role: 'KADER' };
		}
	});

	function handleLogout() {
		localStorage.removeItem('token');
		localStorage.removeItem('user');
		goto('/login');
	}

	let currentPath = $derived($page.url.pathname);
	let currentRole = $derived((user?.role || 'KADER').toUpperCase());

	// --- STRUKTUR MENU BARU BERBASIS KELOMPOK (GROUPING) ---
	let menuGroups = $derived([
		{
			group: 'Utama',
			roles: ['KADER', 'BIDAN'],
			items: [
				{ path: '/dashboard', label: 'Beranda', icon: LayoutDashboard },
				{ path: '/dashboard/jadwal', label: 'Jadwal Kegiatan', icon: CalendarDays }
			]
		},
		{
			group: 'Data Pasien & Layanan',
			roles: ['KADER', 'BIDAN'],
			items: [
				{ path: '/dashboard/balita', label: 'Data Balita', icon: Baby },
				{ path: '/dashboard/ibu-hamil', label: 'Data Ibu Hamil', icon: HeartPulse },
				{ path: '/dashboard/lansia', label: 'Data Lansia', icon: Activity },
				{ path: '/dashboard/pemeriksaan', label: 'Pemeriksaan', icon: Stethoscope }
			]
		},
		{
			group: 'Logistik',
			roles: ['KADER', 'BIDAN'], // <--- FASE 2: HAK AKSES INVENTARIS
			items: [{ path: '/dashboard/inventaris', label: 'Inventaris Obat', icon: Package }]
		},
		{
			group: 'Manajemen Sistem',
			roles: ['ADMIN'], // Khusus Admin
			items: [
				{ path: '/dashboard/laporan', label: 'Laporan Keseluruhan', icon: FileText }, // Admin lihat ini di atas
				{ path: '/dashboard/pengaturan', label: 'Pengaturan Sistem', icon: Settings }
			]
		},
		{
			group: 'Pengguna Internal',
			roles: ['ADMIN'],
			items: [
				// Arahkan ke file +page.svelte yang sudah kita buat sebelumnya
				{ path: '/dashboard/pengguna', label: 'Kelola Staf / Kader', icon: UserPlus }
			]
		},
		{
			group: 'Pengguna Publik',
			roles: ['ADMIN'],
			items: [
				// Rute baru untuk nanti saat kita buat fitur warga
				{ path: '/dashboard/pengguna-publik', label: 'Daftar Warga (User)', icon: Users }
			]
		},
		{
			// Laporan untuk Kader/Bidan ditaruh di bawah agar tidak mencolok
			group: 'Pelaporan',
			roles: ['KADER', 'BIDAN'],
			items: [{ path: '/dashboard/laporan', label: 'Cetak Laporan', icon: FileText }]
		}
	]);

	// Filter grup menu berdasarkan Role
	let visibleGroups = $derived(menuGroups.filter((g) => g.roles.includes(currentRole)));

	// Format tanggal untuk Header
	let todayDate = new Date().toLocaleDateString('id-ID', {
		weekday: 'long',
		day: 'numeric',
		month: 'long',
		year: 'numeric'
	});
</script>

<div class="flex min-h-screen bg-[#f4f7f6]">
	<aside
		class="{isSidebarOpen
			? 'translate-x-0'
			: '-translate-x-full'} fixed inset-y-0 left-0 z-50 flex w-64 flex-col bg-[#117064] text-white transition-transform duration-300 ease-in-out md:relative md:translate-x-0"
	>
		<div class="flex items-center gap-3 border-b border-white/10 px-6 py-5">
			<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-white text-[#117064]">
				<Activity size={20} strokeWidth={3} />
			</div>
			<div>
				<h2 class="text-sm font-bold tracking-wider">POSYANDU</h2>
				<p class="text-[10px] text-teal-100">Sehat Bersama</p>
			</div>
		</div>

		<nav class="scrollbar-hide flex-1 space-y-4 overflow-y-auto px-3 py-4">
			{#each visibleGroups as group (group.group)}
				<div>
					<!-- Header Grup (Label Kecil Transparan) -->
					<h3 class="mb-1 px-3 text-[10px] font-bold tracking-widest text-teal-200/70 uppercase">
						{group.group}
					</h3>
					<div class="space-y-1">
						{#each group.items as menu (menu.path)}
							<a
								href={menu.path}
								class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors {currentPath ===
									menu.path ||
								(menu.path !== '/dashboard' && currentPath.startsWith(menu.path))
									? 'bg-white/20 text-white shadow-sm'
									: 'text-teal-100 hover:bg-white/10 hover:text-white'}"
							>
								<menu.icon size={18} />
								{menu.label}
							</a>
						{/each}
					</div>
				</div>
			{/each}
		</nav>

		<div class="border-t border-white/10 p-4">
			<div class="flex items-center gap-3 rounded-lg bg-white/10 p-3">
				<div
					class="flex h-8 w-8 items-center justify-center rounded-full bg-white text-sm font-bold text-[#117064]"
				>
					{user?.nama_lengkap ? user.nama_lengkap.charAt(0).toUpperCase() : 'P'}
				</div>
				<div class="flex-1 overflow-hidden">
					<p class="truncate text-xs font-bold text-white">{user?.nama_lengkap || 'Petugas'}</p>
					<p class="text-[10px] text-teal-200">{user?.role?.toUpperCase() || 'KADER'}</p>
				</div>
				<button
					onclick={handleLogout}
					class="cursor-pointer text-teal-200 hover:text-white"
					title="Keluar"
				>
					<LogOut size={16} />
				</button>
			</div>
		</div>
	</aside>

	<main class="flex min-w-0 flex-1 flex-col">
		<header class="flex h-16 items-center justify-between bg-white px-6 shadow-sm">
			<div class="flex items-center gap-4">
				<button
					onclick={() => (isSidebarOpen = !isSidebarOpen)}
					class="text-gray-500 hover:text-[#117064] md:hidden"
				>
					<Menu size={24} />
				</button>
				<h1 class="text-xl font-bold text-gray-800">
					{#if currentPath === '/dashboard'}
						Beranda
					{:else if currentPath.includes('inventaris')}
						Inventaris Logistik
					{:else if currentPath.includes('pengguna-publik')}
						Daftar Warga
					{:else}
						Dasbor
					{/if}
				</h1>
			</div>

			<div class="flex items-center gap-4">
				<div class="hidden items-center gap-2 text-sm text-gray-500 md:flex">
					<CalendarDays size={16} />
					<span>{todayDate}</span>
				</div>
				<button
					class="relative cursor-pointer rounded-full border border-gray-200 p-2 text-gray-500 hover:bg-gray-50 hover:text-[#117064]"
				>
					<Bell size={18} />
					<span class="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-red-500 ring-2 ring-white"
					></span>
				</button>
			</div>
		</header>

		<div class="flex-1 overflow-y-auto p-6">
			{@render children()}
		</div>
	</main>
</div>

{#if isSidebarOpen}
	<div
		class="fixed inset-0 z-40 bg-black/50 md:hidden"
		onclick={() => (isSidebarOpen = false)}
	></div>
{/if}
