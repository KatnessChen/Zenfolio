import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: true, // Listen on all addresses including Docker
    port: 3000,
    strictPort: true,
    watch: {
      usePolling: true, // Better file watching in Docker
    },
    hmr: {
      clientPort: 3000, // Ensure HMR works through Docker port mapping
    },
  },
})
