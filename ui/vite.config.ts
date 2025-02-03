import { defineConfig } from 'vite'

export default defineConfig({
    server: {
        proxy: {
            '/api': {
                // target: 'https://localhost:8080', // Go backend
                target: 'http://localhost:9090', // Go proxy backend
                changeOrigin: true,
                secure: false,
            }
        }
    }
})
