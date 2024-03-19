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
      reporter: ['text', 'json'],
    }
  }
})
