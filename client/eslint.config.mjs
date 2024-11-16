// @ts-check
import antfu from '@antfu/eslint-config'
import nuxt from './.nuxt/eslint.config.mjs'

export default nuxt(
  antfu({
    formatters: true,

    ignores: ['*.config.*', 'vite-env.d.ts'],
  }),
)
