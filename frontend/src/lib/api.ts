import { PUBLIC_BASE_API_URL } from '$env/static/public';

// Gunakan URL dari .env, atau fallback jika belum tersetting
const BASE_URL = PUBLIC_BASE_API_URL || 'http://localhost:8080/api/v1';

export async function fetchAPI(endpoint: string, options: RequestInit = {}) {
	// Ambil token JWT dari localStorage yang disimpan saat login
	const token = localStorage.getItem('token');

	// FIX TYPINGS: Gunakan Record<string, string> alih-alih HeadersInit
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		'ngrok-skip-browser-warning': 'true',
		...(options.headers as Record<string, string>)
	};

	// Jika token ada, suntikkan ke header Authorization
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	try {
		const response = await fetch(`${BASE_URL}${endpoint}`, {
			...options,
			headers
		});

		// Tangani jika token kedaluwarsa atau tidak valid (401 Unauthorized)
		if (response.status === 401) {
			localStorage.removeItem('token');
			localStorage.removeItem('user');
			window.location.href = '/login'; // Tendang kembali ke halaman login
			throw new Error('Sesi telah berakhir, silakan login kembali.');
		}

		const data = await response.json();

		if (!response.ok) {
			throw new Error(data.error || 'Terjadi kesalahan pada server');
		}

		return data;
	} catch (error) {
		console.error(`Fetch API Error (${endpoint}):`, error);
		throw error;
	}
}
