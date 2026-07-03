<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchAPI } from '$lib/api';
	import { UserPlus, Key, Edit, ShieldCheck, Search, AlertCircle, X } from 'lucide-svelte';

	interface UserData {
		id: string;
		username: string;
		nama_lengkap: string;
		role: string;
		is_active: boolean;
	}

	// --- STATE MANAGEMENT ---
	let users = $state<UserData[]>([]);
	let isLoading = $state(true);
	let searchQuery = $state('');

	// Modal States
	let showFormModal = $state(false);
	let showResetModal = $state(false);
	let modalMode = $state<'tambah' | 'edit'>('tambah');
	let isSubmitting = $state(false);
	let errorMessage = $state('');

	// Form Data
	let formData = $state({
		id: '',
		username: '',
		password: '',
		nama_lengkap: '',
		role: 'KADER',
		is_active: true
	});

	let resetData = $state({
		id: '',
		username: '',
		new_password: ''
	});

	// --- API CALLS ---
	async function loadUsers() {
		isLoading = true;
		try {
			const res = await fetchAPI('/pengguna');
			users = res.data || [];
		} catch (error: unknown) {
			console.error('Gagal memuat pengguna:', (error as Error).message);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadUsers();
	});

	// --- HANDLERS ---
	function openAddModal() {
		modalMode = 'tambah';
		formData = {
			id: '',
			username: '',
			password: '',
			nama_lengkap: '',
			role: 'KADER',
			is_active: true
		};
		errorMessage = '';
		showFormModal = true;
	}

	function openEditModal(u: UserData) {
		modalMode = 'edit';
		formData = {
			id: u.id,
			username: u.username,
			password: '', // Tidak dikirim saat edit
			nama_lengkap: u.nama_lengkap,
			role: u.role,
			is_active: u.is_active
		};
		errorMessage = '';
		showFormModal = true;
	}

	function openResetModal(u: UserData) {
		resetData = { id: u.id, username: u.username, new_password: '' };
		errorMessage = '';
		showResetModal = true;
	}

	async function handleFormSubmit(e: Event) {
		e.preventDefault();
		isSubmitting = true;
		errorMessage = '';

		try {
			if (modalMode === 'tambah') {
				await fetchAPI('/pengguna', {
					method: 'POST',
					body: JSON.stringify({
						username: formData.username,
						password: formData.password,
						nama_lengkap: formData.nama_lengkap,
						role: formData.role
					})
				});
			} else {
				await fetchAPI(`/pengguna/${formData.id}`, {
					method: 'PUT',
					body: JSON.stringify({
						nama_lengkap: formData.nama_lengkap,
						role: formData.role,
						is_active: formData.is_active
					})
				});
			}
			showFormModal = false;
			loadUsers();
		} catch (error: unknown) {
			errorMessage = (error as Error).message || 'Terjadi kesalahan saat menyimpan data.';
		} finally {
			isSubmitting = false;
		}
	}

	async function handleResetSubmit(e: Event) {
		e.preventDefault();
		isSubmitting = true;
		errorMessage = '';
		try {
			await fetchAPI(`/pengguna/${resetData.id}/reset-password`, {
				method: 'PUT',
				body: JSON.stringify({ new_password: resetData.new_password })
			});
			showResetModal = false;
			alert(`Kata sandi untuk ${resetData.username} berhasil di-reset!`);
		} catch (error: unknown) {
			errorMessage = (error as Error).message || 'Gagal mereset kata sandi.';
		} finally {
			isSubmitting = false;
		}
	}

	// Dynamic Filter
	let filteredUsers = $derived(
		users.filter(
			(u) =>
				u.nama_lengkap.toLowerCase().includes(searchQuery.toLowerCase()) ||
				u.username.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	// Helper Styling Badges
	function getRoleStyle(role: string) {
		if (role === 'ADMIN') return 'bg-slate-800 text-slate-100 border-slate-700';
		if (role === 'BIDAN') return 'bg-blue-50 text-blue-700 border-blue-200';
		return 'bg-teal-50 text-teal-700 border-teal-200';
	}
</script>

<svelte:head>
	<title>Manajemen Pengguna - POSYANDU Sehat Bersama</title>
</svelte:head>

<div class="space-y-6">
	<!-- HEADER KONSOL ADMIN (Desain Slate Gelap) -->
	<div
		class="relative flex flex-col items-start justify-between gap-4 overflow-hidden rounded-2xl bg-slate-900 p-6 shadow-xl md:flex-row md:items-center"
	>
		<!-- Ornamen Background -->
		<div
			class="pointer-events-none absolute top-0 right-0 h-64 w-64 rounded-full bg-teal-500/10 blur-3xl"
		></div>

		<div class="relative z-10">
			<div class="flex items-center gap-3">
				<div
					class="flex h-10 w-10 items-center justify-center rounded-xl border border-slate-700 bg-slate-800"
				>
					<ShieldCheck size={24} class="text-teal-400" />
				</div>
				<div>
					<h1 class="text-xl font-black tracking-wider text-white uppercase">Otoritas Pengguna</h1>
					<p class="text-xs text-slate-400">Pusat kontrol akses & kredensial sistem Posyandu</p>
				</div>
			</div>
		</div>
		<button
			onclick={openAddModal}
			class="relative z-10 flex cursor-pointer items-center gap-2 rounded-xl bg-teal-500 px-5 py-2.5 text-sm font-bold text-slate-900 shadow-lg shadow-teal-500/20 transition hover:bg-teal-400 active:scale-95"
		>
			<UserPlus size={18} />
			REGISTRASI PETUGAS
		</button>
	</div>

	<!-- TABEL & PENCARIAN -->
	<div class="rounded-2xl border border-gray-200 bg-white shadow-sm">
		<div class="border-b border-gray-100 p-5">
			<div class="relative max-w-md">
				<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
					<Search size={18} class="text-gray-400" />
				</div>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Cari nama atau username pengguna..."
					class="block w-full rounded-xl border border-gray-200 bg-gray-50/50 py-2.5 pr-4 pl-11 text-sm text-gray-800 transition focus:border-teal-500 focus:bg-white focus:ring-1 focus:ring-teal-500 focus:outline-none"
				/>
			</div>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead class="bg-gray-50/80 text-xs font-bold tracking-wider text-gray-500 uppercase">
					<tr>
						<th class="px-6 py-4">Nama & Username</th>
						<th class="px-6 py-4 text-center">Tingkat Akses</th>
						<th class="px-6 py-4 text-center">Status</th>
						<th class="px-6 py-4 text-right">Tindakan Khusus</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100">
					{#if isLoading}
						<tr
							><td colspan="4" class="py-10 text-center text-gray-400"
								>Menyinkronkan data kredensial...</td
							></tr
						>
					{:else if filteredUsers.length === 0}
						<tr
							><td colspan="4" class="py-10 text-center text-gray-400"
								>Tidak ada pengguna yang cocok.</td
							></tr
						>
					{/if}

					{#each filteredUsers as u (u.id)}
						<tr class="transition hover:bg-gray-50/50">
							<td class="px-6 py-4">
								<div class="font-bold text-gray-800">{u.nama_lengkap}</div>
								<div class="text-xs font-medium text-gray-500">@{u.username}</div>
							</td>
							<td class="px-6 py-4 text-center">
								<span
									class={`inline-flex rounded-lg border px-2.5 py-1 text-[10px] font-bold tracking-wider uppercase ${getRoleStyle(u.role)}`}
								>
									{u.role}
								</span>
							</td>
							<td class="px-6 py-4 text-center">
								{#if u.is_active}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-green-50 px-3 py-1 text-xs font-bold text-green-700"
									>
										<span class="h-1.5 w-1.5 rounded-full bg-green-500"></span> Aktif
									</span>
								{:else}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-3 py-1 text-xs font-bold text-slate-500"
									>
										<span class="h-1.5 w-1.5 rounded-full bg-slate-400"></span> Nonaktif
									</span>
								{/if}
							</td>
							<td class="px-6 py-4 text-right">
								<div class="flex items-center justify-end gap-2">
									<button
										onclick={() => openEditModal(u)}
										class="cursor-pointer rounded-lg bg-gray-100 p-2 text-gray-600 transition hover:bg-teal-50 hover:text-teal-600"
										title="Edit Profil/Status"
									>
										<Edit size={16} />
									</button>
									<button
										onclick={() => openResetModal(u)}
										class="cursor-pointer rounded-lg bg-gray-100 p-2 text-gray-600 transition hover:bg-amber-50 hover:text-amber-600"
										title="Reset Kata Sandi"
									>
										<Key size={16} />
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

<!-- ============================================== -->
<!-- MODAL TAMBAH / EDIT PENGGUNA (Slate Theme)     -->
<!-- ============================================== -->
{#if showFormModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 p-4 backdrop-blur-sm"
	>
		<div class="w-full max-w-lg overflow-hidden rounded-3xl bg-white shadow-2xl">
			<div
				class="flex items-center justify-between border-b border-gray-100 bg-slate-900 px-6 py-5"
			>
				<div class="flex items-center gap-3">
					<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-slate-800">
						{#if modalMode === 'tambah'}
							<UserPlus size={16} class="text-teal-400" />
						{:else}
							<Edit size={16} class="text-teal-400" />
						{/if}
					</div>
					<h2 class="text-lg font-black tracking-wide text-white uppercase">
						{modalMode === 'tambah' ? 'Registrasi Petugas Baru' : 'Perbarui Kredensial'}
					</h2>
				</div>
				<button
					onclick={() => (showFormModal = false)}
					class="cursor-pointer text-slate-400 hover:text-white"
				>
					<X size={24} />
				</button>
			</div>

			<form onsubmit={handleFormSubmit} class="p-6">
				{#if errorMessage}
					<div
						class="mb-5 flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 p-4 text-red-600"
					>
						<AlertCircle size={18} class="mt-0.5 shrink-0" />
						<p class="text-sm font-medium">{errorMessage}</p>
					</div>
				{/if}

				<div class="space-y-4">
					<div>
						<label class="mb-1.5 block text-xs font-bold text-gray-600">Nama Lengkap</label>
						<input
							type="text"
							required
							bind:value={formData.nama_lengkap}
							class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-teal-500 focus:bg-white focus:ring-1 focus:ring-teal-500 focus:outline-none"
						/>
					</div>

					<div class="grid grid-cols-2 gap-4">
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-600">ID Username</label>
							<input
								type="text"
								required
								disabled={modalMode === 'edit'}
								bind:value={formData.username}
								class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-teal-500 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-100 disabled:text-gray-400"
							/>
						</div>
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-600"
								>Tingkat Akses (Role)</label
							>
							<select
								bind:value={formData.role}
								class="w-full cursor-pointer rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-teal-500 focus:outline-none"
							>
								<option value="KADER">KADER (Operasional)</option>
								<option value="BIDAN">BIDAN (Medis)</option>
								<option value="ADMIN">ADMINISTRATOR (IT)</option>
							</select>
						</div>
					</div>

					{#if modalMode === 'tambah'}
						<div>
							<label class="mb-1.5 block text-xs font-bold text-gray-600">Kata Sandi Awal</label>
							<input
								type="password"
								required
								bind:value={formData.password}
								class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-teal-500 focus:bg-white focus:ring-1 focus:ring-teal-500 focus:outline-none"
							/>
						</div>
					{/if}

					{#if modalMode === 'edit'}
						<div
							class="mt-2 flex items-center justify-between rounded-xl border border-gray-200 bg-gray-50 p-4"
						>
							<div>
								<p class="text-sm font-bold text-gray-800">Status Akun Aktif</p>
								<p class="text-[10px] text-gray-500">
									Matikan jika petugas sudah berhenti bertugas.
								</p>
							</div>
							<label class="relative inline-flex cursor-pointer items-center">
								<input type="checkbox" bind:checked={formData.is_active} class="peer sr-only" />
								<div
									class="peer h-6 w-11 rounded-full bg-gray-300 peer-checked:bg-teal-500 after:absolute after:top-[2px] after:left-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:after:translate-x-full peer-checked:after:border-white"
								></div>
							</label>
						</div>
					{/if}
				</div>

				<div class="mt-8 flex items-center justify-end gap-3 border-t border-gray-100 pt-5">
					<button
						type="button"
						onclick={() => (showFormModal = false)}
						class="cursor-pointer rounded-xl px-5 py-2.5 text-sm font-bold text-gray-600 hover:bg-gray-100"
					>
						Batalkan
					</button>
					<button
						type="submit"
						disabled={isSubmitting}
						class="cursor-pointer rounded-xl bg-slate-900 px-6 py-2.5 text-sm font-bold text-white transition hover:bg-slate-800 disabled:opacity-70"
					>
						{isSubmitting ? 'Menyimpan...' : 'Simpan Kredensial'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- ============================================== -->
<!-- MODAL RESET PASSWORD (Danger/Amber Theme)      -->
<!-- ============================================== -->
{#if showResetModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 p-4 backdrop-blur-sm"
	>
		<div class="w-full max-w-sm overflow-hidden rounded-3xl bg-white shadow-2xl">
			<form onsubmit={handleResetSubmit} class="p-6">
				<div
					class="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-amber-100 text-amber-600"
				>
					<Key size={28} />
				</div>
				<h2 class="mb-2 text-xl font-black text-gray-800">Reset Sandi</h2>
				<p class="mb-6 text-sm text-gray-500">
					Anda akan mereset kata sandi untuk akun <span class="font-bold text-gray-800"
						>@{resetData.username}</span
					>.
				</p>

				{#if errorMessage}
					<p class="mb-4 text-xs font-medium text-red-500">{errorMessage}</p>
				{/if}

				<div>
					<label class="mb-1.5 block text-xs font-bold text-gray-600">Kata Sandi Baru</label>
					<input
						type="text"
						required
						bind:value={resetData.new_password}
						placeholder="Masukkan sandi baru..."
						class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-amber-500 focus:bg-white focus:ring-1 focus:ring-amber-500 focus:outline-none"
					/>
				</div>

				<div class="mt-8 flex items-center gap-3">
					<button
						type="button"
						onclick={() => (showResetModal = false)}
						class="flex-1 cursor-pointer rounded-xl bg-gray-100 py-2.5 text-sm font-bold text-gray-600 hover:bg-gray-200"
					>
						Batal
					</button>
					<button
						type="submit"
						disabled={isSubmitting}
						class="flex-1 cursor-pointer rounded-xl bg-amber-500 py-2.5 text-sm font-bold text-white transition hover:bg-amber-600 disabled:opacity-70"
					>
						{isSubmitting ? 'Proses...' : 'Reset Sandi'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
