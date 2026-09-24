import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 前端 5173，/api 与 /uploads 代理到 Go 后端 8080。
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '127.0.0.1',
    port: 5173,
    strictPort: true,
    proxy: {
      '/api': { target: 'http://127.0.0.1:8080', changeOrigin: true },
      '/uploads': { target: 'http://127.0.0.1:8080', changeOrigin: true }
    }
  }
})
