/// <reference types="vitest" />
/// <reference types="vite/client" />

import { defineConfig } from 'vite'
import svgr from 'vite-plugin-svgr'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react(), svgr()],
  envPrefix: ['VITE_', 'REACT_APP_'],
  server: {
    port: 1234,
    proxy: {
        '/photo/': 'http://localhost:4001/',
        '/graphql': {
             target: 'http://localhost:4001',
             changeOrigin: true,
             secure: false,      
             ws: true,
         }
    },
  },
  esbuild: {
    logOverride: { 'this-is-undefined-in-esm': 'silent' },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './testing/setupTests.ts',
    coverage: {
      reporter: ['text', 'json', 'html'],
    },
  },
})
