import { redirect } from '@sveltejs/kit';

export async function load() {
	// Melakukan pengalihan HTTP 307 secara instan di tingkat server menuju /login
	throw redirect(307, '/login');
}
