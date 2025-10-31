import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    // Прокси для обхода CORS - перенаправляем /api на бэкенд
    proxy: {
      '/api': {
        target: 'http://localhost:8083',
        changeOrigin: true,
        secure: false,
      },
      // Прокси для MinIO изображений
      '/lab1': {
        target: 'http://localhost:9003',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
