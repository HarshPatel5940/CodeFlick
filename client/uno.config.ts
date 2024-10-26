import { defineConfig, presetAttributify, presetUno } from 'unocss'

export default defineConfig({
  theme: {
    fontFamily: {
      satoshi: 'Satoshi-Variable',
    },
  },
  presets: [
    presetUno(),
    presetAttributify(),
  ],
})
