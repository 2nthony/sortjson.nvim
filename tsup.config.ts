import { defineConfig } from 'tsup'

export default defineConfig({
  entry: ['src/cli.ts'],
  outDir: 'lua/sortjson',
  format: ['cjs'],
  target: 'node16',
  minify: true,
})
