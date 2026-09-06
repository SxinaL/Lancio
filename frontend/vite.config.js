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
        // 端口被占用时自动递增选择可用端口，避免与其他进程冲突
        port: 34115,
        strictPort: false,
    },
    resolve: {
    alias: {
        '@': path.resolve(__dirname, 'src'),
    }
 }
});

