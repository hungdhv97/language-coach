import path from 'path'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: '0.0.0.0', // Allow external connections (Docker)
    port: 5173,
    allowedHosts: [
      'app.lexigo.io.vn',
      'localhost',
      '.lexigo.io.vn', // Allow all subdomains
    ],
    watch: {
      usePolling: true, // Required for HMR in Docker
    },
    hmr: {
      host: 'localhost', // HMR host
      port: 5173, // HMR port
    },
  },
})
