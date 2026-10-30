import { defineConfig } from 'vite';
import { resolve } from 'path';
import fs from 'fs';

export default defineConfig({
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    rollupOptions: {
      input: {
        popup: resolve(__dirname, 'src/popup/popup.html'),
        blocked: resolve(__dirname, 'src/blocked-page/blocked.html'),
        background: resolve(__dirname, 'src/background/index.ts'),
      },
      output: {
        entryFileNames: (chunkInfo) => {
          if (chunkInfo.name === 'background') {
            return 'background.js';
          }
          return 'assets/[name].js';
        },
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name].[ext]',
      },
    },
  },
  plugins: [
    {
      name: 'copy-manifest-and-pages',
      closeBundle() {
        if (fs.existsSync('manifest.json')) {
          fs.copyFileSync('manifest.json', 'dist/manifest.json');
        }
        if (fs.existsSync('dist/src/blocked-page/blocked.html')) {
          fs.copyFileSync('dist/src/blocked-page/blocked.html', 'dist/blocked.html');
        }
        if (fs.existsSync('dist/src/popup/popup.html')) {
          fs.copyFileSync('dist/src/popup/popup.html', 'dist/popup.html');
        }
      },
    },
  ],
});
