import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	// TODO remove later
	server: {
		proxy: {
			'/api': 'http://localhost:5000'
		}
	}
});
