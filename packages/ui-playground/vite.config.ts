/// <reference types="vitest" />
import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'

export default defineConfig({
  plugins: [solid()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    testTransformMode: {
      web: ['[jt]sx?'],
    },
    // Vitest requires explicit deps optimization for solid-js
    deps: {
      optimizer: {
        web: {
          include: ['solid-js', '@solidjs/testing-library'],
        },
      },
    },
  },
  resolve: {
    conditions: ['development', 'browser'],
  },
})
