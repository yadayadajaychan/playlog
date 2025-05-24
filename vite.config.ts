import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import fs from 'fs';

const programVersion = fs.readFileSync('VERSION', 'utf-8').trim()

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': 'http://localhost:5000'
		}
	},
	define: {
		__APP_VERSION__: JSON.stringify(programVersion),
	},
});
