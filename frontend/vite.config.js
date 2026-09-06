import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';
export default defineConfig({
    plugins: [vue()],
    build: {
        outDir: 'dist',
        assetsDir: 'assets',
    },
    server: {
        port: 34115,
        strictPort: true,
    },
    resolve: {
    alias: {
        '@': path.resolve(__dirname, 'src'),
    }
 }
});

