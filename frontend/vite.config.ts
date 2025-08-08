import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': { // Requests starting with /api will be proxied
				target: 'http://localhost:8080', // The URL of your backend API
				changeOrigin: true, // Needed for virtual hosted sites
			},
		},
	},
});
