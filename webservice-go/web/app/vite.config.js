/** @type {import('vite').UserConfig} */

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import checker from 'vite-plugin-checker';
import eslint from 'vite-plugin-eslint';

export default defineConfig({
  plugins: [
    react(),
    eslint(),
    checker({ typescript: true }),
  ],
  test: {
    globals: true,
    coverage: {
      exclude: ['.eslintrc.cjs', '*.config.js', '*.config.cjs'],
      reporter: ['text', 'json-summary', 'json'],
      reportOnFailure: true,
    }
  }
})
