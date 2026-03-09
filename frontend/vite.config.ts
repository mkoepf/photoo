import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/thumbnail': {
        target: 'http://127.0.0.1:34115',
        changeOrigin: true,
      }
    }
  }
})
