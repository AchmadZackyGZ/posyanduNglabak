<script lang="ts">
	// Memuat gaya CSS global
	import './layout.css';

	// Menerima komponen halaman anak di Svelte 5
	let { children } = $props();
</script>

{@render children()}
