import { defineConfig } from 'vite'

export default defineConfig({
    server: {
        proxy: {
            '/api': {
                // target: 'http://localhost:8080', // Go backend
                target: 'http://localhost:9090', // Go proxy backend
                changeOrigin: true,
                secure: false,
            }
        }
    }
})
