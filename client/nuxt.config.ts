// https://nuxt.com/docs/api/configuration/nuxt-config
// eslint-disable-next-line no-undef
export default defineNuxtConfig({
  modules: [
    '@nuxt/eslint',
    '@nuxt/icon',
    '@nuxt/ui',
    '@pinia/nuxt',
    'pinia-plugin-persistedstate/nuxt',
    '@vueuse/nuxt',
    '@vueuse/motion/nuxt',
    'nuxt-time',
  ],

  ssr: false,

  imports: {
    dirs: ['./layouts', './components', './pages', './store' ],
  },

  devtools: {
    enabled: true,

    timeline: {
      enabled: true,
    },
  },

  plugins: [
    {
      src: './app/plugins/myPrism.ts',
      mode: 'client'
    }
  ],
  css:['~/assets/css/onedark.css'],

  srcDir: 'app',
  serverDir: 'nuxt-server',
  compatibilityDate: '2024-10-24',

  eslint: {
    config: {
      standalone: false,
      nuxt: {
        sortConfigKeys: true,
      },
    },
  },
})
