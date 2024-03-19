/** @type {import('vite').UserConfig} */

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import checker from 'vite-plugin-checker';

export default defineConfig({
  plugins: [
    react(),
    checker({ typescript: true }),
  ],
  test: {
    globals: true,
    coverage: {
      reporter: ['text', 'json'],
    }
  }
})
